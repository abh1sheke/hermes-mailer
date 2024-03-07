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
	"fmt"
	"net/textproto"
	"strings"
	"sync"

	"github.com/abh1sheke/hermes-mailer/pkg/mailer"
	"github.com/jordan-wright/email"
)

// resultKind represents the outcome of a result.
type resultKind uint

const (
	// success represents a positive result
	success resultKind = iota
	// failure represents a negative (errored) result
	failure
)

// workerResult represents the result of a worker thread's operation.
type workerResult struct {
	kind      resultKind
	sender    string
	error     error
	sent      uint
	receivers []*mailer.Receiver
}

func createEmails(task *task, from string) ([]*email.Email, error) {
	emails := make([]*email.Email, 0, len(task.receivers))

	for _, receiver := range task.receivers {
		e := &email.Email{
			From:    from,
			To:      []string{receiver.Email},
			Subject: task.subject,
			Headers: make(textproto.MIMEHeader),
		}

		if receiver.Bcc != nil {
			e.Cc = receiver.Cc.Data()
		}

		if receiver.Bcc != nil {
			e.Bcc = receiver.Bcc.Data()
		}

		var data map[string]string
		if receiver.Variables != nil {
			data = receiver.Variables.Data()
		} else {
			data = make(map[string]string)
		}

		text := new(strings.Builder)
		err := task.text.Execute(text, data)
		if err != nil {
			return nil, err
		}
		e.Text = []byte(text.String())

		if task.html != nil {
			html := new(strings.Builder)
			err := task.html.Execute(html, data)
			if err != nil {
				return nil, err
			}
			e.HTML = []byte(html.String())
		}

		if task.readReceipt != "" {
			e.Headers.Add("Disposition-Notification-To", task.readReceipt)
			e.Headers.Add("Return-Receipt-To", task.readReceipt)
		}

		emails = append(emails, e)
	}

	return emails, nil
}

func worker(task *task, auth mailer.Auth, res chan workerResult, wg *sync.WaitGroup) {
	defer wg.Done()

	var from string
	if task.sender.Name != "" {
		from = fmt.Sprintf("%s <%s>", task.sender.Name, task.sender.Email)
	} else {
		from = task.sender.Email
	}

	emails, err := createEmails(task, from)
	if err != nil {
		res <- workerResult{kind: failure, error: err, receivers: task.receivers}
		return
	}

	idx, err := mailer.SendEmailsTLS(task.sender, emails, task.host, auth)
	if err != nil {
		res <- workerResult{kind: failure, error: err, receivers: task.receivers[idx:]}
		return
	}

	res <- workerResult{kind: success}
}
