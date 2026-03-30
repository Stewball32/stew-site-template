package websocket

// TODO: Central WebSocket Hub managing clients and rooms.
// Hub struct holds client registry, room memberships, and register/unregister channels.
// Run() processes register/unregister/message channels in a loop.
// Broadcast() sends to all clients. SendToUser() targets by user ID. SendToRoom() targets room members.
// JoinRoom() and LeaveRoom() manage room membership for a client.
