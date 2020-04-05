package coordinates

import (
	"strconv"
	"strings"
)

// Coordinate represents coordinate on battlefield.
// Both X and Y starts with zero.
type Coordinate struct {
	X, Y uint
}

// ConvertCoordinate converts string representation of ship's coordinate
// into internal coordinate value.
func ConvertCoordinate(s string) (c Coordinate, ok bool) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if len(s) < 2 {
		return Coordinate{}, false
	}
	x := s[0]
	if x < 'A' || x > 'Z' {
		return Coordinate{}, false
	}
	x = x - 'A'
	// other symbols are number
	ys := s[1:]
	y, err := strconv.Atoi(ys)
	if err != nil {
		return Coordinate{}, false
	}
	return Coordinate{
		X: uint(x),
		Y: uint(y - 1),
	}, true
}

// GetInnerOuterCells calculates and returns ship cells and
// cells that in close vicinity f ship cells.
func GetInnerOuterCells(c [2]Coordinate) (inner, outer []Coordinate) {
	lx, bx := sortX(c)
	ly, by := sortY(c)

	inner = innerCells(lx, bx, ly, by)
	outer = outerCells(lx, bx, ly, by)
	return inner, outer
}

func outerCells(lx, bx, ly, by uint) []Coordinate {
	inner := innerCells(lx, bx, ly, by)

	if lx != 0 {
		lx--
	}
	if ly != 0 {
		ly--
	}
	bx++
	by++

	distX := 1 + (bx - lx)
	distY := 1 + (by - ly)

	totalCoordinates := make([]Coordinate, 0, distX*distY)

	for x := lx; x <= bx; x++ {
		for y := ly; y <= by; y++ {
			totalCoordinates = append(totalCoordinates, Coordinate{X: x, Y: y})
		}
	}

	// exclude inner from totalOuter
	m := make(map[Coordinate]struct{}, len(totalCoordinates))
	for _, c := range totalCoordinates {
		m[c] = struct{}{}
	}
	for _, c := range inner {
		delete(m, c)
	}
	out := make([]Coordinate, 0, len(m))
	for c := range m {
		out = append(out, c)
	}

	return out
}

func innerCells(lx, bx, ly, by uint) []Coordinate {
	distX := 1 + (bx - lx)
	distY := 1 + (by - ly)

	totalCoordinates := make([]Coordinate, 0, distX*distY)

	for x := lx; x <= bx; x++ {
		for y := ly; y <= by; y++ {
			totalCoordinates = append(totalCoordinates, Coordinate{X: x, Y: y})
		}
	}
	return totalCoordinates
}

func sortX(c [2]Coordinate) (lesser, bigger uint) {
	if c[0].X < c[1].X {
		return c[0].X, c[1].X
	}
	return c[1].X, c[0].X
}

func sortY(c [2]Coordinate) (lesser, bigger uint) {
	if c[0].Y < c[1].Y {
		return c[0].Y, c[1].Y
	}
	return c[1].Y, c[0].Y
}
