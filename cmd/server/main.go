package main

import (
	"context"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/youruser/yourproject/internal/pocketbase/hooks"
	"github.com/youruser/yourproject/internal/pocketbase/oauth"
	"github.com/youruser/yourproject/internal/pocketbase/routes"
	"github.com/youruser/yourproject/internal/pocketbase/schema"

	discordbot "github.com/youruser/yourproject/internal/disgo"
)

func main() {
	app := pocketbase.New()

	var bot *discordbot.Bot

	// Register record lifecycle hooks (callback registration, fires later).
	hooks.RegisterAll(app)

	// OnServe: register schemas and routes (needs running DB / ServeEvent).
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if err := schema.RegisterAll(app); err != nil {
			return err
		}

		if err := oauth.RegisterAll(app); err != nil {
			return err
		}

		routes.RegisterAll(se)

		// TODO: Initialize WebSocket Hub
		// hub := websocket.NewHub()
		// go hub.Run()

		// Start Disgo bot (non-blocking)
		var err error
		bot, err = discordbot.NewBot()
		if err != nil {
			log.Printf("Warning: Discord bot not started: %v", err)
		} else {
			if err := bot.OpenGateway(context.Background()); err != nil {
				log.Printf("Warning: Discord gateway failed: %v", err)
				bot = nil
			} else {
				discordbot.SetInstance(bot)
			}
		}

		return se.Next()
	})

	// OnTerminate: cleanup.
	app.OnTerminate().BindFunc(func(te *core.TerminateEvent) error {
		if bot != nil {
			bot.Close(context.Background())
		}

		log.Println("Server shutting down...")
		return te.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
