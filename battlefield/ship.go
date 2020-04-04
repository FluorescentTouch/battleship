package battlefield

import (
	"strings"

	"my/battleship/coordinates"
)

type ship struct {
	c [2]coordinates.Coordinate
}

func newShip(p1, p2 coordinates.Coordinate) *ship {
	return &ship{
		c: [2]coordinates.Coordinate{p1, p2},
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
