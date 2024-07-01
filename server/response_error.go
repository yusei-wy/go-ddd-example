package server

import (
	"errors"
	"net/http"

	customerror "go_ddd_example/share/custom_error"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

func NewErrorResponse(statusCode int, message string, details string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

func CustomHTTPErrorHandler(err error, ctx echo.Context) {
	if err == nil {
		return
	}

	ctx.Logger().Error(err)

	if ctx.Response().Committed {
		return
	}

	var handlerError customerror.HandlerError
	if !errors.As(err, &handlerError) {
		return
	}

	if handlerError.Context() == customerror.HandlerErrorContextParseError {
		errorResponse := NewErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), handlerError.Message())
		//nolint: errcheck
		ctx.JSON(http.StatusBadRequest, errorResponse)

		return
	}

	var useCaseError customerror.UseCaseError
	if !errors.As(handlerError.Inner(), &useCaseError) {
		return
	}

	statusCode := errorContextToStatusCode(useCaseError.Context())
	errorResponse := NewErrorResponse(statusCode, http.StatusText(statusCode), useCaseError.Error())

	_ = ctx.JSON(statusCode, errorResponse)
}

func errorContextToStatusCode(ctx customerror.UseCaseErrorContext) int {
	switch ctx {
	case customerror.UseCaseErrorContextParseError:
		return http.StatusBadRequest
	case customerror.UseCaseErrorContextNotFound:
		return http.StatusNotFound
	case customerror.UseCaseErrorContextConflict:
		return http.StatusConflict
	case customerror.UseCaseErrorContextInvalidInput:
		return http.StatusUnprocessableEntity
	case customerror.UseCaseErrorContextUnexpected,
		customerror.UsecaseErrorContextDatabase:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
