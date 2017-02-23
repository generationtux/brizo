package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func createVersionsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "4",
		Up: []string{`
      CREATE TABLE versions (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        uuid varchar(36) NOT NULL,
        name varchar(255) NOT NULL,
        slug varchar(255) NOT NULL,
        replicas int(11) NOT NULL DEFAULT '0',
        environment_id int(10) unsigned NOT NULL,
        spec json DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY uix_versions_uuid (uuid)
      );
    `},
		Down: []string{`
      DROP TABLE IF EXISTS versions;
    `},
	}
}
