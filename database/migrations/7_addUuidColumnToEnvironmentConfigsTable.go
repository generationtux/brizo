package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addUUIDColumnToEnvironmentConfigsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "7",
		Up: []string{`
      ALTER TABLE environment_configs
        ADD COLUMN uuid varchar(36) NOT NULL AFTER id
      ;
    `},
		Down: []string{`
      ALTER TABLE environment_configs
        DROP COLUMN uuid
      ;
    `},
	}
}
