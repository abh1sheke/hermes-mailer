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

import "time"

// Stats represents a unique senders statistics including
// total emails successfully sent, emails that have failed to
// send as well as bounces.
type Stats struct {
	skip    bool
	today   uint
	timeout *time.Time
	Sender  string `csv:"sender"`
	Total   uint   `csv:"total"`
	Failed  uint   `csv:"failed"`
	Bounced uint   `csv:"bounced"`
}

func (s *Stats) isTimedOut() bool {
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

func (s *Stats) setTimeout(d time.Duration) {
	t := time.Now().Add(d)
	s.timeout = &t
}

func (s *Stats) increment(num uint) {
	s.today += num
	s.Total += num
}

func (s *Stats) incrementFailed(num uint) {
	s.Failed += num
}

func (s *Stats) incrementBounced(num uint) {
	s.Bounced += num
}
