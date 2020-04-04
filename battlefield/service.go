package battlefield

import (
	"sync"

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

	s.logger.WithField("size", size).Debug("Service: createField started")

	if size < 1 || size > maxFieldSize {
		s.logger.WithField("size", size).
			Error("Field size provided is invalid")
		return errorInvalidFieldSize
	}
	s.f = NewField(size)
	return nil
}
