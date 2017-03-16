package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func updateVersionsUUIDInEnvironmentsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "16",
		Up: []string{`
      UPDATE environments e, versions v SET e.version_uuid = v.uuid WHERE e.version_id = v.id;
    `},
		Down: []string{``},
	}
}
