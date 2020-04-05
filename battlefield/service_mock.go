package battlefield

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// TestifyServiceMock is a mock implementation of Service interface.
type TestifyServiceMock struct {
	mock.Mock
}

// NewTestifyServiceMock creates a new instance of ServiceMock
// and set output on the testing logger.
func NewTestifyServiceMock(t *testing.T) *TestifyServiceMock {
	m := &TestifyServiceMock{}
	m.Test(t)
	return m
}

// createField is mock implementation.
func (r *TestifyServiceMock) createField(size uint) error {
	results := r.Called(size)
	return results.Error(0)
}

// clearField is mock implementation.
func (r *TestifyServiceMock) clearField() error {
	results := r.Called()
	return results.Error(0)
}

// addShipsByCoordinates is mock implementation.
func (r *TestifyServiceMock) addShipsByCoordinates(coords string) error {
	results := r.Called(coords)
	return results.Error(0)
}

// addShipsByCoordinates is mock implementation.
func (r *TestifyServiceMock) shot(coords string) (shotResult, error) {
	results := r.Called(coords)
	return results.Get(0).(shotResult), results.Error(1)
}
