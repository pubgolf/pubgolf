package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var errUnexpectedFormat = errors.New("unexpected format error")

type malformedRequestError struct {
	OriginalErr error
	Status      errorCode
	Msg         string
}

func (mr *malformedRequestError) Error() string {
	return mr.Msg
}

func parseJSONRequest(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"

		return &malformedRequestError{OriginalErr: errUnexpectedFormat, Status: errorCodeMalformedRequest, Msg: msg}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: "Request body contains badly-formed JSON"}

		case errors.As(err, &unmarshalTypeError):
			return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			msg := "Request body contains unknown field " + strings.TrimPrefix(err.Error(), "json: unknown field ")

			return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: "Request body must not be empty"}

		case err.Error() == "http: request body too large":
			return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: "Request body must not be larger than 1MB"}

		default:
			return fmt.Errorf("unexpected error decoding JSON: %w", err)
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return &malformedRequestError{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: "Request body must only contain a single JSON object"}
	}

	return nil
}

func guardParseJSONRequest(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
	if err == nil {
		return true
	}

	var mr *malformedRequestError
	if errors.As(err, &mr) {
		newErrorResponse(ctx, mr.Status, mr.Msg, fmt.Errorf("parse JSON request body: %w", mr.OriginalErr)).Render(w, r)
	} else {
		newErrorResponse(ctx, errorCodeGenericNonRetryable, "Unknown error parsing request body", fmt.Errorf("parse JSON request body: %w", err)).Render(w, r)
	}

	return false
}
