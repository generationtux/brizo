package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateApplicationIDInVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "16",
		Up: []string{`
      UPDATE versions v, environments e SET v.application_id = e.application_id WHERE v.environment_id = e.id;
    `},
		Down: []string{``},
	}
}
