package battlefield

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"my/battleship/coordinates"
)

func TestNewService(t *testing.T) {
	log := logrus.New()
	want := &Service{logger: log}
	got := NewService(log)
	assert.Equal(t, want, got)
}

func TestCreateField(t *testing.T) {
	type args struct {
		size  uint
		field Field
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "success, non-zero valid size",
			args:    args{size: maxFieldSize - 1},
			wantErr: nil,
		},
		{
			name:    "error, zero-size",
			args:    args{size: 0},
			wantErr: errorInvalidFieldSize,
		},
		{
			name:    "error, max size exceeded",
			args:    args{size: maxFieldSize + 1},
			wantErr: errorInvalidFieldSize,
		},
		{
			name: "error, field already set",
			args: args{
				size:  maxFieldSize - 1,
				field: Field{isSet: true},
			},
			wantErr: errorFieldAlreadySet,
		},
	}

	for _, tt := range tests {
		s := Service{logger: logrus.New(), f: tt.args.field}

		t.Run(tt.name, func(t *testing.T) {
			err := s.createField(tt.args.size)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestClearField(t *testing.T) {
	type args struct {
		field Field
	}

	tests := []struct {
		name    string
		args    args
		want    Field
		wantErr error
	}{
		{
			name: "success, field initiated",
			args: args{field: Field{
				field: [][]cell{{}},
				size:  1,
				isSet: true,
			}},
			want:    Field{},
			wantErr: nil,
		},
		{
			name:    "success, field not initiated",
			args:    args{field: Field{}},
			want:    Field{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s := Service{logger: logrus.New(), f: tt.args.field}

		t.Run(tt.name, func(t *testing.T) {
			err := s.clearField()
			assert.Equal(t, tt.want, s.f)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPlaceShip(t *testing.T) {
	type args struct {
		ship  *ship
		field Field
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success, occupy all field",
			args: args{
				field: Field{
					field: [][]cell{{{}}},
					size:  2,
				},
				ship: &ship{
					inner: coordinates.Coordinates{{X: 0, Y: 0}: {}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success, occupy not all field",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				ship: &ship{
					inner: coordinates.Coordinates{{X: 0, Y: 0}: {}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success, nearby non-ship occupied cells",
			args: args{
				field: Field{
					field: [][]cell{{{}, {occupied: true}}, {{occupied: true}, {occupied: true}}},
					size:  2,
				},
				ship: &ship{
					inner: coordinates.Coordinates{{X: 0, Y: 0}: {}},
				},
			},
			wantErr: nil,
		},
		{
			name: "error, out of bonds",
			args: args{
				field: Field{
					field: [][]cell{{{}}},
					size:  1,
				},
				ship: &ship{
					inner: coordinates.Coordinates{{X: 0, Y: 0}: {}, {X: 1, Y: 0}: {}, {X: 0, Y: 1}: {}, {X: 1, Y: 1}: {}},
				},
			},
			wantErr: errorOutOfBonds,
		},
		{
			name: "error, place ship on top of other",
			args: args{
				field: Field{
					field: [][]cell{{{
						occupied: true,
						ship:     &ship{},
					}}},
					size: 1,
				},
				ship: &ship{
					inner: coordinates.Coordinates{{X: 0, Y: 0}: {}},
				},
			},
			wantErr: errorCellIsOccupiedByShip,
		},
		{
			name: "error, place ship on occupied cell",
			args: args{
				field: Field{
					field: [][]cell{{{
						occupied: true,
					}}},
					size: 1,
				},
				ship: &ship{
					inner: coordinates.Coordinates{{X: 0, Y: 0}: {}},
				},
			},
			wantErr: errorCellIsOccupiedNearby,
		},
	}

	for _, tt := range tests {
		s := Service{logger: logrus.New(), f: tt.args.field}

		t.Run(tt.name, func(t *testing.T) {
			err := s.placeShip(tt.args.ship)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestAddShips(t *testing.T) {
	type args struct {
		ships []*ship
		field Field
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success, single ship added",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				ships: []*ship{
					{
						inner: coordinates.Coordinates{
							{X: 0, Y: 0}: {},
							{X: 1, Y: 1}: {},
							{X: 1, Y: 0}: {},
							{X: 0, Y: 1}: {},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "success, two non-overlapping ship added",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
					size:  3,
				},
				ships: []*ship{
					{
						inner: coordinates.Coordinates{{X: 0, Y: 0}: {}, {X: 1, Y: 0}: {}, {X: 2, Y: 0}: {}},
					},
					{
						inner: coordinates.Coordinates{{X: 2, Y: 2}: {}},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "error, ships overlapping",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				ships: []*ship{
					{
						inner: coordinates.Coordinates{{X: 0, Y: 0}: {}, {X: 1, Y: 0}: {}},
					},
					{
						inner: coordinates.Coordinates{{X: 0, Y: 0}: {}, {X: 0, Y: 1}: {}},
					},
				},
			},
			wantErr: errorCellIsOccupiedByShip,
		},
		{
			name: "error, ships close to each other",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				ships: []*ship{
					{
						inner: coordinates.Coordinates{{X: 0, Y: 0}: {}},
						outer: coordinates.Coordinates{
							{X: 1, Y: 0}: {},
							{X: 0, Y: 1}: {},
							{X: 1, Y: 1}: {},
						},
					},
					{
						inner: coordinates.Coordinates{{X: 1, Y: 1}: {}, {X: 1, Y: 1}: {}},
					},
				},
			},
			wantErr: errorCellIsOccupiedNearby,
		},
		{
			name: "error, ship is out of bonds",
			args: args{
				field: Field{
					size: 1,
				},
				ships: []*ship{
					{
						inner: coordinates.Coordinates{{X: 2, Y: 2}: {}, {X: 2, Y: 2}: {}},
					},
				},
			},
			wantErr: errorOutOfBonds,
		},
	}

	for _, tt := range tests {
		s := Service{logger: logrus.New(), f: tt.args.field}

		t.Run(tt.name, func(t *testing.T) {
			err := s.addShips(tt.args.ships)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestAddShipsByCoordinates(t *testing.T) {
	type args struct {
		coords string
		field  Field
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success, single ship added",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				coords: "A1 A1",
			},
			wantErr: nil,
		},
		{
			name: "success, tow non-overlapping ship added",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
					size:  3,
				},
				coords: "A1 A1,C3 C3",
			},
			wantErr: nil,
		},
		{
			name: "error, ships overlapping",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				coords: "A1 A2,A1 B1",
			},
			wantErr: errorCellIsOccupiedByShip,
		},
		{
			name: "error, ships close to each other",
			args: args{
				field: Field{
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				coords: "A1 A1,B2 B2",
			},
			wantErr: errorCellIsOccupiedNearby,
		},
		{
			name: "error, ship is out of bonds",
			args: args{
				field: Field{
					size: 1,
				},
				coords: "A5 B5",
			},
			wantErr: errorOutOfBonds,
		},
		{
			name: "error, ships already placed",
			args: args{
				field: Field{
					shipsAdded: true,
				},
				coords: "A5 B5",
			},
			wantErr: errorShipsAlreadyAdded,
		},
		{
			name: "error, invalid coordinates provided",
			args: args{
				field:  Field{},
				coords: "invalid coordinates",
			},
			wantErr: errorInvalidCoordinate,
		},
	}

	for _, tt := range tests {
		s := Service{logger: logrus.New(), f: tt.args.field}

		t.Run(tt.name, func(t *testing.T) {
			err := s.addShipsByCoordinates(tt.args.coords)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestShot(t *testing.T) {
	type args struct {
		coordinate string
		field      Field
	}
	tests := []struct {
		name    string
		args    args
		want    shotResult
		wantErr error
	}{
		{
			name: "success",
			args: args{
				field: Field{
					field:      [][]cell{{{ship: &ship{aliveCells: 2}}, {}}, {{}, {}}},
					shipsAlive: 2,
					size:       2,
					shipsAdded: true,
				},
				coordinate: "A1",
			},
			want: shotResult{
				Destroy: false,
				Knock:   true,
				End:     false,
			},
			wantErr: nil,
		},
		{
			name: "success, kill not-last ship",
			args: args{
				field: Field{
					field:      [][]cell{{{ship: &ship{aliveCells: 1}}, {}}, {{}, {}}},
					shipsAlive: 2,
					size:       2,
					shipsAdded: true,
				},
				coordinate: "A1",
			},
			want: shotResult{
				Destroy: true,
				Knock:   true,
				End:     false,
			},
			wantErr: nil,
		},
		{
			name: "success, kill last ship",
			args: args{
				field: Field{
					field:      [][]cell{{{ship: &ship{aliveCells: 1}}, {}}, {{}, {}}},
					shipsAlive: 1,
					size:       2,
					shipsAdded: true,
				},
				coordinate: "A1",
			},
			want: shotResult{
				Destroy: true,
				Knock:   true,
				End:     true,
			},
			wantErr: nil,
		},
		{
			name: "success, miss",
			args: args{
				field: Field{
					field:      [][]cell{{{ship: &ship{aliveCells: 1}}, {}}, {{}, {}}},
					shipsAlive: 1,
					size:       2,
					shipsAdded: true,
				},
				coordinate: "A2",
			},
			want: shotResult{
				Destroy: false,
				Knock:   false,
				End:     false,
			},
			wantErr: nil,
		},
		{
			name: "error, shot on already shot cell",
			args: args{
				field: Field{
					field:      [][]cell{{{shot: true}, {}}, {{}, {}}},
					size:       2,
					shipsAdded: true,
				},
				coordinate: "A1",
			},
			want:    shotResult{},
			wantErr: errorCellAlreadyShot,
		},
		{
			name: "error, shot out of bonds",
			args: args{
				field: Field{
					field:      [][]cell{{{}, {}}, {{}, {}}},
					size:       2,
					shipsAdded: true,
				},
				coordinate: "C5",
			},
			want:    shotResult{},
			wantErr: errorOutOfBonds,
		},
		{
			name: "error, invalid coordinate",
			args: args{
				field: Field{
					field:      [][]cell{{{}, {}}, {{}, {}}},
					size:       2,
					shipsAdded: true,
				},
				coordinate: "invalid coordinate",
			},
			want:    shotResult{},
			wantErr: errorInvalidCoordinate,
		},
		{
			name: "error, ships not added",
			args: args{
				field: Field{
					field:      [][]cell{{{}, {}}, {{}, {}}},
					size:       2,
					shipsAdded: false,
				},
				coordinate: "A1",
			},
			want:    shotResult{},
			wantErr: errorShipsNotPlaced,
		},
	}

	for _, tt := range tests {
		s := Service{logger: logrus.New(), f: tt.args.field}

		t.Run(tt.name, func(t *testing.T) {
			got, err := s.shot(tt.args.coordinate)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
