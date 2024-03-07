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

package mailer

import (
	"reflect"
	"testing"
)

func TestReadSender(t *testing.T) {
	data, err := ReadFile[Sender]("../../examples/senders.example.csv")
	if err != nil {
		t.Fatal(err)
	}

  expected := &Sender{
    Email: "emma_anderson@example.com",
    Password: "AndersonPW456",
    Name: "Emma Anderson",
  }

  if !reflect.DeepEqual(expected, data[9]) {
    t.Fatalf("expected: %+v\ngot: %+v\n", expected, data[9])
  }
}

func TestReadReceiver(t *testing.T) {
	data, err := ReadFile[Receiver]("../../examples/receivers.example.csv")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Receiver{
		Email:     "sarah@example.com",
		Cc:        &List{[]string{"tom@example.com"}},
		Bcc:       &List{[]string{"mark@example.com", "emma@example.com"}},
		Variables: &Variables{map[string]string{"name": "Sarah", "location": "Paris"}},
	}

  if !reflect.DeepEqual(expected, data[4]) {
    t.Fatalf("expected: %+v\ngot: %+v\n", expected, data[4])
  }
}
