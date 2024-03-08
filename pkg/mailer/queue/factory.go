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
	"os"
	"text/template"
	"time"
)

// OptFunc represents a function type for configuring a Queue.
type OptFunc func(*Queue) error

// WithWorkers sets the number of simultaneous email sender
// workers for the Queue.
func WithWorkers(num uint8) OptFunc {
	return func(q *Queue) error {
		if num <= 0 {
			return nil
		}
		q.workers = num
		return nil
	}
}

// WithReadReceipts sets the email address which is to receive
// the "Read receipt" notifications from the sent emails.
func WithReadReceipts(e string) OptFunc {
	return func(q *Queue) error {
		q.readReceipts = e
		return nil
	}
}

// WithHTML sets the HTML content for the emails to be sent
// by the Queue.
func WithHTML(file string) OptFunc {
	return func(q *Queue) error {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil
		}

		t, err := template.New("html").Parse(string(b))
		if err != nil {
			return nil
		}
		q.html = t

		return nil
	}
}

// WithRateMinute sets the maximum number of emails that can be sent by
// a single sender in a minute.
func WithRateMinute(rate uint16) OptFunc {
	return func(q *Queue) error {
		if rate <= 0 {
			return nil
		}
		q.perMinute = rate
		return nil
	}
}

// WithRateDaily sets the maximum number of emails that can be sent by
// a single sender in a day.
func WithRateDaily(rate uint16) OptFunc {
	return func(q *Queue) error {
		if rate <= 0 {
			q.perDay = 100
			return nil
		}
		q.perDay = rate
		return nil
	}
}

func defaultQueue() *Queue {
	return &Queue{
		perDay:         100,
		perMinute:      2,
		workers:        2,
		start:          time.Now(),
		status:         make(map[string]*Stats),
		errorThreshold: 6,
		errorCount:     0,
	}
}

// New constructs an instance of [queue.Queue] with the provided options.
func New(senders, receivers, subject, host, textFile string, opts ...OptFunc) (*Queue, error) {
	q := defaultQueue()

	if err := q.load(senders, receivers); err != nil {
		return nil, err
	}

	q.subject = subject
	q.host = host

	for _, sender := range q.senders {
		q.status[sender.Email] = &Stats{Sender: sender.Email}
	}

	b, err := os.ReadFile(textFile)
	if err != nil {
		return nil, err
	}

	t, err := template.New("text").Parse(string(b))
	if err != nil {
		return nil, err
	}
	q.text = t

	for _, optFn := range opts {
		if err := optFn(q); err != nil {
			return nil, err
		}
	}

	return q, nil
}
