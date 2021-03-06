package battlefield

import (
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewEndpoints(t *testing.T) {
	l := logrus.New()
	s := &Service{logger: l}
	want := Endpoints{service: s, logger: l}
	got := NewEndpoints(l, s)
	assert.Equal(t, want, got)
}

func TestCreateFieldResponse_StatusCode(t *testing.T) {
	want := http.StatusCreated
	got := CreateFieldResponse{}.StatusCode()
	assert.Equal(t, want, got)
}

func TestCreateFieldEndpoint(t *testing.T) {
	type args struct {
		req CreateFieldRequest
	}

	tests := []struct {
		name    string
		args    args
		want    CreateFieldResponse
		wantErr error
	}{
		{
			name:    "success, non-zero valid size",
			args:    args{req: CreateFieldRequest{Size: maxFieldSize - 1}},
			want:    CreateFieldResponse{},
			wantErr: nil,
		},
		{
			name:    "error, zero-size",
			args:    args{req: CreateFieldRequest{Size: 0}},
			want:    CreateFieldResponse{},
			wantErr: errorInvalidFieldSize,
		},
		{
			name:    "error, max size exceeded",
			args:    args{req: CreateFieldRequest{Size: maxFieldSize + 1}},
			want:    CreateFieldResponse{},
			wantErr: errorInvalidFieldSize,
		},
	}

	for _, tt := range tests {
		l := logrus.New()
		e := Endpoints{
			logger:  l,
			service: &Service{logger: l},
		}

		t.Run(tt.name, func(t *testing.T) {
			resp, err := e.createFieldEndpoint(tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestClearFieldResponse_StatusCode(t *testing.T) {
	want := http.StatusOK
	got := ClearFieldResponse{}.StatusCode()
	assert.Equal(t, want, got)
}

func TestClearFieldEndpoint(t *testing.T) {
	type args struct {
		f Field
	}

	tests := []struct {
		name    string
		args    args
		want    ClearFieldResponse
		wantErr error
	}{
		{
			name:    "success, field is set",
			args:    args{f: Field{isSet: true}},
			want:    ClearFieldResponse{},
			wantErr: nil,
		},
		{
			name:    "success, field is not set",
			args:    args{},
			want:    ClearFieldResponse{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		l := logrus.New()
		e := Endpoints{
			logger:  l,
			service: &Service{logger: l},
		}

		t.Run(tt.name, func(t *testing.T) {
			resp, err := e.clearFieldEndpoint()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestAddShipsResponse_StatusCode(t *testing.T) {
	want := http.StatusCreated
	got := AddShipsResponse{}.StatusCode()
	assert.Equal(t, want, got)
}

func TestAddShipsEndpoint(t *testing.T) {
	type args struct {
		field Field
		req   AddShipsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    AddShipsResponse
		wantErr error
	}{
		{
			name: "success",
			args: args{
				field: Field{
					field: [][]cell{{{}}},
					size:  1,
				},
				req: AddShipsRequest{Coords: "A1 A1"},
			},
			want:    AddShipsResponse{},
			wantErr: nil,
		},
		{
			name: "error",
			args: args{
				field: Field{
					shipsAdded: true,
				},
				req: AddShipsRequest{Coords: "A1 A1"},
			},
			want:    AddShipsResponse{},
			wantErr: errorShipsAlreadyAdded,
		},
	}

	for _, tt := range tests {
		l := logrus.New()
		e := Endpoints{
			logger: l,
			service: &Service{
				f:      tt.args.field,
				logger: l,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			resp, err := e.addShipsEndpoint(tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestShotsResponse_StatusCode(t *testing.T) {
	want := http.StatusOK
	got := ShotResponse{}.StatusCode()
	assert.Equal(t, want, got)
}

func TestShotEndpoint(t *testing.T) {
	type args struct {
		field Field
		req   ShotRequest
	}
	tests := []struct {
		name    string
		args    args
		want    ShotResponse
		wantErr error
	}{
		{
			name: "success",
			args: args{
				field: Field{
					field:      [][]cell{{{ship: &ship{aliveCells: 1}}}},
					size:       1,
					shipsAlive: 1,
					shipsAdded: true,
				},
				req: ShotRequest{Coord: "A1"},
			},
			want: ShotResponse{
				Knock:   true,
				Destroy: true,
				End:     true,
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: args{
				field: Field{
					shipsAdded: false,
				},
				req: ShotRequest{Coord: "A1"},
			},
			want:    ShotResponse{},
			wantErr: errorShipsNotPlaced,
		},
	}

	for _, tt := range tests {
		l := logrus.New()
		e := Endpoints{
			logger: l,
			service: &Service{
				f:      tt.args.field,
				logger: l,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			resp, err := e.shotEndpoint(tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestStateResponse_StatusCode(t *testing.T) {
	want := http.StatusOK
	got := StateResponse{}.StatusCode()
	assert.Equal(t, want, got)
}

func TestStateEndpoint(t *testing.T) {
	type args struct {
		f Field
	}

	tests := []struct {
		name string
		args args
		want StateResponse
	}{
		{
			name: "success",
			args: args{f: Field{
				state: state{
					knocked:   1,
					destroyed: 1,
					shipCount: 3,
					shotCount: 5,
				},
			},
			},
			want: StateResponse{
				Knocked:   1,
				Destroyed: 1,
				ShipCount: 3,
				ShotCount: 5,
			},
		},
		{
			name: "success, game is not started",
			args: args{f: Field{}},
			want: StateResponse{
				Knocked:   0,
				Destroyed: 0,
				ShipCount: 0,
				ShotCount: 0,
			},
		},
	}

	for _, tt := range tests {
		l := logrus.New()
		e := Endpoints{
			logger:  l,
			service: &Service{logger: l, f: tt.args.f},
		}

		t.Run(tt.name, func(t *testing.T) {
			resp := e.stateEndpoint()
			assert.Equal(t, tt.want, resp)
		})
	}
}
