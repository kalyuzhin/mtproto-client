package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/examples"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/telegram/query/dialogs"
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

	rtc := config.MustLoad()

	client := telegram.NewClient(
		rtc.GetAppID(),
		rtc.GetAppHash(),
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: rtc.GetSessionFile()},
			Middlewares: []telegram.Middleware{
				floodwait.NewSimpleWaiter(),
				ratelimit.New(20, 1),
			},
		})
	err := client.Run(ctx, func(ctx context.Context) error {
		errRun := authCode(ctx, client, rtc)
		if errRun != nil {
			log.Fatal(errRun)
		}
		errRun = getChats(ctx, client)
		if errRun != nil {
			return errRun
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func authQR(ctx context.Context, client *telegram.Client, password string) error {
	status, errRun := client.Auth().Status(ctx)
	if errRun != nil {
		return errRun
	}

	if status.Authorized {
		log.Println("Reusing stored session, no login required")
	} else {
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
	}

	return nil
}

func authCode(ctx context.Context, client *telegram.Client, rtConfig rtConfig) error {
	authFlow := auth.NewFlow(auth.Constant(rtConfig.GetPhoneNumber(),
		rtConfig.GetPassword(), examples.Terminal{}), auth.SendCodeOptions{})

	status, err := client.Auth().Status(ctx)
	if err != nil {
		return err
	}

	if status.Authorized {
		log.Println("Reusing stored session, no login required")
	} else {
		log.Println("No valid session, logging in with code")
		if err = client.Auth().IfNecessary(ctx, authFlow); err != nil {
			return err
		}
	}

	return nil
}

func getChats(ctx context.Context, client *telegram.Client) error {
	api := client.API()
	count := 0
	return query.GetDialogs(api).ForEach(ctx,
		func(ctx context.Context, elem dialogs.Elem) error {
			if elem.Deleted() {
				return nil
			}
			count++
			log.Printf("№%d – %s", count, dialogName(elem))

			return nil
		})
}

func dialogName(elem dialogs.Elem) string {
	peer, ok := elem.Dialog.(*tg.Dialog)
	if !ok {
		return "?"
	}
	switch p := peer.Peer.(type) {
	case *tg.PeerUser:
		if u, ok := elem.Entities.User(p.UserID); ok {
			return fmt.Sprintf("Name:%s \nLastName: %s \nUsername: %s\n", u.FirstName, u.LastName, u.Username)
		}
	case *tg.PeerChat:
		if c, ok := elem.Entities.Chat(p.ChatID); ok {
			return c.Title
		}
	case *tg.PeerChannel:
		if c, ok := elem.Entities.Channel(p.ChannelID); ok {
			return c.Title
		}
	}
	return "?"
}
