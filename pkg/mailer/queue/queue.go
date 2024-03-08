// Copyright 2024 Abhisheke Acharya
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/abh1sheke/hermes-mailer/pkg/mailer"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

type task struct {
	sender                     *mailer.Sender
	receivers                  []*mailer.Receiver
	subject, host, readReceipt string
	text, html                 *template.Template
}

// Queue represents a worker queue performing email send operations.
type Queue struct {
	senders                     []*mailer.Sender
	receivers                   []*mailer.Receiver
	subject, host, readReceipts string
	text, html                  *template.Template
	perMinute, perDay           uint16
	start                       time.Time
	status                      map[string]*Stats
	workers                     uint8
	auth                        mailer.Auth
	failures                    []*mailer.Receiver
	errorThreshold, errorCount  uint8
}

func (q *Queue) isTomorrow() bool {
	now := time.Now()
	if now.Unix() > q.start.Add(24*time.Hour).Unix() {
		q.start = now
		return true
	}
	return false
}

func (q *Queue) load(senders, receivers string) error {
	s, err := mailer.ReadFile[mailer.Sender](senders)
	if err != nil {
		return err
	}

	r, err := mailer.ReadFile[mailer.Receiver](receivers)
	if err != nil {
		return err
	}
	q.senders = s
	q.receivers = r

	return nil
}

func (q *Queue) collectResults(res chan workerResult, wg *sync.WaitGroup) error {
	wg.Wait()
	close(res)

	for res := range res {
		switch res.kind {
		case success:
			if q.errorCount > 0 {
				q.errorCount = 0
			}
			status := q.status[res.sender]
			status.increment(res.sent)
			status.setTimeout(1 * time.Minute)
			log.Debug().Str("sender", res.sender).Uint("sent", res.sent).Msg("send success")

		case failure:
			log.Error().Str("from", res.sender).Uint("sent", res.sent).Err(res.error).Msg("send failure")
			status := q.status[res.sender]
			status.incrementFailed(uint(len(res.receivers)))
			q.failures = append(q.failures, res.receivers...)
			q.errorCount++

			if q.errorCount >= q.errorThreshold {
				return errors.New("queue has errored too many times")
			}
		}
	}
	return nil
}

func (q *Queue) saveFailures() (err error) {
	if len(q.failures) == 0 {
		return nil
	}
	var pwd string
	pwd, err = os.Getwd()
	if err != nil {
		return err
	}

	filename := filepath.Join(pwd, "errored_receivers.csv")
	log.Info().Str("file", filename).Msg("saving failed receivers")
	var file *os.File
	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("could not save errors")
		}
		file.Close()
	}()

	err = gocsv.MarshalFile(&q.failures, file)

	return
}

// Run initates the mail queue and performs all the specified
// operations based off of the given queue configuration.
func (q *Queue) Run() error {
	wg := new(sync.WaitGroup)
	var receiverPtr, senderPtr, skips int
	for receiverPtr < len(q.receivers) {
		res := make(chan workerResult, q.workers)
		for i := 0; i < int(q.workers); i++ {
			if receiverPtr > len(q.receivers) {
				break
			}

			sender := q.senders[senderPtr]
			status := q.status[sender.Email]
			if status.skip {
				log.Warn().Msgf("skipping risky sender: %s", sender.Email)
				time.Sleep(2 * time.Second)
				senderPtr = (senderPtr + 1) % len(q.senders)
				continue
			}

			if status.isTimedOut() {
				dur := time.Until(*status.timeout)
				log.Warn().
					Str("sender", sender.Email).
					Str("dur", fmt.Sprintf("%.2f min", dur.Minutes())).
					Msgf("sender timed out")

				skips++
				if skips >= len(q.senders) {
					log.Warn().
						Str("dur", fmt.Sprintf("%.2f min", dur.Minutes())).
						Msg("sleeping due to repeated skips")
					time.Sleep(dur)
					skips = 0
				}
				time.Sleep(2 * time.Second)

				senderPtr = (senderPtr + 1) % len(q.senders)
				continue
			}

			if !q.isTomorrow() && status.today >= uint(q.perDay) {
				status.setTimeout(24 * time.Hour)

				dur := time.Until(*status.timeout)
				log.Warn().
					Str("sender", sender.Email).
					Str("dur", fmt.Sprintf("%.2f min", dur.Minutes())).
					Msgf("daily limit exceeded")

				skips++
				if skips >= len(q.senders) {
					log.Warn().
						Str("dur", fmt.Sprintf("%.2f min", dur.Minutes())).
						Msg("sleeping due to repeated skips")
					time.Sleep(dur)
					skips = 0
				}
				time.Sleep(2 * time.Second)

				senderPtr = (senderPtr + 1) % len(q.senders)
				continue
			}

			var receivers []*mailer.Receiver
			end := receiverPtr + int(q.perMinute)

			if end > len(q.receivers) {
				receivers = q.receivers[receiverPtr:]
			} else {
				receivers = q.receivers[receiverPtr:end]
			}
			receiverPtr = end

			task := &task{
				sender:    sender,
				receivers: receivers,
				subject:   q.subject,
				host:      q.host,
				text:      q.text,
				html:      q.html,
			}

			wg.Add(1)
			go worker(task, q.auth, res, wg)

			senderPtr = (senderPtr + 1) % len(q.senders)
		}

		if err := q.collectResults(res, wg); err != nil {
			return errors.Join(err, q.saveFailures())
		}
	}

	return q.saveFailures()
}
