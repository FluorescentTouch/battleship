package battlefield

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type service interface {
	createField(size uint) error
	clearField() error
	addShipsByCoordinates(coords string) error
	shot(coordinate string) (shotResult, error)
	state() state
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
	Coords string `json:"Coordinates"`
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

// ShotRequest collect params for shot request.
type ShotRequest struct {
	Coord string `json:"coord"`
}

// ShotResponse created contains params for shot response.
type ShotResponse struct {
	Destroy bool `json:"destroy"`
	Knock   bool `json:"knock"`
	End     bool `json:"end"`
}

// StatusCode implements StatusCoder.
func (r ShotResponse) StatusCode() int {
	return http.StatusOK
}

func (e Endpoints) shotEndpoint(req ShotRequest) (ShotResponse, error) {
	e.logger.Debug("Endpoints: shotEndpoint started")

	res, err := e.service.shot(req.Coord)
	if err != nil {
		return ShotResponse{}, err
	}
	return ShotResponse{
		Destroy: res.Destroy,
		Knock:   res.Knock,
		End:     res.End,
	}, nil
}

// StateResponse defines state response.
type StateResponse struct {
	ShipCount int `json:"ship_count"`
	Destroyed int `json:"destroyed"`
	Knocked   int `json:"knocked"`
	ShotCount int `json:"shot_count"`
}

// StatusCode implements StatusCoder.
func (r StateResponse) StatusCode() int {
	return http.StatusOK
}

func (e Endpoints) stateEndpoint() StateResponse {
	e.logger.Debug("Endpoints: stateEndpoint started")

	state := e.service.state()
	return StateResponse{
		ShipCount: state.shipCount,
		Destroyed: state.destroyed,
		Knocked:   state.knocked,
		ShotCount: state.shotCount,
	}
}
