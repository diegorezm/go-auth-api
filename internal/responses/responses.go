package responses

import (
	"net/http"

	"github.com/diegorezm/ticketing/internal/json"
)

func Ok[T any](w http.ResponseWriter, status int, data T) {
	json.Write(w, status, ok(status, data))
}

func OkMessage(w http.ResponseWriter, status int, message string) {
	json.Write(w, status, okMessage(status, message))
}

func Fail(w http.ResponseWriter, status int, message string) {
	json.Write(w, status, failResponse(status, message))
}

func ok[T any](code int, data T) Success[T] {
	return Success[T]{
		Code: code,
		Data: &data,
	}
}

func okMessage(code int, message string) Success[any] {
	return Success[any]{
		Code:    code,
		Message: &message,
	}
}

func failResponse(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
