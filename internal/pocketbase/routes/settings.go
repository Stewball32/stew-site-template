package routes

import (
	"os"

	"github.com/pocketbase/pocketbase/core"
)

func init() {
	register(registerSettingsRoutes)
}

func registerSettingsRoutes(se *core.ServeEvent) {
	serve := func(e *core.RequestEvent) error {
		return e.FileFS(os.DirFS("pb_public"), "settings/index.html")
	}

	se.Router.GET("/settings", serve)
	se.Router.GET("/settings/", serve)
}
