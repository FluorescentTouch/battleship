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

	if s.f.isSet && !s.f.gameIsOver {
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

	s.logger.WithField("coords", coords).
		Debug("Service: addShipsByCoordinates started")

	if s.f.shipsAdded {
		return errorShipsAlreadyAdded
	}

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
	s.f.shipsAlive = len(ships)
	s.f.state.shipCount = len(ships)
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
	// occupy ship cells
	for c := range sh.inner {
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
	for c := range sh.outer {
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

func (s *Service) shot(coordinate string) (shotResult, error) {
	s.Lock()
	defer s.Unlock()

	s.logger.WithField("coordinate", coordinate).
		Debug("Service: shot started")

	if !s.f.shipsAdded {
		return shotResult{}, errorShipsNotPlaced
	}

	c, ok := coordinates.ConvertCoordinate(coordinate)
	if !ok {
		s.logger.WithField("coordinate", coordinate).
			Error("shot: invalid coordinate provided")
		return shotResult{}, errorInvalidCoordinate
	}

	if c.X >= s.f.size || c.Y >= s.f.size {
		return shotResult{}, errorOutOfBonds
	}

	cell := s.f.field[c.X][c.Y]
	if cell.shot {
		return shotResult{}, errorCellAlreadyShot
	}
	cell.shot = true

	res := shotResult{}

	if cell.ship != nil {
		res.Knock = true
		cell.ship.aliveCells--

		if !cell.ship.isKnocked {
			cell.ship.isKnocked = true
			// update global state
			s.f.state.knocked++
		}

		if cell.ship.aliveCells == 0 {
			res.Destroy = true

			// update global state
			s.f.shipsAlive--
			s.f.state.knocked--
			s.f.state.destroyed++
		}
	}

	if s.f.shipsAlive == 0 {
		s.f.gameIsOver = true
		res.End = true
	}
	s.f.field[c.X][c.Y] = cell

	// update global state
	s.f.state.shotCount++
	return res, nil
}

func (s *Service) state() state {
	s.RLock()
	defer s.RUnlock()

	s.logger.Debug("Service: state started")

	return s.f.state
}
