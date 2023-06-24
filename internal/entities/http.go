package entities

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

var (
	ErrDecode = NewHttpErr(http.StatusBadRequest, "problem to decode your input", nil)
)

type HttpErr struct {
	HTTPCode int    `mapstructure:"code" json:"http_code,omitempty"`
	Message  string `mapstructure:"message" json:"message,omitempty"`
	Detail   any    `mapstructure:"detail,omitempty" swaggerignore:"true" json:"detail,omitempty"`
}

func (e *HttpErr) Error() string {
	return fmt.Sprintf("code: %v - message: %v - detail: %v", e.HTTPCode, e.Message, string(toJSON(e.Detail)))
}

func NewHttpErr(httpCode int, message string, detail any) error {
	return &HttpErr{
		HTTPCode: httpCode,
		Message:  message,
		Detail:   detail,
	}
}

func toJSON(v any) []byte {
	bt, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return bt
}
