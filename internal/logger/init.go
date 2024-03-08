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
