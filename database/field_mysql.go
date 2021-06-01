// +build mysql

package database

import (
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	Service	string
	StorageMethod string
	Path string	`gorm:"type:varchar(500)"`
	OwnerID uint
	OwnerType string
}