//go:build !dev

package disgo

// SubsystemEnabled is true in test/prod builds (pre, prod). Those tiers start
// the Discord bot when a token is configured. The dev build compiles it out
// (enabled_dev.go). See docs/DEPLOYMENTS.md.
const SubsystemEnabled = true
