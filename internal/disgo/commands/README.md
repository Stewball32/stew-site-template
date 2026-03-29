# internal/disgo/commands

Discord application command (slash command) definitions and handlers.

## Responsibilities

Each file defines one or more slash commands and their interaction handler logic. Commands are collected and registered with the Discord API at bot startup.

## Expected Files

- One `.go` file per command or command group (e.g., `ping.go`, `info.go`)
- A `commands.go` that exports the full list of `discord.ApplicationCommandCreate` definitions and a map of handler functions

## Example Pattern

```go
// ping.go
var PingCommand = discord.SlashCommandCreate{
    Name:        "ping",
    Description: "Replies with pong",
}

func HandlePing(event *events.ApplicationCommandInteractionCreate) error {
    return event.CreateMessage(discord.MessageCreate{Content: "pong!"})
}
```
