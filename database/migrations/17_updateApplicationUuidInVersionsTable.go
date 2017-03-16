package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateApplicationUUIDInVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "17",
		Up: []string{`
      UPDATE versions v, applications a SET v.application_uuid = a.uuid WHERE v.application_id = a.id;
    `},
		Down: []string{``},
	}
}
