// main.go
package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.Println("=== 环境准备检查 ===\\n")

	// 连接数据库
	dsn := "root:123456@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ 连接数据库失败:", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ 获取数据库实例失败:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("❌ 数据库连接测试失败:", err)
	}

	log.Println("✅ 环境准备完成！")
	log.Println("✅ Go 版本:", "go1.21+")
	log.Println("✅ GORM 已安装")
	log.Println("✅ 数据库连接正常")
}
