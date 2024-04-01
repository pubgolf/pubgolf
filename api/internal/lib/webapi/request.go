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

type malformedRequest struct {
	OriginalErr error
	Status      errorCode
	Msg         string
}

func (mr *malformedRequest) Error() string {
	return mr.Msg
}

func parseJSONRequest(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"
		return &malformedRequest{OriginalErr: errors.New("unexpected format error"), Status: errorCodeMalformedRequest, Msg: msg}
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
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{OriginalErr: err, Status: errorCodeMalformedRequest, Msg: msg}
	}

	return nil
}

func guardParseJSONRequest(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) (ok bool) {
	if err == nil {
		return true
	}

	var mr *malformedRequest
	if errors.As(err, &mr) {
		newErrorResponse(ctx, mr.Status, mr.Msg, fmt.Errorf("parse JSON request body: %w", mr.OriginalErr)).Render(w, r)
	} else {
		newErrorResponse(ctx, errorCodeGenericNonRetryable, "Unknown error parsing request body", fmt.Errorf("parse JSON request body: %w", err)).Render(w, r)
	}
	return false
}
