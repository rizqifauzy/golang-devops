package repositories

import (
	"telegram-inventory/database"

	"gorm.io/gorm"
)

var DB *gorm.DB = database.Init()
