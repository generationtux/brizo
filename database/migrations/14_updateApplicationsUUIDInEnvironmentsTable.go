package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateApplicationsUUIDInEnvironmentsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "14",
		Up: []string{`
      UPDATE environments e, applications a SET e.application_uuid = a.uuid WHERE e.application_id = a.id;
    `},
		Down: []string{``},
	}
}
