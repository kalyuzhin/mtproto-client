package main

import (
	"context"
	"log"

	"github.com/gotd/td/examples"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/kalyuzhin/mtproto-client/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rtConfig := config.MustLoad()

	authFlow := auth.NewFlow(auth.Constant(rtConfig.GetPhoneNumber(),
		rtConfig.GetPassword(), examples.Terminal{}), auth.SendCodeOptions{})

	client := telegram.NewClient(
		rtConfig.GetAppID(),
		rtConfig.GetAppHash(),
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: rtConfig.GetSessionFile()},
		})
	err := client.Run(ctx, func(ctx context.Context) error {
		status, errRun := client.Auth().Status(ctx)
		if errRun != nil {
			return errRun
		}

		if status.Authorized {
			log.Println("Reusing stored session, no login required")
		} else {
			log.Println("No valid session, logging in with code")
			if errRun = client.Auth().IfNecessary(ctx, authFlow); errRun != nil {
				return errRun
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
