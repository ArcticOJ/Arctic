//go:build headless

package main

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/seed"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/ArcticOJ/blizzard/v0/migrations"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/migrate"
	"os"
	"strings"
)

var migrator *migrate.Migrator

var _init = &cobra.Command{
	Use:   "init",
	Short: "create migration tables",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrator.Init(cmd.Context())
	},
}

var _migrate = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := migrator.Lock(cmd.Context()); err != nil {
			return err
		}
		defer migrator.Unlock(cmd.Context())
		group, err := migrator.Migrate(cmd.Context())
		if err != nil {
			return err
		}
		if group.IsZero() {
			logger.Global.Info().Msg("there are no new migrations to run, database is up to date.")
			return nil
		}
		logger.Global.Info().Msgf("migrated to %s", group)
		return nil
	},
}

var rollback = &cobra.Command{
	Use:   "rollback",
	Short: "rollback the last migration group",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := migrator.Lock(cmd.Context()); err != nil {
			return err
		}
		defer migrator.Unlock(cmd.Context())
		group, err := migrator.Rollback(cmd.Context())
		if err != nil {
			return err
		}
		if group.IsZero() {
			logger.Global.Info().Msg("there are no groups to roll back")
			return nil
		}
		logger.Global.Info().Msgf("rolled back %s", group)
		return nil
	},
}

var lock = &cobra.Command{
	Use:   "lock",
	Short: "lock migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrator.Lock(cmd.Context())
	},
}

var unlock = &cobra.Command{
	Use:   "unlock",
	Short: "unlock migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrator.Unlock(cmd.Context())
	},
}

var createGo = &cobra.Command{
	Use:   "create_go",
	Short: "create Go migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, "_")
		mf, err := migrator.CreateGoMigration(cmd.Context(), name)
		if err != nil {
			return err
		}
		logger.Global.Info().Msgf("created migration %s (%s)", mf.Name, mf.Path)
		return nil
	},
}

var createSQL = &cobra.Command{
	Use:   "create_sql",
	Short: "create up and down SQL migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, "_")
		files, err := migrator.CreateSQLMigrations(cmd.Context(), name)
		if err != nil {
			return err
		}
		for _, mf := range files {
			logger.Global.Info().Msgf("created migration %s (%s)", mf.Name, mf.Path)
		}
		return nil
	},
}

var status = &cobra.Command{
	Use:   "status",
	Short: "print migrations status",
	RunE: func(cmd *cobra.Command, args []string) error {
		ms, err := migrator.MigrationsWithStatus(cmd.Context())
		if err != nil {
			return err
		}
		logger.Global.Info().Msgf(
			"migrations: %s\n"+
				"unapplied migrations: %s\n"+
				"last migration group: %s", ms, ms.Unapplied(), ms.LastGroup())
		return nil
	},
}

var markApplied = &cobra.Command{
	Use:   "mark_applied",
	Short: "mark migrations as applied without actually running them",
	RunE: func(cmd *cobra.Command, args []string) error {
		group, err := migrator.Migrate(cmd.Context(), migrate.WithNopMigration())
		if err != nil {
			return err
		}
		if group.IsZero() {
			logger.Global.Info().Msg("there are no new migrations to mark as applied")
			return nil
		}
		logger.Global.Info().Msgf("marked as applied %s", group)
		return nil
	},
}

var reset = &cobra.Command{
	Use:   "reset",
	Short: "recreate all tables and seed with example data",
	RunE: func(cmd *cobra.Command, args []string) error {
		if e := seed.DropAll(db.Database, cmd.Context()); e != nil {
			return e
		}
		if e := seed.CreateAll(db.Database, cmd.Context()); e != nil {
			return e
		}
		if e := migrator.Reset(cmd.Context()); e != nil {
			return e
		}
		if _, e := migrator.Migrate(cmd.Context(), migrate.WithNopMigration()); e != nil {
			return e
		}
		fixture := dbfixture.New(migrator.DB())
		return fixture.Load(cmd.Context(), os.DirFS("."), "fixture.yml")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "database migration helper",
}

func init() {
	migrator = migrate.NewMigrator(
		db.Database,
		migrations.Migrations,
		migrate.WithTableName("arctic_migrations"),
		migrate.WithLocksTableName("arctic_migration_locks"))
	migrateCmd.AddCommand(
		_init,
		_migrate,
		rollback,
		lock,
		unlock,
		createGo,
		createSQL,
		status,
		markApplied,
		reset,
	)
}
