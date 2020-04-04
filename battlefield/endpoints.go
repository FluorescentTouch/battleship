package battlefield

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type service interface {
	createField(size uint) error
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
