package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addUuidColumnToEnvironmentConfigsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "7",
		Up: []string{`
      ALTER TABLE environment_configs
	ADD COLUMN uuid varchar(36) NOT NULL AFTER id DEFAULT uuid(),
        ADD CONSTRAINT uix_environment_configs_uuid UNIQUE (uuid)
      ;
    `},
		Down: []string{`
      ALTER TABLE environment_configs
        DROP CONSTRAINT uix_environment_configs_uuid,
        DROP COLUMN uuid
      ;
    `},
	}
}
