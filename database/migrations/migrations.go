package migrations

import (
	"github.com/generationtux/brizo/database"
	"github.com/rubenv/sql-migrate"
)

// All gets all migrations to run
func All() []*migrate.Migration {
	return []*migrate.Migration{
		createUsersTable(),
		createApplicationsTable(),
		createEnvironmentsTable(),
		createVersionsTable(),
		createEnvironmentConfigsTable(),
		createAccessTokensTable(),
	}
}

// Run execute migrations forward
func Run() error {
	conn, err := database.Connect()
	if err != nil {
		return err
	}

	m := &migrate.MemoryMigrationSource{Migrations: All()}
	_, err = migrate.Exec(conn.DB(), "mysql", m, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
