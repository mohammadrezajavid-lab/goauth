package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	Limit  int           `koanf:"limit"`
	Window time.Duration `koanf:"window"`
}

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	config   Config
}

func NewRateLimiter(config Config) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
	}
}

// RateLimitMiddleware is the Echo middleware function.
func (rl *RateLimiter) RateLimitMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req goauth.GenerateOTPRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
		}

		// Store the bound request in the context so the handler doesn't need to bind it again.
		c.Set("request", &req)

		phoneNumber := req.PhoneNumber

		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Clean up old requests that are outside the time window.
		now := time.Now()
		var recentRequests []time.Time
		for _, t := range rl.requests[phoneNumber] {
			if now.Sub(t) < rl.config.Window {
				recentRequests = append(recentRequests, t)
			}
		}
		rl.requests[phoneNumber] = recentRequests

		// Check if the limit is exceeded.
		if len(rl.requests[phoneNumber]) >= rl.config.Limit {
			return c.JSON(http.StatusTooManyRequests, echo.Map{"error": "you have made too many requests, please try again later"})
		}

		// Add the current request timestamp.
		rl.requests[phoneNumber] = append(rl.requests[phoneNumber], now)

		return next(c)
	}
}
