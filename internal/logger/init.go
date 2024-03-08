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

package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02/01/06 03:04PM",
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

// Init sets the log level and initialises the global logger.
func Init(level zerolog.Level) error {
	if level < zerolog.TraceLevel && level > zerolog.NoLevel {
		return fmt.Errorf("expected values between -1 and 6, got: %v", level)
	}
	zerolog.SetGlobalLevel(level)

	return nil
}
