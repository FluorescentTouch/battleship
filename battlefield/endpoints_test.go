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
	l := logrus.New()
	e := Endpoints{
		logger:  l,
		service: &Service{logger: l},
	}

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
		t.Run(tt.name, func(t *testing.T) {
			resp, err := e.createFieldEndpoint(tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}
