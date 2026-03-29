# internal/disgo

Discord bot implementation using [Disgo](https://github.com/disgoorg/disgo).

## Responsibilities

Creates and manages the Disgo bot client. Registers slash commands and attaches event listeners. The bot runs in its own goroutine alongside PocketBase.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `commands/` | Slash command definitions and interaction handlers |
| `events/` | Discord gateway event listeners (messages, reactions, etc.) |

## Expected Files

- `bot.go` — constructs the Disgo client, registers commands, attaches event handlers, and exposes a `Start(ctx context.Context) error` function
- `config.go` (optional) — bot-specific config (token, guild IDs, etc.)

## Bot Lifecycle

The bot is started from `cmd/server/main.go` in a goroutine before `app.Start()`:

```go
go func() {
    if err := disgobot.Start(ctx); err != nil {
        log.Fatal(err)
    }
}()
```
