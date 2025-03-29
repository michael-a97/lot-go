package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"lot/config"
	"lot/internal/entity"
)

func ConnectDb() *gorm.DB {
	dbUserName, err := config.Config("dbUsername")
	if err != nil {
		log.Fatal("Please sepecify a `dbUsername`")
	}
	dbPassword, err := config.Config("dbPassword")
	if err != nil {
		log.Fatal("Please sepecify a `dbPassword`")
	}
	host, err := config.Config("host")
	if err != nil {
		log.Fatal("Please sepecify a `host`")
	}
	dbPort, err := config.Config("dbPort")
	if err != nil {
		log.Fatal("Please sepecify a `dbPort`")
	}
	dbName, err := config.Config("dbName")
	if err != nil {
		log.Fatal("Please sepecify a `dbName`")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUserName,
		dbPassword,
		host,
		dbPort,
		dbName,
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
