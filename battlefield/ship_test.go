package battlefield

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"my/battleship/coordinates"
)

func TestNewShip(t *testing.T) {
	c := [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 1, Y: 1}}
	want := &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 1, Y: 1}}}
	got := newShip(c[0], c[1])
	assert.Equal(t, want, got)
}

func TestMakeShipsFromCoords(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []*ship
		wantErr error
	}{
		{
			name: "success, single ship",
			args: "A1 A1",
			want: []*ship{{
				c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}},
			}},
			wantErr: nil,
		},
		{
			name: "success, multiple ships",
			args: "A1 A1,B3 C4",
			want: []*ship{
				{
					c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}},
				},
				{
					c: [2]coordinates.Coordinate{{X: 1, Y: 2}, {X: 2, Y: 3}},
				},
			},
			wantErr: nil,
		},
		{
			name:    "error, zero ships",
			args:    "",
			want:    nil,
			wantErr: errorInvalidCoordinate,
		},
		{
			name:    "error, single coordinate provided",
			args:    "A1",
			want:    nil,
			wantErr: errorInvalidCoordinate,
		},
		{
			name:    "error, first coordinate is invalid",
			args:    "A B2",
			want:    nil,
			wantErr: errorInvalidCoordinate,
		},
		{
			name:    "error, second coordinate is invalid",
			args:    "A1 2",
			want:    nil,
			wantErr: errorInvalidCoordinate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ship, err := makeShipsFromCoords(tt.args)
			assert.Equal(t, tt.want, ship)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
