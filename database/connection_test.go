package database

import (
	"testing"

	"github.com/generationtux/brizo/config"
	"github.com/stretchr/testify/assert"
)

func TestGetDBAddress(t *testing.T) {
	config.App.MysqlUser = "foo"
	config.App.MysqlPassword = "bar"
	config.App.MysqlHost = "db.example.com"
	config.App.MysqlPort = "3307"
	config.App.MysqlDatabase = "somedb"

	expected := "foo:bar@(db.example.com:3307)/somedb?charset=utf8&parseTime=True&loc=Local"
	actual := getDBAddress()
	assert.Equal(t, expected, actual)
}
func TestGetDBConnectionError(t *testing.T) {
	config.App.MysqlUser = "foo"
	config.App.MysqlPassword = "bar"
	config.App.MysqlHost = "db.example.com"
	config.App.MysqlPort = "3307"
	config.App.MysqlDatabase = "somedb"

	_, err := Connect()
	assert.Contains(t, err.Error(), "no such host")
}
