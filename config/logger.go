package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	format string
	level  string
}

func (l LoggerConfig) Setup() {
	level, err := logrus.ParseLevel(l.level)
	if err != nil {
		logrus.Fatalln("Unable to parse log level", err.Error())
	}

	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(level)

	switch strings.ToLower(l.format) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
}

func SetupTestLogger() {
	logrus.SetOutput(ioutil.Discard)
}
