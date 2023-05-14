package main

import (
	"context"
	"encoding/json"
	"fmt"
	"messenger/services/api/pkg/helpers/models"
	"os"

	"messenger/services/api/internal/handlers"
	_ "os"
	"os/signal"
	"syscall"

	. "messenger/services/api/pkg/helpers/logger"
)

var Configuration models.Configuration

func init() {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, _ := os.Open(pwd + "/config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	Configuration.PathToServerCrt = pwd + Configuration.PathToServerCrt
	Configuration.PathToServerKey = pwd + Configuration.PathToServerKey
}

func main() {
	Logger.Info("Mess running")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	s := handlers.NewServer(Configuration)
	s.Init()

	err := s.Run(ctx)
	if err != nil {
		Logger.Fatal(err)
	}
}
