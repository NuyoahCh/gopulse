package config

import "os"

// Config 应用配置
type Config struct {
	Server ServerConfig
	Eino   EinoConfig
	Log    LogConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Host string
}

// EinoConfig Eino 配置
type EinoConfig struct {
	APIKey  string
	APIBase string
	Model   string
}

// LogConfig 日志配置
type LogConfig struct {
	Level string
}

// Load 加载配置
func Load() *Config {
	// 加载服务器配置
	return &Config{
		// 加载服务器配置
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		// 加载Eino配置
		Eino: EinoConfig{
			APIKey:  getEnv("ARK_API_KEY", ""),
			APIBase: getEnv("ARK_API_BASE_URL", "https://ark.cn-beijing.volces.com/api/v3"),
			Model:   getEnv("ARK_MODEL_NAME", "doubao-seed-1-6-lite-251015"),
		},
		// 加载日志配置
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
