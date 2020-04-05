package coordinates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCoordinate(t *testing.T) {
	tests := []struct {
		name   string
		args   string
		want   Coordinate
		wantOK bool
	}{
		{
			name: "success, digit is less than 10",
			args: "B5",
			want: Coordinate{
				X: 1,
				Y: 4,
			},
			wantOK: true,
		},
		{
			name: "success, digit is bigger than 10",
			args: "B50",
			want: Coordinate{
				X: 1,
				Y: 49,
			},
			wantOK: true,
		},
		{
			name:   "error, single letter input",
			args:   "B",
			wantOK: false,
		},
		{
			name:   "error, single number input",
			args:   "1",
			wantOK: false,
		},
		{
			name:   "error, non-alphabetic letter",
			args:   "*50",
			wantOK: false,
		},
		{
			name:   "error, non-numeric second symbol",
			args:   "AZ",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ConvertCoordinate(tt.args)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOK, ok)
		})
	}
}

func TestGetInnerOuterCells(t *testing.T) {
	tests := []struct {
		name      string
		args      [2]Coordinate
		wantInner map[Coordinate]struct{}
		wantOuter map[Coordinate]struct{}
	}{
		{
			name:      "not zero",
			args:      [2]Coordinate{{1, 1}, {1, 1}},
			wantInner: map[Coordinate]struct{}{{1, 1}: {}},
			wantOuter: map[Coordinate]struct{}{{0, 0}: {}, {0, 1}: {}, {0, 2}: {}, {1, 2}: {}, {2, 2}: {}, {2, 1}: {}, {2, 0}: {}, {1, 0}: {}},
		},
		{
			name:      "x is zero",
			args:      [2]Coordinate{{0, 1}, {0, 1}},
			wantInner: map[Coordinate]struct{}{{0, 1}: {}},
			wantOuter: map[Coordinate]struct{}{{0, 0}: {}, {1, 0}: {}, {1, 1}: {}, {1, 2}: {}, {0, 2}: {}},
		},
		{
			name:      "y is zero",
			args:      [2]Coordinate{{1, 0}, {1, 0}},
			wantInner: map[Coordinate]struct{}{{1, 0}: {}},
			wantOuter: map[Coordinate]struct{}{{0, 0}: {}, {0, 1}: {}, {1, 1}: {}, {2, 1}: {}, {2, 0}: {}},
		},
		{
			name:      "both coordinates are zero",
			args:      [2]Coordinate{{0, 0}, {0, 0}},
			wantInner: map[Coordinate]struct{}{{0, 0}: {}},
			wantOuter: map[Coordinate]struct{}{{0, 1}: {}, {1, 1}: {}, {1, 0}: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, out := GetInnerOuterCells(tt.args)
			assert.Equal(t, tt.wantInner, in)
			assert.Equal(t, tt.wantOuter, out)
		})
	}
}

func TestOuterCells(t *testing.T) {
	type args struct {
		lx, bx, ly, by uint
	}

	tests := []struct {
		name string
		args args
		want map[Coordinate]struct{}
	}{
		{
			name: "not zero",
			args: args{
				lx: 1,
				bx: 1,
				ly: 1,
				by: 1,
			},
			want: map[Coordinate]struct{}{{0, 0}: {}, {0, 1}: {}, {0, 2}: {}, {1, 2}: {}, {2, 2}: {}, {2, 1}: {}, {2, 0}: {}, {1, 0}: {}},
		},
		{
			name: "x is zero",
			args: args{
				lx: 0,
				bx: 0,
				ly: 1,
				by: 1,
			},
			want: map[Coordinate]struct{}{{0, 0}: {}, {1, 0}: {}, {1, 1}: {}, {1, 2}: {}, {0, 2}: {}},
		},
		{
			name: "y is zero",
			args: args{
				lx: 1,
				bx: 1,
				ly: 0,
				by: 0,
			},
			want: map[Coordinate]struct{}{{0, 0}: {}, {0, 1}: {}, {1, 1}: {}, {2, 1}: {}, {2, 0}: {}},
		},
		{
			name: "both coordinates are zero",
			args: args{
				lx: 0,
				bx: 0,
				ly: 0,
				by: 0,
			},
			want: map[Coordinate]struct{}{{0, 1}: {}, {1, 1}: {}, {1, 0}: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := outerCells(
				tt.args.lx,
				tt.args.bx,
				tt.args.ly,
				tt.args.by,
			)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInnerCells(t *testing.T) {
	type args struct {
		lx, bx, ly, by uint
	}

	tests := []struct {
		name string
		args args
		want map[Coordinate]struct{}
	}{
		{
			name: "single cell",
			args: args{
				lx: 1,
				bx: 1,
				ly: 1,
				by: 1,
			},
			want: map[Coordinate]struct{}{{1, 1}: {}},
		},
		{
			name: "multiple cells, left to right",
			args: args{
				lx: 1,
				bx: 3,
				ly: 1,
				by: 1,
			},
			want: map[Coordinate]struct{}{{1, 1}: {}, {2, 1}: {}, {3, 1}: {}},
		},
		{
			name: "multiple cells, up to bottom",
			args: args{
				lx: 1,
				bx: 1,
				ly: 1,
				by: 3,
			},
			want: map[Coordinate]struct{}{{1, 1}: {}, {1, 2}: {}, {1, 3}: {}},
		},
		{
			name: "multiple cells, diagonal",
			args: args{
				lx: 1,
				bx: 2,
				ly: 1,
				by: 2,
			},
			want: map[Coordinate]struct{}{{1, 1}: {}, {1, 2}: {}, {2, 1}: {}, {2, 2}: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := innerCells(
				tt.args.lx,
				tt.args.bx,
				tt.args.ly,
				tt.args.by,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSortY(t *testing.T) {
	tests := []struct {
		name  string
		args  [2]Coordinate
		wantB uint
		wantL uint
	}{
		{
			name:  "single cell",
			args:  [2]Coordinate{{1, 1}, {1, 1}},
			wantB: 1,
			wantL: 1,
		},
		{
			name:  "multiple cells, orientated bottom to up",
			args:  [2]Coordinate{{1, 1}, {3, 3}},
			wantB: 3,
			wantL: 1,
		},
		{
			name:  "multiple cells, orientated up to bottom",
			args:  [2]Coordinate{{3, 3}, {1, 1}},
			wantB: 3,
			wantL: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, b := sortY(tt.args)
			assert.Equal(t, tt.wantL, l)
			assert.Equal(t, tt.wantB, b)
		})
	}
}

func TestSortX(t *testing.T) {
	tests := []struct {
		name  string
		args  [2]Coordinate
		wantB uint
		wantL uint
	}{
		{
			name:  "single cell",
			args:  [2]Coordinate{{1, 1}, {1, 1}},
			wantB: 1,
			wantL: 1,
		},
		{
			name:  "multiple cells, orientated left to right",
			args:  [2]Coordinate{{1, 1}, {3, 3}},
			wantB: 3,
			wantL: 1,
		},
		{
			name:  "multiple cells, orientated right to left",
			args:  [2]Coordinate{{3, 3}, {1, 1}},
			wantB: 3,
			wantL: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, b := sortX(tt.args)
			assert.Equal(t, tt.wantL, l)
			assert.Equal(t, tt.wantB, b)
		})
	}
}
