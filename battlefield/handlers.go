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
// @Failure 400 {string} battlefield.HTTPError
// @Failure 500 {string} string
// @Router /create-matrix [post]
// @Param model body battlefield.CreateFieldRequest true "size"
func (h Handlers) CreateBattleField(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handlers: CreateBattleField started")

	req := CreateFieldRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Errorf("Handlers: CreateBattleField: can't decode request: %v", err)
		handleErrorResponse(w, errorInvalidCreateFieldParams)
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
