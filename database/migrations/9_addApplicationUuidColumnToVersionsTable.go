package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addApplicationUuidColumnToVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "9",
		Up: []string{`
      ALTER TABLE versions
        ADD COLUMN application_uuid varchar(36) NOT NULL AFTER replicas,
		ADD CONSTRAINT uix_application_uuid UNIQUE (uuid),
		ADD UNIQUE INDEX uix_versions_name_application_uuid (name, application_uuid)
      ;
    `},
		Down: []string{`
	  ALTER TABLE versions
	    DROP CONSTRAINT uix_application_uuid,
		DROP INDEX uix_versions_name_application_uuid,
	    DROP COLUMN application_uuid
	  ;
	`},
	}
}
