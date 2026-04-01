package interfaces

import "github.com/disgoorg/snowflake/v2"

// DiscordService abstracts Discord API access for guards and resolvers.
// Implemented by the disgo Bot struct (structural typing, no import needed).
type DiscordService interface {
	IsMember(guildID, userID snowflake.ID) (bool, error)
	MemberRoles(guildID, userID snowflake.ID) ([]snowflake.ID, error)
}
