package websocket

// TODO: Single WebSocket connection with read/write pumps.
// Client struct holds the websocket.Conn, Hub reference, send channel, and optional user record.
// readPump() reads messages from the browser and forwards them to the Hub for routing.
// writePump() sends queued messages from the Hub to the browser.
