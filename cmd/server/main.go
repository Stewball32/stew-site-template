package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/youruser/yourproject/internal/pocketbase/hooks"
	"github.com/youruser/yourproject/internal/pocketbase/oauth"
	"github.com/youruser/yourproject/internal/pocketbase/routes"
	"github.com/youruser/yourproject/internal/pocketbase/schema"
)

func main() {
	app := pocketbase.New()

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

		// TODO: Start Disgo bot (non-blocking)
		// bot, err := disgo.NewBot()
		// if err != nil { return err }
		// if err := bot.OpenGateway(context.Background()); err != nil { return err }

		return se.Next()
	})

	// OnTerminate: cleanup.
	app.OnTerminate().BindFunc(func(te *core.TerminateEvent) error {
		// TODO: Shut down Disgo bot
		// bot.Close(context.Background())

		log.Println("Server shutting down...")
		return te.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
