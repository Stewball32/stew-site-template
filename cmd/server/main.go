package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/youruser/yourproject/internal/guards"
	"github.com/youruser/yourproject/internal/pocketbase/hooks"
	"github.com/youruser/yourproject/internal/pocketbase/migrateconf"
	"github.com/youruser/yourproject/internal/pocketbase/oauth"
	"github.com/youruser/yourproject/internal/pocketbase/routes"
	"github.com/youruser/yourproject/internal/pocketbase/seed"
	ws "github.com/youruser/yourproject/internal/websocket"

	discordbot "github.com/youruser/yourproject/internal/disgo"
	"github.com/youruser/yourproject/internal/disgo/commands"
	pb "github.com/youruser/yourproject/internal/pocketbase"
	_ "github.com/youruser/yourproject/internal/websocket/handlers" // self-registering WS handlers
	_ "github.com/youruser/yourproject/internal/websocket/rooms"    // self-registering WS room types
	_ "github.com/youruser/yourproject/migrations"                  // self-registering DB migrations
)

func main() {
	if port := os.Getenv("PUBLIC_PB_PORT"); port != "" && !hasHTTPFlag(os.Args) {
		os.Args = append(os.Args, "--http=0.0.0.0:"+port)
	}

	app := pocketbase.New()

	// Database migrations are the SOURCE OF TRUTH for schema (docs/MIGRATIONS.md).
	// Pending migrations in migrations/ are applied automatically on boot, tracked
	// in the _migrations table. Automigrate writes a migration file whenever the
	// schema changes — enabled in dev builds only, so test/prod can only ever
	// APPLY reviewed migrations, never author new ones from a live admin edit.
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Dir:          "migrations",
		Automigrate:  migrateconf.Automigrate,
		TemplateLang: migratecmd.TemplateLangGo,
	})

	var bot *discordbot.Bot
	var hub *ws.Hub

	// Register record lifecycle hooks (callback registration, fires later).
	hooks.RegisterAll(app)

	// OnServe: register schemas and routes (needs running DB / ServeEvent).
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// NOTE: collections are created by MIGRATIONS (migrations/), applied on
		// boot before this hook — not by on-serve schema code. See docs/MIGRATIONS.md.
		if err := oauth.RegisterAll(app); err != nil {
			return err
		}

		if err := seed.Run(app); err != nil {
			return err
		}

		routes.RegisterAll(se)

		hub = ws.NewHub(app)
		go hub.Run()
		ws.SetInstance(hub)
		se.Router.GET("/api/ws", ws.NewHandler(hub, app))

		// Start Disgo bot (non-blocking). The DEV tier compiles the bot subsystem
		// out entirely (discordbot.SubsystemEnabled == false) — it is the
		// frontend/HMR tier and never opens a gateway. pre (test) is the lowest
		// tier that runs a bot. See docs/DEPLOYMENTS.md.
		if discordbot.SubsystemEnabled {
			b, err := discordbot.NewBot()
			if err != nil {
				log.Printf("Warning: Discord bot not started: %v", err)
			} else if err := b.OpenGateway(context.Background()); err != nil {
				log.Printf("Warning: Discord gateway failed: %v", err)
			} else {
				bot = b
				discordbot.SetInstance(bot)
			}
		} else {
			log.Println("Discord bot: not started (dev tier runs no bot)")
		}

		// Wire up cross-system Services for guards, resolvers, and actions.
		pbSvc := pb.NewService(app)
		svc := &guards.Services{
			App: app,
			WS:  hub,
			PB:  pbSvc,
		}
		if bot != nil {
			svc.Discord = bot
			bot.SetServices(svc)
		}
		hub.SetServices(svc)
		hooks.SetServices(svc)
		commands.SetServices(svc)

		return se.Next()
	})

	// OnTerminate: cleanup.
	app.OnTerminate().BindFunc(func(te *core.TerminateEvent) error {
		if hub != nil {
			hub.Stop()
		}

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

func hasHTTPFlag(args []string) bool {
	for _, a := range args {
		if a == "--http" || strings.HasPrefix(a, "--http=") {
			return true
		}
	}
	return false
}
