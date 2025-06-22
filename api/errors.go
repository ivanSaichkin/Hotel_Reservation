package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return ctx.Status(apiError.Code).JSON(apiError)
	}

	apiError := NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func (e Error) Error() string {
	return e.Err
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid ID given",
	}
}

func ErrUnAthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "Unauthorized request",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}
}

func ErrNotfound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + "resorce not found",
	}
}
