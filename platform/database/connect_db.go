package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"lot/config"
	"lot/pkg/entity"
)

func ConnectDb() *gorm.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Config("dbUsername"),
		config.Config("dbPassword"),
		config.Config("host"),
		config.Config("dbPort"),
		config.Config("dbName"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the db")
	} else {
		fmt.Println("Connected to the db")
	}
	db.AutoMigrate(entity.Role{}, entity.User{}, entity.RefreshToken{})
	return db
}
