package main

import (
	"context"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/kalyuzhin/mtproto-client/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rtConfig := config.MustLoad()

	client := telegram.NewClient(rtConfig.GetAppID(), rtConfig.GetAppHash(), telegram.Options{})
	err := client.Run(ctx, func(ctx context.Context) error {
		client.Auth()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
