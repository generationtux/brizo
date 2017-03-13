package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateUUIDColumnForEnvironmentConfigsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "8",
		Up: []string{`
      UPDATE environment_configs
        SET uuid = uuid()
        WHERE uuid = ''
      ;
    `},
		Down: []string{``},
	}
}
