//go:build dev

package disgo

// SubsystemEnabled is false in dev builds.
//
// The dev tier runs NO Discord bot — it is the frontend/HMR iteration tier and
// must never open a gateway connection. main.go skips bot startup entirely when
// this is false, so it's an architectural fact of the dev build, not a runtime
// toggle. `pre` (test) is the lowest tier that runs a bot.
// See docs/DEPLOYMENTS.md.
const SubsystemEnabled = false
