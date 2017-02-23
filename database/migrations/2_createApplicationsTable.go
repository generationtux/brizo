package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func createApplicationsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "2",
		Up: []string{`
      CREATE TABLE applications (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        uuid varchar(36) NOT NULL,
        name varchar(255) NOT NULL,
        slug varchar(255) NOT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY uix_applications_uuid (uuid),
        UNIQUE KEY uix_applications_name (name),
        UNIQUE KEY uix_applications_slug (slug)
      )
    `},
		Down: []string{`
      DROP TABLE IF EXISTS applications;
    `},
	}
}
