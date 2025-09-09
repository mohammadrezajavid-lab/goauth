package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	nethttp "net/http"
)

// handleServiceError maps service layer errors to appropriate HTTP responses.
func handleServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, goauth.ErrOTPNotFound), errors.Is(err, goauth.ErrInvalidOTPCode):
		return c.JSON(nethttp.StatusUnauthorized, echo.Map{
			"error": "Invalid phone number or OTP code.",
		})

	case errors.Is(err, goauth.ErrUserNotFound):
		return c.JSON(nethttp.StatusNotFound, echo.Map{
			"error": err.Error(),
		})

	case errors.Is(err, goauth.ErrUserAlreadyExists):
		return c.JSON(nethttp.StatusConflict, echo.Map{
			"error": err.Error(),
		})
	}

	c.Logger().Error("An internal server error occurred: ", err)
	return c.JSON(nethttp.StatusInternalServerError, echo.Map{
		"error": "An unexpected error occurred on the server.",
	})
}
