package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func createEnvironmentsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "3",
		Up: []string{`
      CREATE TABLE environments (
		  id int(10) unsigned NOT NULL AUTO_INCREMENT,
          created_at timestamp NULL DEFAULT NULL,
          updated_at timestamp NULL DEFAULT NULL,
          uuid varchar(36) NOT NULL,
          name varchar(255) NOT NULL,
          slug varchar(255) NOT NULL,
          application_id bigint(20) unsigned DEFAULT NULL,
          PRIMARY KEY (id),
          UNIQUE KEY uix_environments_uuid (uuid)
      );
    `},
		Down: []string{`
      DROP TABLE IF EXISTS environments;
    `},
	}
}
