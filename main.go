package main

import (
	"context"
	"os"
	"os/signal"

	"telegram/config"
	"telegram/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	config.Validate()

	level, err := log.ParseLevel(config.Args.LOG_LEVEL)
	if err != nil {
		log.Fatal(err)
	}

	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(level)

	go utils.StartTelegramBot(ctx)

	<-c
	cancel()
}
