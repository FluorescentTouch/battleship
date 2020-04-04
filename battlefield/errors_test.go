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
			name: "errorInvalidCreateFieldParams",
			e:    errorInvalidCreateFieldParams,
			want: "invalid input params",
		},
		{
			name: "errorInvalidFieldSize",
			e:    errorInvalidFieldSize,
			want: "Field size is invalid",
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
			name: "errorInvalidCreateFieldParams",
			e:    errorInvalidCreateFieldParams,
			want: http.StatusBadRequest,
		},
		{
			name: "errorInvalidFieldSize",
			e:    errorInvalidFieldSize,
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
			name:    "errorInvalidCreateFieldParams",
			e:       errorInvalidCreateFieldParams,
			want:    `{"err":"invalid input params"}`,
			wantErr: nil,
		},
		{
			name:    "errorInvalidFieldSize",
			e:       errorInvalidFieldSize,
			want:    `{"err":"Field size is invalid"}`,
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