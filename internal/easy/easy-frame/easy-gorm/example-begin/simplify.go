package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 用户模型
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;size:50;not null"`
	Email    string `gorm:"size:100;not null"`
}

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=True&loc=Local"
	// CREATE DATABASE IF NOT EXISTS community CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 自动迁移（创建表）
	db.AutoMigrate(&User{})

	// 创建用户 - 简单直观
	user := User{Username: "alice", Email: "alice@example.com"}
	db.Create(&user)
	log.Printf("创建用户: %s (ID: %d)", user.Username, user.ID)

	// 查询用户 - 链式调用
	var foundUser User
	db.Where("username = ?", "alice").First(&foundUser)
	log.Printf("找到用户: %s", foundUser.Email)
}
