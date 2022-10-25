package database

import (
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
	"fmt"
)

var Db *gorm.DB
func InitDb() *gorm.DB { // OOP constructor
	Db = connectDB()
	return Db
}

func connectDB() (*gorm.DB) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err !=nil {
		fmt.Println("Error...")
		return nil
	}
	return db
}
