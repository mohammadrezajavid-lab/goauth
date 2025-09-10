package goauthapp

import (
	"context"
	"errors"
	"fmt"
	authHTTP "github.com/mohammadrezajavid-lab/goauth/goauthapp/delivery/http"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/delivery/http/middleware"
	userRepository "github.com/mohammadrezajavid-lab/goauth/goauthapp/repository/database"
	otpRepository "github.com/mohammadrezajavid-lab/goauth/goauthapp/repository/memory"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	authService "github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	"github.com/mohammadrezajavid-lab/goauth/pkg/database"
	"github.com/mohammadrezajavid-lab/goauth/pkg/httpserver"
	"github.com/mohammadrezajavid-lab/goauth/pkg/logger"
	"github.com/mohammadrezajavid-lab/goauth/pkg/token"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	HTTPServer  authHTTP.Server
	AuthSvc     authService.Service
	AuthHandler authHTTP.Handler
	Config      Config
}

func Setup(config Config, postgresConn *database.Database) *Application {
	log := logger.L()

	httpServer, hErr := httpserver.New(config.HTTPServer)
	if hErr != nil {
		log.Error("Failed to initialize HTTP server", slog.String("error", hErr.Error()))
		panic(hErr)
	}

	otpRepo := otpRepository.NewGoCacheOTPRepository(config.OTPCache)
	userRepo := userRepository.NewUserRepository(postgresConn.Pool)
	tokenManager, jErr := token.NewJWTMaker(config.JWT)
	if jErr != nil {
		log.Error("Failed initial tokenManager.", slog.String("error", jErr.Error()))
		panic(jErr)
	}
	authSvc := authService.NewService(userRepo, otpRepo, tokenManager)
	authValidator := goauth.NewValidator()
	authHandler := authHTTP.NewHandler(authSvc, authValidator)
	rateLimiter := middleware.NewRateLimiter(config.RateLimiter)
	authHTTPServer := authHTTP.New(httpServer, authHandler, rateLimiter, tokenManager)

	return &Application{
		HTTPServer:  authHTTPServer,
		AuthSvc:     authSvc,
		AuthHandler: authHandler,
		Config:      config,
	}
}

func (app *Application) Start() {
	log := logger.L()
	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.startHTTPServer(&wg)

	log.Info("Auth application started.")

	<-ctx.Done()
	log.Info("Shutdown signal received...")

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), app.Config.TotalShutdownTimeout)
	defer cancel()

	if app.shutdownServers(shutdownTimeoutCtx) {
		log.Info("Servers shutdown gracefully")
	} else {
		log.Warn("Shutdown timed out, exiting application")
		os.Exit(1)
	}

	wg.Wait()
	log.Info("Auth Application stopped")
}

func (app *Application) startHTTPServer(wg *sync.WaitGroup) {
	log := logger.L()
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("HTTP server starting...",
			slog.String("host", app.Config.HTTPServer.Host),
			slog.Int("port", app.Config.HTTPServer.Port))

		if err := app.HTTPServer.Serve(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error(
					"HTTP server failed",
					slog.Int("port", app.Config.HTTPServer.Port),
					slog.Any("error", err),
				)
				panic(err)
			}
		}

		log.Info("HTTP server stopped",
			slog.String("host", app.Config.HTTPServer.Host),
			slog.Int("port", app.Config.HTTPServer.Port))
	}()
}

func (app *Application) shutdownServers(ctx context.Context) bool {
	log := logger.L()
	log.Info("Starting server shutdown process...")

	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup

		shutdownWg.Add(1)
		go app.shutdownHTTPServer(ctx, &shutdownWg)

		shutdownWg.Wait()
		close(shutdownDone)
		log.Info("All servers have been shutdown successfully.")
	}()

	select {
	case <-shutdownDone:
		return true
	case <-ctx.Done():
		return false
	}
}

func (app *Application) shutdownHTTPServer(parentCtx context.Context, wg *sync.WaitGroup) {
	log := logger.L()
	defer wg.Done()
	log.Info("Starting graceful shutdown for HTTP server", "port", app.Config.HTTPServer.Port)

	httpCtx, cancel := context.WithTimeout(parentCtx, app.Config.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := app.HTTPServer.Stop(httpCtx); err != nil {
		log.Error("HTTP server graceful shutdown failed", "error", err)
	} else {
		log.Info("HTTP server shutdown successfully")
	}
}
