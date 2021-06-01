// +build mysql

package database

import (
	em "github.com/Etpmls/Etpmls-Micro/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "LIKE"
)

func (this *database) runDatabase() {
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=" + url.QueryEscape(timezone)

	//Connect Database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
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