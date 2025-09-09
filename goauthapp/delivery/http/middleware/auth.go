package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadrezajavid-lab/goauth/pkg/token"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type AuthMiddleware struct {
	tokenMaker token.Maker
}

func NewAuthMiddleware(tokenMaker token.Maker) *AuthMiddleware {
	return &AuthMiddleware{tokenMaker: tokenMaker}
}

// RequireAuth is an Echo middleware function for requiring JWT authentication.
func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authorizationHeaderKey)
		if len(authHeader) == 0 {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authorization header is not provided"})
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization header format"})
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unsupported authorization type"})
		}

		accessToken := fields[1]
		claims, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
		}

		// Store claims in context
		c.Set(authorizationPayloadKey, claims)
		return next(c)
	}
}
