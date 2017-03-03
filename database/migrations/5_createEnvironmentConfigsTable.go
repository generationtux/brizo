package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func createEnvironmentConfigsTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "5",
		Up: []string{`
      CREATE TABLE environment_configs (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        name varchar(255) NOT NULL,
        value varchar(255) NOT NULL,
        environment_uuid varchar(255) NOT NULL,
        PRIMARY KEY (id)
      )
    `},
		Down: []string{`
      DROP TABLE IF EXISTS environment_configs;
    `},
	}
}
