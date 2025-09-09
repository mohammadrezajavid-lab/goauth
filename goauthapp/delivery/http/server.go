package http

import (
	"context"
	_ "github.com/mohammadrezajavid-lab/goauth/docs"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/delivery/http/middleware"
	"github.com/mohammadrezajavid-lab/goauth/pkg/httpserver"
	"github.com/mohammadrezajavid-lab/goauth/pkg/token"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	HTTPServer  *httpserver.Server
	Handler     Handler
	RateLimiter *middleware.RateLimiter
	TokenMaker  token.Maker
}

func New(server *httpserver.Server, handler Handler, rateLimiter *middleware.RateLimiter, tokenMaker token.Maker) Server {
	return Server{
		HTTPServer:  server,
		Handler:     handler,
		RateLimiter: rateLimiter,
		TokenMaker:  tokenMaker,
	}
}

func (s Server) Serve() error {
	s.RegisterRoutes()
	if err := s.HTTPServer.Start(); err != nil {
		return err
	}
	return nil
}

func (s Server) RegisterRoutes() {
	router := s.HTTPServer.GetRouter()
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := router.Group("/v1")

	// --- Auth Routes (Public) ---
	authGroup := v1.Group("/auth")
	authGroup.POST("/generateotp", s.Handler.GenerateOTPCode, s.RateLimiter.RateLimitMiddleware)
	authGroup.POST("/verify", s.Handler.VerifyAndLoginOrRegister)

	// --- User Management Routes (Protected by JWT) ---
	authMiddleware := middleware.NewAuthMiddleware(s.TokenMaker)
	usersGroup := v1.Group("/users")
	usersGroup.Use(authMiddleware.RequireAuth)

	usersGroup.GET("/:id", s.Handler.GetUser)
	usersGroup.GET("", s.Handler.ListUsers)
}

func (s Server) Stop(ctx context.Context) error {
	return s.HTTPServer.Stop(ctx)
}
