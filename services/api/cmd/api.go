package main

import (
	"context"
	"fmt"
	"log"
	"messenger/services/api/internal/server"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Printf("Mess running")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s := server.NewServer()
	s.Init()

	err := s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
