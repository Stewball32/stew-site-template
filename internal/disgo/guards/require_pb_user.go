package guards

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	pbactions "github.com/youruser/yourproject/internal/pocketbase/actions"
)

// RequirePBUser checks that the Discord user has a linked PocketBase account.
// Returns the PocketBase user record or an error.
func RequirePBUser(app *pocketbase.PocketBase, e *handler.CommandEvent) (*core.Record, error) {
	return pbactions.FindUserByDiscordID(app, e.User().ID.String())
}
