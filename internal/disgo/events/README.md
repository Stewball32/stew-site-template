# internal/disgo/events

Discord gateway event listeners.

## Responsibilities

Handles Discord gateway events beyond slash commands — message creates, reactions, guild member joins, presence updates, etc. Event listeners are registered on the Disgo event manager during bot setup.

## Expected Files

- One `.go` file per event domain (e.g., `messages.go`, `members.go`)
- Each exports a `Register*Listeners(client bot.Client)` function

## Example Pattern

```go
// messages.go
func RegisterMessageListeners(client bot.Client) {
    client.EventManager().AddEventListeners(&events.ListenerAdapter{
        OnMessageCreate: func(event *events.MessageCreate) {
            // handle incoming messages
        },
    })
}
```
