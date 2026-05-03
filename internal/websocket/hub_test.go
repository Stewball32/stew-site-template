package websocket

import "testing"

func TestNewHubEmptyState(t *testing.T) {
	h := NewHub(nil)

	if h.IsConnected("nobody") {
		t.Error("fresh hub should report no users connected")
	}
	if h.IsInRoom("nobody", "lobby") {
		t.Error("fresh hub should report no users in any room")
	}
	if rooms := h.UserRooms("nobody"); len(rooms) != 0 {
		t.Errorf("fresh hub UserRooms = %v, want empty", rooms)
	}
}
