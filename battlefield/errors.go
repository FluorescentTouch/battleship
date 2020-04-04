package battlefield

import (
	"encoding/json"
)

// HTTPError represents json error with http code and error.
type HTTPError struct {
	Err  string `json:"err"`
	Code int    `json:"-"`
}

// Error is implementation of error interface.
func (e HTTPError) Error() string {
	return e.Err
}

// StatusCode is implementation of StatusCoder interface.
func (e HTTPError) StatusCode() int {
	return e.Code
}

// MarshalJSON is implementation of Marshaller interface.
func (e HTTPError) MarshalJSON() ([]byte, error) {
	type t struct {
		Err  string `json:"err,omitempty"`
		Code int    `json:"-"`
	}

	// Casting is needed because json.Marshal(e) creates infinite recursion
	// to the MarshalJSON method of HTTPError.
	return json.Marshal(t(e))
}

var (
	errorInvalidCreateFieldParams = HTTPError{
		Err:  "invalid input params",
		Code: 400,
	}

	errorInvalidFieldSize = HTTPError{
		Err:  "Field size is invalid",
		Code: 400,
	}

	errorFieldAlreadySet = HTTPError{
		Err:  "field is already set",
		Code: 409,
	}
)
