package apierror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	Errors = map[error]struct {
		Code    int
		Message string
	}{
		ErrInternalServer: {
			http.StatusInternalServerError,
			ErrInternalServer.Error(),
		},
		ErrBadRequest: {
			http.StatusBadRequest,
			ErrBadRequest.Error(),
		},
		ErrUnableToCreateLink: {
			http.StatusConflict,
			ErrUnableToCreateLink.Error(),
		},
		ErrLinkNotFound: {
			http.StatusNotFound,
			ErrLinkNotFound.Error(),
		},
		ErrURLNotValid: {
			http.StatusBadRequest,
			ErrURLNotValid.Error(),
		},
	}
)

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrUnableToCreateLink = errors.New("unable to create link")

	ErrLinkNotFound = errors.New("link not found")
	ErrURLNotValid  = errors.New("url is not valid")
)

type APIError struct {
	err           error
	internalError error
}

func (ae APIError) Error() string {
	if ae.internalError == nil {
		return fmt.Sprintf("[error]: %s", ae.err)
	}

	return fmt.Sprintf("[error]: %s", ae.internalError)
}

func (ae APIError) Unwrap() error {
	return ae.err
}

func NewAPIError(err, internal error) *APIError {
	return &APIError{
		err:           err,
		internalError: internal,
	}
}

func BadRequestError() *APIError {
	return NewAPIError(ErrBadRequest, nil)
}

func InternalError(err error) *APIError {
	return NewAPIError(ErrInternalServer, err)
}

func NotFoundError() *APIError {
	return NewAPIError(ErrLinkNotFound, nil)
}
