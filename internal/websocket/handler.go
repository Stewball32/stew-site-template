package websocket

// TODO: HTTP upgrade endpoint for WebSocket connections.
// NewHandler(hub, app) returns an http.Handler that upgrades requests via websocket.Accept().
// Checks optional ?token= query param, validates JWT with app.FindAuthRecordByToken().
// Valid token tags the Client with the user record; missing/invalid token allows anonymous access.
