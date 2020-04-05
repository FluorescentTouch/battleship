package battlefield

import (
	"strings"

	"my/battleship/coordinates"
)

type ship struct {
	c          [2]coordinates.Coordinate
	inner      coordinates.Coordinates
	outer      coordinates.Coordinates
	aliveCells int
	isKnocked  bool
}

func newShip(p1, p2 coordinates.Coordinate) *ship {
	c := [2]coordinates.Coordinate{p1, p2}
	in, out := coordinates.GetInnerOuterCells(c)

	return &ship{
		c:          c,
		inner:      in,
		outer:      out,
		aliveCells: len(in),
	}
}

func makeShipsFromCoords(coords string) ([]*ship, error) {
	if len(coords) == 0 {
		return nil, errorInvalidCoordinate
	}

	s := strings.Split(coords, ",")

	ships := make([]*ship, 0, len(s))

	for _, sc := range s {
		l := strings.Split(sc, " ")
		if len(l) != 2 {
			return nil, errorInvalidCoordinate
		}
		p1, ok := coordinates.ConvertCoordinate(l[0])
		if !ok {
			return nil, errorInvalidCoordinate
		}
		p2, ok := coordinates.ConvertCoordinate(l[1])
		if !ok {
			return nil, errorInvalidCoordinate
		}
		ships = append(ships, newShip(p1, p2))
	}
	return ships, nil
}
