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
