package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addApplicationAndEnvironmentColumnsToVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "13",
		Up: []string{`
      ALTER TABLE versions
	  	  ADD COLUMN application_id int(10) unsigned NOT NULL AFTER replicas,
		    ADD COLUMN environment_uuid varchar(36) NOT NULL AFTER environment_id
      ;
    `},
		Down: []string{`
      ALTER TABLE versions
		    DROP COLUMN application_id,
		    DROP COLUMN environment_uuid
      ;
    `},
	}
}
