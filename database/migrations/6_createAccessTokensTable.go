package migrations

import (
	"github.com/rubenv/sql-migrate"
)

func createAccessTokensTable() *migrate.Migration {
	return &migrate.Migration{
		Id: "6",
		Up: []string{`
      CREATE TABLE access_tokens (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        token varchar(128) NOT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY uix_access_tokens_token (token)
      );
    `},
		Down: []string{`
      DROP TABLE IF EXISTS access_tokens;
    `},
	}
}
