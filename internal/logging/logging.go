package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

var (
	Logger              zerolog.Logger
	EnpassVaultLogLevel = logrus.Level(zerolog.ErrorLevel)
)

func Init() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}
