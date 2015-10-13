package main

import (
	"os"
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

const (
	DEBUG = "debug"

)

var log = logrus.New()

func InitLogger(logfile *string, loglevel *string) {
	log.Formatter = new(logrus.TextFormatter)
	log.Hooks.Add(&SourceFileHook{
		LogLevel: logrus.InfoLevel,
	})

	determineLevel(loglevel)
	setOutput(logfile)
}

func determineLevel(loglevel *string) {
	level, err := logrus.ParseLevel(strings.ToLower(*loglevel))
	if err != nil {
		return
	}

	log.Level = level
}

func setOutput(logfile *string) {
	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
		if err != nil {
			log.Panic("Cannot open logfile", err)
			return
		}

		log.Out = f
	}
}

type SourceFileHook struct {
	LogLevel logrus.Level
}

func (hook *SourceFileHook) Fire(entry *logrus.Entry) (_ error) {
	for skip := 4; skip < 9; skip++ {
		_, file, line, _ := runtime.Caller(skip)
		split := strings.Split(file, "/")
		if l := len(split); l > 1 {
			pkg := split[l-2]
			if pkg != "logrus" {
				file = fmt.Sprintf("%s/%s:%d", split[l-2], split[l-1], line)
				entry.Data["src"] = file
				return
			}
		}
	}

	return
}

func (hook *SourceFileHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, hook.LogLevel+1)
	for i, _ := range levels {
		levels[i] = logrus.Level(i)
	}
	return levels
}
