package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addUUIDColumnUniqueConstraintToEnvironmentConfigsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "9",
		Up: []string{`
      ALTER TABLE environment_configs
        ADD CONSTRAINT uix_environment_configs_uuid UNIQUE (uuid)
      ;
    `},
		Down: []string{`
      ALTER TABLE environment_configs
        DROP CONSTRAINT uix_environment_configs_uuid
      ;
    `},
	}
}
