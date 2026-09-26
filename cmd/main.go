package main

import (
	"context"
	"log"
	"os"

	"github.com/gotd/td/examples"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/kalyuzhin/mtproto-client/internal/config"
	"github.com/mdp/qrterminal/v3"
)

type rtConfig interface {
	GetPhoneNumber() string
	GetPassword() string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rtConfig := config.MustLoad()

	client := telegram.NewClient(
		rtConfig.GetAppID(),
		rtConfig.GetAppHash(),
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: rtConfig.GetSessionFile()},
		})
	err := client.Run(ctx, func(ctx context.Context) error {
		errRun := authQR(ctx, client, rtConfig.GetPassword())
		if errRun != nil {
			log.Fatal(errRun)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func authQR(ctx context.Context, client *telegram.Client, password string) error {
	dispatcher := tg.NewUpdateDispatcher()
	loggedIn := qrlogin.OnLoginToken(&dispatcher)

	show := func(ctx context.Context, token qrlogin.Token) error {
		qrterminal.Generate(token.URL(), qrterminal.L, os.Stderr)
		return nil
	}

	if _, err := client.QR().Auth(ctx, loggedIn, show); err != nil {
		if !tgerr.Is(err, "SESSION_PASSWORD_NEEDED") {
			return err
		}
		// Prompt for the 2FA password and finish.
		if _, err := client.Auth().Password(ctx, password); err != nil {
			return err
		}
	}

	return nil
}

func authCode(ctx context.Context, client *telegram.Client, rtConfig rtConfig) error {
	authFlow := auth.NewFlow(auth.Constant(rtConfig.GetPhoneNumber(),
		rtConfig.GetPassword(), examples.Terminal{}), auth.SendCodeOptions{})

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
}
