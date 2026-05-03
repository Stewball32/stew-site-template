package guards

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func pass(_ *Services, _ *core.Record) error  { return nil }
func fail(_ *Services, _ *core.Record) error  { return errors.New("denied") }
func fail2(_ *Services, _ *core.Record) error { return errors.New("denied 2") }

func TestAny(t *testing.T) {
	tests := []struct {
		name    string
		guards  []GuardFunc
		wantErr bool
	}{
		{"empty passes nil", nil, false},
		{"first passes", []GuardFunc{pass, fail}, false},
		{"second passes", []GuardFunc{fail, pass}, false},
		{"all fail returns last error", []GuardFunc{fail, fail2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Any(tt.guards...)(nil, nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name    string
		guards  []GuardFunc
		wantErr bool
	}{
		{"empty passes", nil, false},
		{"all pass", []GuardFunc{pass, pass}, false},
		{"first fails short-circuits", []GuardFunc{fail, pass}, true},
		{"last fails", []GuardFunc{pass, fail}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := All(tt.guards...)(nil, nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}
