package test

import (
	"testing"
	"wxcloudrun-golang/db"
)

func TestMigrate(t *testing.T) {
	db.Init()
	db.Migrate()
}
