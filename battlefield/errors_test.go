package battlefield

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    HTTPError
		want string
	}{
		{
			name: "errorInvalidInputParams",
			e:    errorInvalidInputParams,
			want: "invalid input params",
		},
		{
			name: "errorInvalidFieldSize",
			e:    errorInvalidFieldSize,
			want: "field size is invalid",
		},
		{
			name: "errorFieldAlreadySet",
			e:    errorFieldAlreadySet,
			want: "field is already set",
		},
		{
			name: "errorInvalidCoordinate",
			e:    errorInvalidCoordinate,
			want: "invalid coordinate provided",
		},
		{
			name: "errorCellIsOccupiedByShip",
			e:    errorCellIsOccupiedByShip,
			want: "can't place ships on top of each other",
		},
		{
			name: "errorCellIsOccupiedNearby",
			e:    errorCellIsOccupiedNearby,
			want: "can't place ships close to each other",
		},
		{
			name: "errorShipsAlreadyAdded",
			e:    errorShipsAlreadyAdded,
			want: "ships are already added",
		},
		{
			name: "errorOutOfBonds",
			e:    errorOutOfBonds,
			want: "out of bonds",
		},
		{
			name: "errorCellAlreadyShot",
			e:    errorCellAlreadyShot,
			want: "cell was already shot",
		},
		{
			name: "errorShipsNotPlaced",
			e:    errorShipsNotPlaced,
			want: "ships not placed yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHTTPError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		e    HTTPError
		want int
	}{
		{
			name: "errorInvalidInputParams",
			e:    errorInvalidInputParams,
			want: http.StatusBadRequest,
		},
		{
			name: "errorInvalidFieldSize",
			e:    errorInvalidFieldSize,
			want: http.StatusBadRequest,
		},
		{
			name: "errorFieldAlreadySet",
			e:    errorFieldAlreadySet,
			want: http.StatusConflict,
		},
		{
			name: "errorInvalidCoordinate",
			e:    errorInvalidCoordinate,
			want: http.StatusBadRequest,
		},
		{
			name: "errorCellIsOccupiedByShip",
			e:    errorCellIsOccupiedByShip,
			want: http.StatusBadRequest,
		},
		{
			name: "errorCellIsOccupiedNearby",
			e:    errorCellIsOccupiedNearby,
			want: http.StatusBadRequest,
		},
		{
			name: "errorShipsAlreadyAdded",
			e:    errorShipsAlreadyAdded,
			want: http.StatusBadRequest,
		},
		{
			name: "errorOutOfBonds",
			e:    errorOutOfBonds,
			want: http.StatusBadRequest,
		},
		{
			name: "errorCellAlreadyShot",
			e:    errorCellAlreadyShot,
			want: http.StatusBadRequest,
		},
		{
			name: "errorShipsNotPlaced",
			e:    errorShipsNotPlaced,
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.StatusCode()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHTTPError_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		e       HTTPError
		want    string
		wantErr error
	}{
		{
			name:    "errorInvalidInputParams",
			e:       errorInvalidInputParams,
			want:    `{"err":"invalid input params"}`,
			wantErr: nil,
		},
		{
			name:    "errorInvalidFieldSize",
			e:       errorInvalidFieldSize,
			want:    `{"err":"field size is invalid"}`,
			wantErr: nil,
		},
		{
			name:    "errorFieldAlreadySet",
			e:       errorFieldAlreadySet,
			want:    `{"err":"field is already set"}`,
			wantErr: nil,
		},
		{
			name:    "errorInvalidCoordinate",
			e:       errorInvalidCoordinate,
			want:    `{"err":"invalid coordinate provided"}`,
			wantErr: nil,
		},
		{
			name:    "errorCellIsOccupiedByShip",
			e:       errorCellIsOccupiedByShip,
			want:    `{"err":"can't place ships on top of each other"}`,
			wantErr: nil,
		},
		{
			name:    "errorCellIsOccupiedNearby",
			e:       errorCellIsOccupiedNearby,
			want:    `{"err":"can't place ships close to each other"}`,
			wantErr: nil,
		},
		{
			name:    "errorShipsAlreadyAdded",
			e:       errorShipsAlreadyAdded,
			want:    `{"err":"ships are already added"}`,
			wantErr: nil,
		},
		{
			name:    "errorOutOfBonds",
			e:       errorOutOfBonds,
			want:    `{"err":"out of bonds"}`,
			wantErr: nil,
		},
		{
			name:    "errorCellAlreadyShot",
			e:       errorCellAlreadyShot,
			want:    `{"err":"cell was already shot"}`,
			wantErr: nil,
		},
		{
			name:    "errorShipsNotPlaced",
			e:       errorShipsNotPlaced,
			want:    `{"err":"ships not placed yet"}`,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.MarshalJSON()
			assert.JSONEq(t, tt.want, string(got))
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
