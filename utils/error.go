package utils

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// ErrorResponse
type errorResponse struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp"`
}

func (resp *errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	resp.Timestamp = time.Now().Unix()
	return nil
}

func ErrorResponse(message string, code int) render.Renderer {
	return &errorResponse{Message: message, Code: code}
}

// ServiceError
type serviceError struct {
	Code    int
	Message string
}

func (e *serviceError) Error() string {
	return e.Message
}

func ServiceError(message string, code int) error {
	return &serviceError{Message: message, Code: code}
}

func ServiceErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusInternalServerError
	message := err.Error()

	serr, ok := err.(*serviceError)
	if ok {
		code = serr.Code
		message = serr.Message
	}
	render.Render(w, r, ErrorResponse(message, code))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func RecoverIfError(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		err := rec.(error)
		log.Printf("Detected error: %s`n", err.Error())
		ServiceErrorResponse(w, r, err)
	}
}
