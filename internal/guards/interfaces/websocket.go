package interfaces

// WebSocketService abstracts WebSocket Hub state queries for guards and resolvers.
// Implemented by the websocket Hub struct (structural typing, no import needed).
type WebSocketService interface {
	IsConnected(userID string) bool
	IsInRoom(userID string, room string) bool
	UserRooms(userID string) []string
}
