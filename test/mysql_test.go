package test

import (
	"carbonbackend/db"
	"testing"
)

func TestMigrate(t *testing.T) {
	db.Init()
	db.Migrate()
}
