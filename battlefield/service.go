package battlefield

import (
	"sync"

	"my/battleship/coordinates"

	"github.com/sirupsen/logrus"
)

// Service describes the battlefield service operations.
type Service struct {
	f Field

	logger *logrus.Logger
	sync.RWMutex
}

// NewService creates new Service.
func NewService(l *logrus.Logger) *Service {
	l.Infof("MAXIMUM FIELD SIZE POSSIBLE IS %d", maxFieldSize)
	return &Service{logger: l}
}

func (s *Service) createField(size uint) error {
	s.Lock()
	defer s.Unlock()

	if s.f.isSet {
		return errorFieldAlreadySet
	}

	s.logger.WithField("size", size).Debug("Service: createField started")

	if size < 1 || size > maxFieldSize {
		s.logger.WithField("size", size).
			Error("Field size provided is invalid")
		return errorInvalidFieldSize
	}
	s.f = NewField(size)
	return nil
}

func (s *Service) clearField() error {
	s.Lock()
	defer s.Unlock()

	s.logger.Debug("Service: clearField started")

	s.f = Field{}
	return nil
}

func (s *Service) addShipsByCoordinates(coords string) error {
	s.Lock()
	defer s.Unlock()

	if s.f.shipsAdded {
		return errorShipsAlreadyAdded
	}

	s.logger.WithField("coords", coords).
		Debug("Service: addShipsByCoordinates started")

	ships, err := makeShipsFromCoords(coords)
	if err != nil {
		s.logger.WithField("coords", coords).
			Error("addShipsByCoordinates: invalid coordinates provided")
		return err
	}
	err = s.addShips(ships)
	if err != nil {
		s.logger.WithField("coords", coords).
			Error("addShipsByCoordinates: can't add ships")
		return err
	}
	s.f.shipsAdded = true
	return nil
}

func (s *Service) addShips(ships []*ship) error {
	for _, ship := range ships {
		err := s.placeShip(ship)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) placeShip(sh *ship) error {
	inner, outer := coordinates.GetInnerOuterCells(sh.c)
	// occupy ship cells
	for _, c := range inner {
		if c.X >= s.f.size || c.Y >= s.f.size {
			return errorOutOfBonds
		}
		cell := s.f.field[c.X][c.Y]
		if cell.occupied {
			if cell.ship != nil {
				return errorCellIsOccupiedByShip
			}
			return errorCellIsOccupiedNearby
		}
		cell.occupied = true
		cell.ship = sh
		s.f.field[c.X][c.Y] = cell
	}
	// occupy nearby cells, skip if out of bonds
	for _, c := range outer {
		if c.X >= s.f.size || c.Y >= s.f.size {
			continue
		}
		cell := s.f.field[c.X][c.Y]
		if cell.occupied {
			continue
		}
		cell.occupied = true
		s.f.field[c.X][c.Y] = cell
	}
	return nil
}
