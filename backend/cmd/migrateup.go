package cmd

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required to support file protocol
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// migrateupCmd represents the migrateup command
var migrateupCmd = &cobra.Command{
	Use:          "migrateup",
	Short:        "Performs database `up` migration",
	Long:         "Connects to database using the configured connection properties & performs `up` migration",
	SilenceUsage: true, // do not print usage info in case of error
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Info("starting database migration")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			conf.Database.Host, conf.Database.Username, conf.Database.Password, conf.Database.DBName, conf.Database.Port)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}

		driver, err := postgres.WithInstance(db, &postgres.Config{
			MultiStatementEnabled: true,
		})
		if err != nil {
			return err
		}
		m, err := migrate.NewWithDatabaseInstance("file://migrations", conf.Database.DBName, driver)
		if err != nil {
			return err
		}
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				logrus.Info("database migration skipped. no changes detected")
				return nil
			}
			return err
		}
		logrus.Info("database migration completed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
