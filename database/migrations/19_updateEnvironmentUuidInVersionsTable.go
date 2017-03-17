package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateEnvironmentUUIDInVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "19",
		Up: []string{`
      UPDATE versions v, environments e SET v.environment_uuid = e.uuid WHERE v.environment_id = e.id;
    `},
		Down: []string{``},
	}
}
