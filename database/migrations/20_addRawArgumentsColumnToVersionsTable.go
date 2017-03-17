package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func addRawArgumentsColumnToVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "20",
		Up: []string{`
      ALTER TABLE versions
	  	  ADD COLUMN raw_arguments json
      ;
    `},
		Down: []string{`
      ALTER TABLE versions
		    DROP COLUMN raw_arguments
      ;
    `},
	}
}
