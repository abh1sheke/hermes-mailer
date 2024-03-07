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
	"errors"
	"fmt"
	"net/smtp"
)

// Auth is an intermediary representation of a
// type of supported [Auth] authentication mechanism
type Auth uint8

const (
	// Plain represents [Auth] PLAIN authentication
	Plain Auth = iota
	// Login represents [Auth] LOGIN authentication
	Login
	// CRAMMD5 represents [Auth] CRAM-MD5 authentication
	CRAMMD5
)

type loginAuth struct {
	username, password, host string
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown 'fromServer': %q", string(fromServer))
		}
	}

	return nil, nil
}

// LoginAuth returns an [Auth] that implements the LOGIN authentication
// mechanism. The returned Auth uses the given username and password to
// authenticate to host.
//
// LoginAuth will only send the credentials if the connection is using TLS
// or is connected to localhost. Otherwise authentication will fail with an
// error, without sending the credentials.
func LoginAuth(username, password, host string) smtp.Auth {
	return &loginAuth{username, password, host}
}
