package command

import (
	"github.com/mohammadrezajavid-lab/goauth/goauthapp"
	"github.com/mohammadrezajavid-lab/goauth/pkg/database"
	"github.com/mohammadrezajavid-lab/goauth/pkg/logger"
	"github.com/mohammadrezajavid-lab/goauth/pkg/migrator"
	"github.com/spf13/cobra"
	"log/slog"
)

var migrateUp bool
var migrateDown bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the auth service",
	Long:  `This command starts the main auth service.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	serveCmd.Flags().BoolVar(&migrateUp, "migrate-up", false, "Run migrations up before starting the server")
	serveCmd.Flags().BoolVar(&migrateDown, "migrate-down", false, "Run migrations down before starting the server")
	serveCmd.MarkFlagsMutuallyExclusive("migrate-up", "migrate-down")
	RootCmd.AddCommand(serveCmd)
}

func serve() {
	cfg := loadAppConfig()

	// Initialize logger
	if err := logger.Init(cfg.Logger); err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		return
	}
	defer func() {
		if err := logger.Close(); err != nil {
			slog.Warn("logger close error", "error", err)
		}
	}()
	slog.Info("Logger initialized successfully.")

	// Run migrations if flags are set
	if migrateUp || migrateDown {
		mgr := migrator.New(cfg.PostgresDB, cfg.PathOfMigration)
		if migrateUp {
			slog.Info("Running migrations up...")
			mgr.Up()
			slog.Info("Migrations up completed.")
		}
		if migrateDown {
			slog.Info("Running migrations down...")
			mgr.Down()
			slog.Info("Migrations down completed.")
		}
	}

	slog.Info("Starting Auth Service...")

	// Connect to the database
	databaseConn, cnErr := database.Connect(cfg.PostgresDB)
	if cnErr != nil {
		slog.Error("fatal error occurred", "reason", "failed to connect to database", "error", cnErr)
		return
	}
	defer databaseConn.Close()

	// Setup and start the application
	app := goauthapp.Setup(cfg, databaseConn)
	app.Start()

	slog.Info("Auth service stopped.")
}
