package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func createUsersTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "1",
		Up: []string{`
      CREATE TABLE users (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        username varchar(255) NOT NULL,
        name varchar(255) DEFAULT NULL,
        email varchar(255) DEFAULT NULL,
        github_username varchar(255) DEFAULT NULL,
        github_token varchar(255) DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY uix_users_username (username)
      );
    `},
		Down: []string{`
      DROP TABLE IF EXISTS users;
    `},
	}
}
