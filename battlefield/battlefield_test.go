package battlefield

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewField(t *testing.T) {
	type args struct {
		size uint
	}

	tests := []struct {
		name string
		args args
		want Field
	}{
		{
			name: "success,zero-size",
			args: args{size: 0},
			want: Field{
				field: [][]cell{},
				size:  0,
			},
		},
		{
			name: "success, non-zero size",
			args: args{size: 2},
			want: Field{
				field: [][]cell{{{}, {}}, {{}, {}}},
				size:  2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewField(tt.args.size)
			assert.Equal(t, tt.want, got)
		})
	}
}
