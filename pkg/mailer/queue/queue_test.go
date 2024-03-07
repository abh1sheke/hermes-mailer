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
	"strconv"
	"sync"
	"testing"
	"text/template"

	"github.com/abh1sheke/hermes-mailer/pkg/mailer"
)

func TestWorker(t *testing.T) {
	senderEmail, ok := os.LookupEnv("SENDER")
	if !ok {
		t.Fatalf("'SENDER' value missing from ENV")
	}

	senderPass, ok := os.LookupEnv("PASS")
	if !ok {
		t.Fatalf("'PASS' value missing from ENV")
	}

	receiverEmail, ok := os.LookupEnv("RECEIVER")
	if !ok {
		t.Fatalf("'RECEIVER' value missing from ENV")
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		t.Fatalf("'HOST' value missing from ENV")
	}

	sender := &mailer.Sender{
		Email:    senderEmail,
		Password: senderPass,
	}

	receiver := &mailer.Receiver{Email: receiverEmail}

	msg := "Greetings,\n\nThis is a test email sent from Go.\n\nThank you,\nAbhisheke."
	text, err := template.New("text").Parse(msg)
	if err != nil {
		t.Fatal(err)
	}

	task := &task{
		sender:    sender,
		receivers: []*mailer.Receiver{receiver},
		subject:   "This is a test email",
		host:      host,
		text:      text,
	}

	wg := new(sync.WaitGroup)
	res := make(chan workerResult, 1)

	wg.Add(1)
	go worker(task, mailer.Plain, res, wg)

	wg.Wait()
	close(res)

	for r := range res {
		if r.kind == failure {
			t.Fatal(r.error)
		}
	}
}

func TestQueue(t *testing.T) {
	senders, ok := os.LookupEnv("SENDERS")
	if !ok {
		t.Fatalf("'SENDERS' value missing from ENV")
	}
	receivers, ok := os.LookupEnv("RECEIVERS")
	if !ok {
		t.Fatalf("'RECEIVERS' value missing from ENV")
	}

	var workers uint8
	s, ok := os.LookupEnv("WORKERS")
	if !ok {
		workers = 2
	} else {
		var err error
		w, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			t.Fatal(err)
		}
		workers = uint8(w)
	}

	q, err := New(
		senders,
		receivers,
		"This is to test the queue functionality",
		"",
		"../../../examples/text_templ.txt",
		WithHTML("../../../examples/html_templ.example.html"),
		WithWorkers(workers),
		WithReadReceipts(os.Getenv("READ_RECEIPTS")),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := q.Run(); err != nil {
		t.Fatal(err)
	}
}
