package routes

import (
	"os"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	register(registerStaticRoute)
}

// registerStaticRoute serves pb_public/ as a catch-all SPA fallback.
// More specific routes registered elsewhere take priority over this wildcard.
func registerStaticRoute(se *core.ServeEvent) {
	se.Router.GET("/{path...}", apis.Static(os.DirFS("pb_public"), true))
}
