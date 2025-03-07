package responsex

import (
	"net/http"
)

type ApiSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuccessOpt func(*ApiSuccess)

func WithSuccessMessage(msg string) SuccessOpt {
	return func(e *ApiSuccess) {
		e.Message = msg
	}
}

func NewApiSuccess(code int, opt ...SuccessOpt) *ApiSuccess {
	e := &ApiSuccess{
		Code:    code,
		Message: http.StatusText(code),
	}

	for _, o := range opt {
		o(e)
	}
	return e
}
func (e *ApiSuccess) WithMessage(msg string) *ApiSuccess {
	e.Message = msg
	return e
}
func (e *ApiSuccess) WithData(data any) *ApiSuccess {
	e.Data = data
	return e
}
