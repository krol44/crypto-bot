package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type Struct struct {
	IsDev    bool
	BotToken string
}

var config Struct

func init() {
	dev := os.Getenv("dev") == "true"
	botToken := os.Getenv("bot_token")

	config = Struct{
		dev,
		botToken,
	}

	logSetup()
}

func logSetup() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
	if l, err := log.ParseLevel("debug"); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
		log.SetOutput(os.Stdout)
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
