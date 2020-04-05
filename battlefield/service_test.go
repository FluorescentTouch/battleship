package battlefield

import (
	"my/battleship/coordinates"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
					field: [][]cell{{{}, {}}, {{}, {}}},
					size:  2,
				},
				ship: &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 1, Y: 1}}},
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
				ship: &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}}},
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
				ship: &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}}},
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
				ship: &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 1, Y: 1}}},
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
				ship: &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}}},
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
				ship: &ship{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}}},
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
				ships: []*ship{{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 1, Y: 1}}}},
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
				ships: []*ship{
					{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 2, Y: 0}}},
					{c: [2]coordinates.Coordinate{{X: 2, Y: 2}, {X: 2, Y: 2}}},
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
					{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 1, Y: 0}}},
					{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 1}}},
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
					{c: [2]coordinates.Coordinate{{X: 0, Y: 0}, {X: 0, Y: 0}}},
					{c: [2]coordinates.Coordinate{{X: 1, Y: 1}, {X: 1, Y: 1}}},
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
					{c: [2]coordinates.Coordinate{{X: 2, Y: 2}, {X: 2, Y: 2}}},
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
