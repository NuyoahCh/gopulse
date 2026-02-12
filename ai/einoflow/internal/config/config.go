package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// 服务配置
	ServerHost string
	ServerPort int
	LogLevel   string
	LogFormat  string

	// 数据库配置
	DBPath string

	// OpenAI 配置
	OpenAIKey     string
	OpenAIBaseURL string

	// 字节豆包配置
	ArkAPIKey         string
	ArkBaseURL        string
	ArkEmbeddingModel string

	// Anthropic 配置
	AnthropicKey string

	// 向量存储配置
	VectorStoreType string
	VectorDim       int
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:        getEnvInt("SERVER_PORT", 8080),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "json"),
		DBPath:            getEnv("DB_PATH", "./data/einoflow.db"),
		OpenAIKey:         getEnv("OPENAI_API_KEY", ""),
		OpenAIBaseURL:     getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		ArkAPIKey:         getEnv("ARK_API_KEY", ""), // 从 .env 文件读取，不设置默认值
		ArkBaseURL:        getEnv("ARK_BASE_URL", "https://ark.cn-beijing.volces.com/api/v3"),
		ArkEmbeddingModel: getEnv("ARK_EMBEDDING_MODEL", "doubao-embedding-large-text-250515"),
		AnthropicKey:      getEnv("ANTHROPIC_API_KEY", ""),
		VectorStoreType:   getEnv("VECTOR_STORE_TYPE", "memory"),
		VectorDim:         getEnvInt("VECTOR_DIM", 1536),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) Validate() error {
	// 验证至少配置了一个 LLM API Key
	if c.OpenAIKey == "" && c.ArkAPIKey == "" && c.AnthropicKey == "" {
		return fmt.Errorf("at least one LLM API key must be configured (OPENAI_API_KEY, ARK_API_KEY, or ANTHROPIC_API_KEY)")
	}

	// 验证服务器端口
	if c.ServerPort < 1024 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d (must be between 1024 and 65535)", c.ServerPort)
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be one of: debug, info, warn, error, fatal)", c.LogLevel)
	}

	// 验证日志格式
	if c.LogFormat != "json" && c.LogFormat != "text" {
		return fmt.Errorf("invalid log format: %s (must be 'json' or 'text')", c.LogFormat)
	}

	// 验证向量维度
	if c.VectorDim <= 0 || c.VectorDim > 10000 {
		return fmt.Errorf("invalid vector dimension: %d (must be between 1 and 10000)", c.VectorDim)
	}

	// 验证向量存储类型
	if c.VectorStoreType != "memory" && c.VectorStoreType != "persistent" {
		return fmt.Errorf("invalid vector store type: %s (must be 'memory' or 'persistent')", c.VectorStoreType)
	}

	// 验证数据库路径
	if c.DBPath == "" {
		return fmt.Errorf("database path cannot be empty")
	}

	return nil
}
