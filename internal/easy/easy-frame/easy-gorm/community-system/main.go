package main

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;size:50;not null"`
	Email     string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	log.Println("=== GORM 最小可运行示例 ===\\n")

	// 1. 连接数据库
	dsn := "root:123456@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ 连接数据库失败:", err)
	}
	log.Println("✅ 数据库连接成功")

	// 2. 自动迁移（创建表）
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("❌ 自动迁移失败:", err)
	}
	log.Println("✅ 表结构创建成功")

	// 3. Create - 创建用户
	user := User{
		Username: "alice",
		Email:    "alice@example.com",
	}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatal("❌ 创建用户失败:", result.Error)
	}
	log.Printf("✅ 创建用户成功: %s (ID: %d)\\n", user.Username, user.ID)

	// 4. Read - 查询用户
	var foundUser User
	db.First(&foundUser, user.ID)
	log.Printf("✅ 查询用户: %s - %s\\n", foundUser.Username, foundUser.Email)

	// 5. Update - 更新用户
	foundUser.Email = "alice.new@example.com"
	db.Save(&foundUser)
	log.Printf("✅ 更新用户邮箱: %s\\n", foundUser.Email)

	// 6. Delete - 删除用户
	db.Delete(&foundUser)
	log.Printf("✅ 删除用户: %s\\n", foundUser.Username)

	// 验证删除
	var deletedUser User
	result = db.First(&deletedUser, user.ID)
	if result.Error == gorm.ErrRecordNotFound {
		log.Println("✅ 用户已成功删除")
	}
}
