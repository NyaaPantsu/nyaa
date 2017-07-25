package oauth_client

import (
	"net/http"

	"github.com/pkg/errors"
)

type RichError struct {
	Status int
	error
}

var (
	ErrNotFound = &RichError{
		Status: http.StatusNotFound,
		error:  errors.New("Not found"),
	}
)
