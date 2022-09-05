package main

import (
	"context"
	"os"
	"os/signal"

	"telegram/internal/controller"
)

// create struct for json bodies
// Interface for GetOriginString
// Remove explicit status checks in handlers

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go controller.StartTelegramBot(ctx)

	<-c
	cancel()
}
