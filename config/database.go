package config

import (
	"log"
	"workshoptdd/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database!")
	}

	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		return nil
	}

	return db
}
