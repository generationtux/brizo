package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateVersionsIDInEnvironmentsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "15",
		Up: []string{`
      UPDATE environments e, versions v SET e.version_id = v.id WHERE e.application_uuid = v.application_uuid;
    `},
		Down: []string{``},
	}
}
