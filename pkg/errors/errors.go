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
		ErrUrlNotValid: {
			http.StatusBadRequest,
			ErrUrlNotValid.Error(),
		},
	}
)

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrUnableToCreateLink = errors.New("unable to create link")

	ErrLinkNotFound = errors.New("link not found")
	ErrUrlNotValid  = errors.New("url is not valid")
)

type ApiError struct {
	err           error
	internalError error
}

func (ae ApiError) Error() string {
	if ae.internalError == nil {
		return fmt.Sprintf("[error]: %s", ae.err)
	}
	return fmt.Sprintf("[error]: %s", ae.internalError)
}

func (ae ApiError) Unwrap() error {
	return ae.err
}

func NewAppError(err, internal error) *ApiError {
	return &ApiError{
		err:           err,
		internalError: internal,
	}
}

func BadRequestError() *ApiError {
	return NewAppError(ErrBadRequest, nil)
}

func InternalError(err error) *ApiError {
	return NewAppError(ErrInternalServer, err)
}

func NotFoundError() *ApiError {
	return NewAppError(ErrLinkNotFound, nil)
}
