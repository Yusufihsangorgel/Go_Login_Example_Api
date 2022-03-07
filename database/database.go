package database

import (
	"backendtest/models"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func Connect() *gorm.DB {

	//database bağlantısını oluşturuyoruz ve gorm.DB tipinde bir değişkene atıyoruz
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic(err)

	}

	//veritabanına tablo oluşturuyoruz
	db.AutoMigrate(&models.User{})

	return db

}
