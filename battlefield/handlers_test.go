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
				body:   fmt.Sprintf(`{"size": %d}`, maxFieldSize-1),
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
				body:   fmt.Sprintf(`{"size": %d}`, maxFieldSize+1),
			},
			setup: func() {
				testifyServiceMock.On(
					"createField",
					maxFieldSize+1,
				).Return(errorInvalidFieldSize).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"err":"Field size is invalid"}`,
		},
		{
			name: "error, service general error",
			args: args{
				url:    "/create-matrix",
				method: http.MethodPost,
				body:   fmt.Sprintf(`{"size": %d}`, maxFieldSize+1),
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