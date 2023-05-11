package main

import (
	"context"
	"fmt"
	"os"

	"messenger/services/api/internal/handlers"
	_ "os"
	"os/signal"
	"syscall"

	. "messenger/services/api/pkg/helpers/logger"
)

func main() {
	Logger.Info("Mess running")

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	s := handlers.NewServer()
	s.Init()

	err = s.Run(ctx, fmt.Sprintf("%s/%s", pwd, "server.crt"), fmt.Sprintf("%s/%s", pwd, "server.key"))
	if err != nil {
		Logger.Fatal(err)
	}
}
