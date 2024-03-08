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
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

// SendEmailsTLS sends a list of email messages, defined as
// instances of [github.com/jordan-wright/email.Email], over SMTP.
//
// It returns the slice index of the email message that failed to be
// sent as well as the reason for the failure.
//
// # Arguments:
//
//   - sender: An instance of [mailer.Sender]
//
//   - emails: A slice of [github.com/jordan-wright/email.Email] instances
//
//   - host: The senders SMTP host address
//
//   - auth: An instance of [mailer.Auth] (authentication mechanism such as PLAIN, LOGIN, etc,.)
func SendEmailsTLS(sender *Sender, emails []*email.Email, host string, auth Auth) (int, error) {
	var a smtp.Auth
	switch auth {
	case Plain:
		a = smtp.PlainAuth("", sender.Email, sender.Password, host)
	case Login:
		a = LoginAuth(sender.Email, sender.Password, host)
	case CRAMMD5:
		a = smtp.CRAMMD5Auth(sender.Email, sender.Password)
	}

	addr := fmt.Sprintf("%s:587", host)
	log.Debug().Msgf("creating conn pool for: %s", sender.Email)
	pool, err := email.NewPool(addr, len(emails), a)
	if err != nil {
		return 0, err
	}

	for i, email := range emails {
		log.Info().Str("from", sender.Email).Str("to", email.To[0]).Msg("sending email")
		_ = pool.Send(email, 10*time.Second)
		if err != nil {
			return i, err
		}
	}

	return 0, nil
}
