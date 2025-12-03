package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 用户模型
type UserOther struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;size:50;not null"`
	Email    string `gorm:"size:100;not null"`
}

func main() {
	log.Println("=== GORM 简介与特点 ===\\n")

	dsn := "root:123456@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 自动迁移
	db.AutoMigrate(&UserOther{})

	// 创建用户
	user := UserOther{Username: "alice", Email: "alice@example.com"}
	db.Create(&user)
	log.Printf("✅ 创建用户: %s", user.Username)

	// 查询用户
	var foundUser UserOther
	db.First(&foundUser, user.ID)
	log.Printf("✅ 查询用户: %s", foundUser.Email)
}
