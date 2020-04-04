package battlefield

import (
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
	s := Service{logger: logrus.New()}

	type args struct {
		size uint
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.createField(tt.args.size)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
