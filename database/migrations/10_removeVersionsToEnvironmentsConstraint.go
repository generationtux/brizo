package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func removeVersionsToEnvironmentsConstraint() *migrate.Migration {
	return &migrate.Migration{
		Id: "10",
		Up: []string{`
      ALTER TABLE versions
        DROP INDEX uix_versions_name_environment_id
      ;
    `},
		Down: []string{`
      ALTER TABLE versions
        ADD UNIQUE INDEX uix_versions_name_environment_id (name, environment_id)
      ;
    `},
	}
}
