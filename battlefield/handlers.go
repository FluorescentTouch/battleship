package battlefield

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Handlers collects all handlers that responds to an HTTP request.
type Handlers struct {
	e      Endpoints
	logger *logrus.Logger
}

// StatusCoder describes responses that carry status Code with themselves.
type StatusCoder interface {
	StatusCode() int
}

// NewHandlers creates new Handlers.
func NewHandlers(l *logrus.Logger, e Endpoints) Handlers {
	return Handlers{logger: l, e: e}
}

// CreateBattleField handles request for creating battlefield
// @Title CreateBattleField
// @Tags BattleField
// @Accept json
// @Description create new battlefield with provided size
// @Summary create new battlefield
// @Success 201
// @Failure 400 {object} battlefield.HTTPError
// @Failure 409 {object} battlefield.HTTPError
// @Failure 500 {string} string
// @Router /create-matrix [post]
// @Param model body battlefield.CreateFieldRequest true "createParams"
func (h Handlers) CreateBattleField(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handlers: CreateBattleField started")

	req := CreateFieldRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Errorf("Handlers: CreateBattleField: can't decode request: %v", err)
		handleErrorResponse(w, errorInvalidInputParams)
		return
	}
	resp, err := h.e.createFieldEndpoint(req)
	if err != nil {
		h.logger.Errorf("Handlers: CreateBattleField: can't create Field: %v", err)
		handleErrorResponse(w, err)
		return
	}

	h.logger.Infof("NEW BATTLEFIELD CREATED WITH SIZE %d", req.Size)
	handleOKResponse(w, resp)
}

// ClearBattleField handles request for clearing battlefield
// @Title ClearBattleField
// @Tags BattleField
// @Accept json
// @Description clear the battlefield
// @Summary clear the battlefield
// @Success 200
// @Failure 500 {string} string
// @Router /clear [post]
func (h Handlers) ClearBattleField(w http.ResponseWriter, _ *http.Request) {
	h.logger.Debug("Handlers: ClearBattleField started")

	resp, err := h.e.clearFieldEndpoint()
	if err != nil {
		h.logger.Errorf("Handlers: ClearBattleField: can't clear Field: %v", err)
		handleErrorResponse(w, err)
		return
	}

	h.logger.Info("BATTLEFIELD HAS BEEN CLEARED")
	handleOKResponse(w, resp)
}

// AddShips handles request for adding ships to battlefield
// @Title AddShips
// @Tags Ships
// @Accept json
// @Description add ships to battlefield
// @Description input params should be like this:
// @Description "A1 B2,C4 C6,E7 F8" where first coordinate is one corner of ship, second - other.
// @Description ships can be square or rectangular
// @Description ships can't be placed on top of each other and near each other.
// @Summary add ships to battlefield
// @Success 200
// @Failure 400 {object} battlefield.HTTPError
// @Failure 409 {object} battlefield.HTTPError
// @Failure 500 {string} string
// @Router /ship [post]
// @Param model body battlefield.AddShipsRequest true "coordinates"
func (h Handlers) AddShips(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handlers: AddShips started")

	req := AddShipsRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Errorf("Handlers: AddShips: can't decode request: %v", err)
		handleErrorResponse(w, errorInvalidInputParams)
		return
	}
	resp, err := h.e.addShipsEndpoint(req)
	if err != nil {
		h.logger.Errorf("Handlers: AddShips: can't add ships: %v", err)
		handleErrorResponse(w, err)
		return
	}

	h.logger.Infof("SHIPS ADDED")
	handleOKResponse(w, resp)
}

func handleErrorResponse(w http.ResponseWriter, err error) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			contentType, body = "application/json; charset=utf-8", jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if sc, ok := err.(StatusCoder); ok {
		w.WriteHeader(sc.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = w.Write(body)
}

func handleOKResponse(w http.ResponseWriter, resp interface{}) {
	if sc, ok := resp.(StatusCoder); ok {
		w.WriteHeader(sc.StatusCode())
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)
	return
}
