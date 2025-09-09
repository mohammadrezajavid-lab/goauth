package command

import (
	"github.com/mohammadrezajavid-lab/goauth/pkg/migrator"
	"github.com/spf13/cobra"
	"log"
)

var up bool
var down bool

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `This command runs the database migrations for the auth service.`,
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func init() {
	migrateCmd.Flags().BoolVar(&up, "up", false, "Run migrations up")
	migrateCmd.Flags().BoolVar(&down, "down", false, "Run migrations down")
	migrateCmd.MarkFlagsMutuallyExclusive("up", "down")
	RootCmd.AddCommand(migrateCmd)
}

func migrate() {
	cfg := loadAppConfig()

	mgr := migrator.New(cfg.PostgresDB, cfg.PathOfMigration)

	if up && down {
		log.Fatalf("Flags --up and --down are mutually exclusive")
	}

	if up {
		log.Println("Running migrations up...")
		mgr.Up()
		log.Println("Migrations up completed.")
	} else if down {
		log.Println("Running migrations down...")
		mgr.Down()
		log.Println("Migrations down completed.")
	} else {
		log.Println("Please specify a migration direction with --up or --down")
	}
}
