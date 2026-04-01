package guards

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/youruser/yourproject/internal/guards/interfaces"
)

// Services bundles all system access a guard or resolver may need.
// Fields may be nil if the corresponding system is not running.
type Services struct {
	App     core.App
	Discord interfaces.DiscordService
	WS      interfaces.WebSocketService
}
