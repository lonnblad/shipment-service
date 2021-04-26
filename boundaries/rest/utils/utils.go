package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	RegexpUUID                     = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}"
	ErrResponseInternalServerError = `{"error": {"message": "internal server error"}}`
	ErrResponseNotFound            = `{"error": {"message": "not found"}}`
	ErrResponseMethodNotAllowed    = `{"error": {"message": "method not allowed"}}`
)

// UnmarshalRequest will take an io.ReadCloser and an interface{}.
//
// It will try to decode the data read from the ReadCloser to the
// provided interface{} and then close ReadCloser.
func UnmarshalRequest(body io.ReadCloser, v interface{}) error {
	defer body.Close()

	if err := json.NewDecoder(body).Decode(v); err != nil {
		return fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	return nil
}

// MarshalAndWriteJSONResponse will take an http.ResponseWriter, a
// statusCode and an interface{}.
//
// It will attempt to marshal the interface{} and then call
// WriteJSONResponse with the serialized version.
func MarshalAndWriteJSONResponse(w http.ResponseWriter, statusCode int, v interface{}) []byte {
	response, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to marshal response body of type: %T, error: %s", v, err.Error())

		statusCode = http.StatusInternalServerError
		response = []byte(ErrResponseInternalServerError)
	}

	WriteJSONResponse(w, statusCode, response)

	return response
}

// ErrorResponse is a common structure for communicating error messages.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// WrapErrorAndWriteJSONResponse will take an http.ResponseWriter, a
// statusCode and an error.
//
// It will wrap the error and then call MarshalAndWriteJSONResponse
// with the wrapped error.
func WrapErrorAndWriteJSONResponse(w http.ResponseWriter, statusCode int, err error) {
	resp := ErrorResponse{}
	resp.Error.Message = err.Error()

	MarshalAndWriteJSONResponse(w, statusCode, resp)
}

// WriteJSONResponse will take an http.ResponseWriter, a
// statusCode and a body as a byte slice.
//
// It will set headers and status code and finally write the
// body to the ResponseWriter.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(body); err != nil {
		log.Printf("Failed to write response: %s", err.Error())
	}
}

// HandlerNotFound is an http.Handler which will return a 404,
// with a predefined message.
func HandlerNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusNotFound
		resp := []byte(ErrResponseNotFound)
		WriteJSONResponse(w, code, resp)
	})
}

// HandlerMethodNotAllowed is an http.Handler which will return
// a 405, with a predefined message.
func HandlerMethodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusMethodNotAllowed
		resp := []byte(ErrResponseMethodNotAllowed)
		WriteJSONResponse(w, code, resp)
	})
}

// MiddlewareRecovery will capture any panics and write an
// Internal Server Error response on the http.ResponseWriter.
//
// This Middleware is good to use to not make sure the server
// doesn't crash on panics, like nil pointer or index out of bounds.
func MiddlewareRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if panic := recover(); panic != nil {
				err, ok := panic.(error)
				if !ok {
					err = fmt.Errorf("panic: %v", panic)
				}

				WrapErrorAndWriteJSONResponse(w, http.StatusInternalServerError, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
