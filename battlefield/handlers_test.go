package battlefield

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewHandlers(t *testing.T) {
	l := logrus.New()
	s := &Service{logger: l}
	e := Endpoints{service: s, logger: l}
	want := Handlers{logger: l, e: e}
	got := NewHandlers(l, e)
	assert.Equal(t, want, got)
}

func TestHandlers_CreateBattleField(t *testing.T) {
	testifyServiceMock := NewTestifyServiceMock(t)

	type args struct {
		method string
		url    string
		body   string
	}

	tests := []struct {
		name       string
		args       args
		setup      func()
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			args: args{
				url:    "/create-matrix",
				method: http.MethodPost,
				body:   fmt.Sprintf(`{"range": %d}`, maxFieldSize-1),
			},
			setup: func() {
				testifyServiceMock.On(
					"createField",
					maxFieldSize-1,
				).Return(nil).Once()
			},
			wantStatus: http.StatusCreated,
			wantBody:   "{}",
		},
		{
			name: "error, invalid request body",
			args: args{
				url:    "/create-matrix",
				method: http.MethodPost,
				body:   "{totally not a valid json]",
			},
			setup:      func() {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"invalid input params"}`,
		},
		{
			name: "error, service errorInvalidFieldSize error",
			args: args{
				url:    "/create-matrix",
				method: http.MethodPost,
				body:   fmt.Sprintf(`{"range": %d}`, maxFieldSize+1),
			},
			setup: func() {
				testifyServiceMock.On(
					"createField",
					maxFieldSize+1,
				).Return(errorInvalidFieldSize).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"field size is invalid"}`,
		},
		{
			name: "error, service errorFieldAlreadySet error",
			args: args{
				url:    "/create-matrix",
				method: http.MethodPost,
				body:   fmt.Sprintf(`{"range": %d}`, maxFieldSize-1),
			},
			setup: func() {
				testifyServiceMock.On(
					"createField",
					maxFieldSize-1,
				).Return(errorFieldAlreadySet).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody:   `{"err":"field is already set"}`,
		},
		{
			name: "error, service general error",
			args: args{
				url:    "/create-matrix",
				method: http.MethodPost,
				body:   fmt.Sprintf(`{"range": %d}`, maxFieldSize+1),
			},
			setup: func() {
				testifyServiceMock.On(
					"createField",
					maxFieldSize+1,
				).Return(errors.New("something went wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "something went wrong",
		},
	}

	logger := logrus.New()
	r := mux.NewRouter()

	endpoints := NewEndpoints(logger, testifyServiceMock)
	handlers := NewHandlers(logger, endpoints)

	r.HandleFunc("/create-matrix", handlers.CreateBattleField)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer testifyServiceMock.AssertExpectations(t)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				tt.args.method,
				tt.args.url,
				strings.NewReader(tt.args.body),
			)
			r.ServeHTTP(res, req)

			assert.Equal(t, res.Code, tt.wantStatus)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(res.Body.String()))
		})
	}
}

func TestHandlers_ClearBattleField(t *testing.T) {
	testifyServiceMock := NewTestifyServiceMock(t)

	type args struct {
		method string
		url    string
	}

	tests := []struct {
		name       string
		args       args
		setup      func()
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			args: args{
				url:    "/clear",
				method: http.MethodPost,
			},
			setup: func() {
				testifyServiceMock.On(
					"clearField",
				).Return(nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   "{}",
		},
		{
			name: "error, service general error",
			args: args{
				url:    "/clear",
				method: http.MethodPost,
			},
			setup: func() {
				testifyServiceMock.On(
					"clearField",
				).Return(errors.New("something went wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "something went wrong",
		},
	}

	logger := logrus.New()
	r := mux.NewRouter()

	endpoints := NewEndpoints(logger, testifyServiceMock)
	handlers := NewHandlers(logger, endpoints)

	r.HandleFunc("/clear", handlers.ClearBattleField)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer testifyServiceMock.AssertExpectations(t)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				tt.args.method,
				tt.args.url,
				nil,
			)
			r.ServeHTTP(res, req)

			assert.Equal(t, res.Code, tt.wantStatus)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(res.Body.String()))
		})
	}
}

func TestHandlers_AddShips(t *testing.T) {
	testifyServiceMock := NewTestifyServiceMock(t)

	type args struct {
		method string
		url    string
		body   string
	}

	tests := []struct {
		name       string
		args       args
		setup      func()
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			args: args{
				url:    "/ship",
				method: http.MethodPost,
				body:   `{"coords": "A1 A1"}`,
			},
			setup: func() {
				testifyServiceMock.On(
					"addShipsByCoordinates",
					"A1 A1",
				).Return(nil).Once()
			},
			wantStatus: http.StatusCreated,
			wantBody:   "{}",
		},
		{
			name: "error, invalid request body",
			args: args{
				url:    "/ship",
				method: http.MethodPost,
				body:   "{totally not a valid json]",
			},
			setup:      func() {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"invalid input params"}`,
		},
		{
			name: "error, service error",
			args: args{
				url:    "/ship",
				method: http.MethodPost,
				body:   `{"coords": "A1 A1"}`,
			},
			setup: func() {
				testifyServiceMock.On(
					"addShipsByCoordinates",
					"A1 A1",
				).Return(errorShipsAlreadyAdded).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"ships are already added"}`,
		},
		{
			name: "error, service general error",
			args: args{
				url:    "/ship",
				method: http.MethodPost,
				body:   `{"coords": "A1 A1"}`,
			},
			setup: func() {
				testifyServiceMock.On(
					"addShipsByCoordinates",
					"A1 A1",
				).Return(errors.New("something went wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "something went wrong",
		},
	}

	logger := logrus.New()
	r := mux.NewRouter()

	endpoints := NewEndpoints(logger, testifyServiceMock)
	handlers := NewHandlers(logger, endpoints)

	r.HandleFunc("/ship", handlers.AddShips)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer testifyServiceMock.AssertExpectations(t)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				tt.args.method,
				tt.args.url,
				strings.NewReader(tt.args.body),
			)
			r.ServeHTTP(res, req)

			assert.Equal(t, res.Code, tt.wantStatus)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(res.Body.String()))
		})
	}
}

func TestHandlers_Shot(t *testing.T) {
	testifyServiceMock := NewTestifyServiceMock(t)

	type args struct {
		method string
		url    string
		body   string
	}

	tests := []struct {
		name       string
		args       args
		setup      func()
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			args: args{
				url:    "/shot",
				method: http.MethodPost,
				body:   `{"coord": "A1"}`,
			},
			setup: func() {
				testifyServiceMock.On(
					"shot",
					"A1",
				).Return(shotResult{
					Destroy: true,
					Knock:   true,
					End:     true,
				}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"destroy":true,"knock":true,"end":true}`,
		},
		{
			name: "error, invalid request body",
			args: args{
				url:    "/shot",
				method: http.MethodPost,
				body:   "{totally not a valid json]",
			},
			setup:      func() {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"invalid input params"}`,
		},
		{
			name: "error, service error",
			args: args{
				url:    "/shot",
				method: http.MethodPost,
				body:   `{"coord": "A1"}`,
			},
			setup: func() {
				testifyServiceMock.On(
					"shot",
					"A1",
				).Return(shotResult{}, errorShipsNotPlaced).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"ships not placed yet"}`,
		},
		{
			name: "error, service general error",
			args: args{
				url:    "/shot",
				method: http.MethodPost,
				body:   `{"coord": "A1"}`,
			},
			setup: func() {
				testifyServiceMock.On(
					"shot",
					"A1",
				).Return(shotResult{}, errors.New("something went wrong")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "something went wrong",
		},
	}

	logger := logrus.New()
	r := mux.NewRouter()

	endpoints := NewEndpoints(logger, testifyServiceMock)
	handlers := NewHandlers(logger, endpoints)

	r.HandleFunc("/shot", handlers.Shot)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer testifyServiceMock.AssertExpectations(t)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				tt.args.method,
				tt.args.url,
				strings.NewReader(tt.args.body),
			)
			r.ServeHTTP(res, req)

			assert.Equal(t, tt.wantStatus, res.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(res.Body.String()))
		})
	}
}
