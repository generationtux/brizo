package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addApplicationAndVersionColumnsToEnvironmentsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "12",
		Up: []string{`
      ALTER TABLE environments
        ADD COLUMN application_uuid varchar(36) DEFAULT NULL AFTER application_id,
        ADD COLUMN version_id int(10) unsigned DEFAULT NULL AFTER application_uuid,
        ADD COLUMN version_uuid varchar(36) NOT NULL AFTER version_id
      ;
    `},
		Down: []string{`
      ALTER TABLE environments
        DROP COLUMN application_uuid,
        DROP COLUMN version_id,
		DROP COLUMN version_uuid
      ;
    `},
	}
}
