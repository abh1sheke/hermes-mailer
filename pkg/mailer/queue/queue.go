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
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/abh1sheke/hermes-mailer/pkg/mailer"
	"github.com/gocarina/gocsv"
)

type task struct {
	sender                     *mailer.Sender
	receivers                  []*mailer.Receiver
	subject, host, readReceipt string
	text, html                 *template.Template
}

type status struct {
	timeout      *time.Time
	total, today uint
	skip         bool
}

func (s *status) isTimedOut() bool {
	if s.timeout == nil {
		return false
	}

	now := time.Now()

	if now.Unix() > s.timeout.Unix() {
		s.timeout = nil
		return false
	}

	return true
}

func (s *status) setTimeout(d time.Duration) {
	t := time.Now().Add(d)
	s.timeout = &t
}

func (s *status) increment(num uint) {
	s.today += num
	s.total += num
}

// Queue represents a worker queue performing email send operations.
type Queue struct {
	senders                     []*mailer.Sender
	receivers                   []*mailer.Receiver
	subject, host, readReceipts string
	text, html                  *template.Template
	perMinute, perDay           uint16
	start                       time.Time
	status                      map[string]*status
	workers                     uint8
	auth                        mailer.Auth
	failures                    []*mailer.Receiver
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

	var err error

	for res := range res {
		switch res.kind {
		case success:
			status := q.status[res.sender]
			status.increment(res.sent)
			status.setTimeout(1 * time.Minute)

		case failure:
			//err = res.Error
			q.failures = append(q.failures, res.receivers...)
		}
	}

	return err
}

func (q *Queue) saveFailures() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filename := filepath.Join(pwd, "errored_receivers.csv")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	return gocsv.MarshalFile(&q.failures, f)
}

// Run initates the mail queue and performs all the specified
// operations based off of the given queue configuration.
func (q *Queue) Run() error {
	wg := new(sync.WaitGroup)

	receiverPtr := 0
	for receiverPtr < len(q.receivers) {
		senderPtr := 0
		res := make(chan workerResult, q.workers)
		for i := 0; i < int(q.workers); i++ {
			if receiverPtr > len(q.receivers) {
				break
			}

			sender := q.senders[senderPtr]
			status := q.status[sender.Email]
			if status.skip {
				continue
			}

			if status.isTimedOut() {
				continue
			}

			if !q.isTomorrow() && status.today > uint(q.perDay) {
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

	return nil
}
