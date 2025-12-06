package tools

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseTool 数据库查询工具
type DatabaseTool struct {
	db *gorm.DB
}

// NewDatabaseTool 创建数据库工具
func NewDatabaseTool(dbPath string) (*DatabaseTool, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return &DatabaseTool{db: db}, nil
}

// QueryData 查询数据
// @tool Query data from database
func (t *DatabaseTool) QueryData(ctx context.Context, table string, condition string) (string, error) {
	var results []map[string]interface{}

	query := t.db.Table(table)
	if condition != "" {
		query = query.Where(condition)
	}

	if err := query.Limit(10).Find(&results).Error; err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}

	if len(results) == 0 {
		return "No data found", nil
	}

	// 格式化结果
	output := fmt.Sprintf("Found %d records:\n", len(results))
	for i, row := range results {
		output += fmt.Sprintf("%d. %v\n", i+1, row)
	}

	return output, nil
}

// InsertData 插入数据
// @tool Insert data into database
func (t *DatabaseTool) InsertData(ctx context.Context, table string, data map[string]interface{}) (string, error) {
	if err := t.db.Table(table).Create(data).Error; err != nil {
		return "", fmt.Errorf("insert failed: %w", err)
	}

	return "Data inserted successfully", nil
}
