package guards

import (
	"fmt"

	"github.com/disgoorg/disgo/handler"
	"github.com/pocketbase/pocketbase"
)

// RequirePBRole checks that the Discord user has a linked PocketBase account
// with the specified role. Returns the user record or an error.
func RequirePBRole(app *pocketbase.PocketBase, e *handler.CommandEvent, role string) error {
	record, err := RequirePBUser(app, e)
	if err != nil {
		return err
	}

	userRole := record.GetString("role")
	if userRole != role {
		return fmt.Errorf("user %s does not have required role %q", e.User().Username, role)
	}

	return nil
}
