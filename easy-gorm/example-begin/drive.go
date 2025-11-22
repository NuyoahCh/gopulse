package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type User01 struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;size:50;not null"`
	Email    string `gorm:"size:100;not null"`
}

func main() {
	// MySQL DSN 格式
	dsn := "root:123456@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User01{})
	log.Println("✅ MySQL 连接成功")
}
