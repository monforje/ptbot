package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"ptbot/internal/app"
	"syscall"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigs
		log.Println("received stop signal")
		a.Stop()
		cancel()
	}()

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
