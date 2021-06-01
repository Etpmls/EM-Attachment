// +build postgresql

package database

import (
	em "github.com/Etpmls/Etpmls-Micro/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "ILIKE"
)

func (this *database) runDatabase() {

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=" + timezone

	//Connect Database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: prefix,
		},
	})
	if err != nil {
		em.LogPanic.Path("Unable to connect to the database!", err)
	}

	err = DB.AutoMigrate(migrate...)
	if err != nil {
		em.LogInfo.Path("Failed to create database!", err)
	}

}