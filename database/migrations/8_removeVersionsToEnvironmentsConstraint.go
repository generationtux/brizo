package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func removeVersionsToEnvironmentsConstraint() *migrate.Migration {
	return &migrate.Migration{
		Id: "8",
		Up: []string{`
	      ALTER TABLE versions
	        DROP INDEX uix_versions_name_environment_id
	      ;
    `},
		Down: []string{``},
	}
}
