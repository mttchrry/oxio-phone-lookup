package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mttchrry/oxio-phone-lookup/pkg/utils/errors"
)

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	fmt.Printf("\nerror occurred in request: %v", err)

	switch {
	case errors.Is(err, errors.ErrInvalidRequest):
		fallthrough
	case errors.Is(err, errors.ErrValidation):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, errors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, errors.ErrUnknown):
		fallthrough
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	errJSON := struct {
		Error string `json:"error"`
	}{
		Error: strings.Split(err.Error(), errors.ErrSeperator)[0], // TODO we may need to strip additional error information
	}

	data, err := json.Marshal(errJSON)
	if err != nil {
		fmt.Printf("failed to serialize error response: %v", err)
		data = []byte(`{"error": "internal server error"}`)
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Printf("failed to write error response: %v", err)
	}
}
