package main

import (
	"context"
	"fmt"
	"log"
	"messenger/services/api/internal/handlers"
	_ "os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Printf("Mess running")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	s := handlers.NewServer()
	s.Init()

	err := s.Run(ctx, "", "")
	if err != nil {
		log.Fatal(err)
	}
}
