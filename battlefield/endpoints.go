package battlefield

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type service interface {
	createField(size uint) error
	clearField() error
	addShipsByCoordinates(coords string) error
}

// NewEndpoints creates new Endpoints.
func NewEndpoints(l *logrus.Logger, s service) Endpoints {
	return Endpoints{logger: l, service: s}
}

// Endpoints collects all of the endpoints.
type Endpoints struct {
	service service
	logger  *logrus.Logger
}

// CreateFieldRequest collect params for createField request.
type CreateFieldRequest struct {
	Size uint `json:"range"`
}

// CreateFieldResponse created for swagger docs.
type CreateFieldResponse struct{}

// StatusCode implements StatusCoder.
func (r CreateFieldResponse) StatusCode() int {
	return http.StatusCreated
}

func (e Endpoints) createFieldEndpoint(r CreateFieldRequest) (CreateFieldResponse, error) {
	e.logger.WithField("CreateFieldRequest", r).Debug("Endpoints: createFieldEndpoint started")

	err := e.service.createField(r.Size)
	return CreateFieldResponse{}, err
}

// ClearFieldResponse defined for swagger docs.
type ClearFieldResponse struct{}

// StatusCode implements StatusCoder.
func (r ClearFieldResponse) StatusCode() int {
	return http.StatusOK
}

func (e Endpoints) clearFieldEndpoint() (ClearFieldResponse, error) {
	e.logger.Debug("Endpoints: clearFieldEndpoint started")

	err := e.service.clearField()
	return ClearFieldResponse{}, err
}

// AddShipsRequest collect params for addShips request.
type AddShipsRequest struct {
	Coords string `json:"coords"`
}

// AddShipsResponse created for swagger docs.
type AddShipsResponse struct{}

// StatusCode implements StatusCoder.
func (r AddShipsResponse) StatusCode() int {
	return http.StatusCreated
}

func (e Endpoints) addShipsEndpoint(req AddShipsRequest) (AddShipsResponse, error) {
	e.logger.Debug("Endpoints: addShipsEndpoint started")

	err := e.service.addShipsByCoordinates(req.Coords)
	return AddShipsResponse{}, err
}
