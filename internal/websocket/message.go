package websocket

// TODO: Wire format for WebSocket messages.
// Message struct with Type, Room, Target (user ID), and Payload (json.RawMessage) fields.
// Types: "broadcast", "room", "direct", "join_room", "leave_room".
// Hub inspects Type to decide routing; Payload is opaque project-specific data.
