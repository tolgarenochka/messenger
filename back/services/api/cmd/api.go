package main

import (
	"context"

	"messenger/services/api/internal/handlers"
	_ "os"
	"os/signal"
	"syscall"

	. "messenger/services/api/pkg/helpers/logger"
)

func main() {
	Logger.Info("Mess running")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	s := handlers.NewServer()
	s.Init()

	err := s.Run(ctx, "", "")
	if err != nil {
		Logger.Fatal(err)
	}
}
