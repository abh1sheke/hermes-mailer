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
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

// List is a convenience type that unmarshals a CSV string and converts it
// into a of slice of strings, and vice versa.
type List struct {
	data []string
}

// UnmarshalCSV is a helper method which unmarshals a CSV string
// into a List struct with []string.
func (l *List) UnmarshalCSV(csv string) error {
	if len(csv) == 0 {
		return nil
	}

	l.data = strings.Split(csv, ";")
	return nil
}

// MarshalCSV is a helper method which marshals the data
// from a List struct into a CSV string.
func (l *List) MarshalCSV() (string, error) {
	if l.data == nil {
		return "", nil
	}

	builder := new(strings.Builder)
	for i, v := range l.data {
		if _, err := builder.WriteString(v); err != nil {
			return "", err
		}
		if i < len(l.data)-1 {
			if _, err := builder.WriteRune(';'); err != nil {
				return "", err
			}
		}
	}

	return builder.String(), nil
}

// Data is a getter method for the underlying struct data.
func (l *List) Data() []string {
	return l.data
}

// Variables is a convenience type that unmarshals a CSV string
// with the format "KEY=VALUE" and converts it into a of map of
// string keys and values, and vice versa.
type Variables struct {
	data map[string]string
}

// UnmarshalCSV is a helper method which unmarshals a CSV string
// into a Variables struct with map[string]string data.
func (v *Variables) UnmarshalCSV(csv string) error {
	if len(csv) == 0 {
		return nil
	}

	split := strings.Split(csv, ";")
	data := make(map[string]string, len(split))

	for _, pair := range split {
		key, val, ok := strings.Cut(pair, "=")
		if !ok {
			return fmt.Errorf("KEY=VALUE pair: %q, is of invalid format", pair)
		}
		data[key] = val
	}

	v.data = data

	return nil
}

// MarshalCSV is a helper method which marshals the data
// from a Variables struct into a CSV string.
func (v *Variables) MarshalCSV() (string, error) {
	if v.data == nil {
		return "", nil
	}

	builder := new(strings.Builder)
	l := len(v.data) - 1
	i := 0
	for k, v := range v.data {
		if _, err := builder.WriteString(k + "=" + v); err != nil {
			return "", err
		}
		if i < l {
			if _, err := builder.WriteRune(';'); err != nil {
				return "", err
			}
		}
	}

	return builder.String(), nil
}

// Data is a getter method for the underlying struct data.
func (v *Variables) Data() map[string]string {
	return v.data
}

// CSVData is an empty interface for types which
// support marshalling to/unmarshalling from raw
// CSV data.
type CSVData interface{}

// Sender represents an email sender with the necessary
// information for authentication.
type Sender struct {
	Email    string `csv:"email"`
	Password string `csv:"password"`
	Name     string `csv:"name"`
}

// Receiver represents the recipient of an email message.
type Receiver struct {
	Email     string     `csv:"email"`
	Cc        *List      `csv:"cc"`
	Bcc       *List      `csv:"bcc"`
	Variables *Variables `csv:"variables"`
}

// ReadFile reads CSV data and returns the unmarshalled data
// of type CSVData.
func ReadFile[T CSVData](file string) ([]*T, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	log.Debug().Str("file", file).Msg("reading csv file...")
	var data []*T
	if err := gocsv.UnmarshalFile(f, &data); err != nil {
		return nil, err
	}

	return data, nil
}
