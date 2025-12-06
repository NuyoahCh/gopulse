package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// FileTool 文件操作工具
type FileTool struct {
	baseDir string
}

// NewFileTool 创建文件工具
func NewFileTool(baseDir string) *FileTool {
	return &FileTool{baseDir: baseDir}
}

// ReadFile 读取文件
// @tool Read content from a file
func (t *FileTool) ReadFile(ctx context.Context, filename string) (string, error) {
	path := filepath.Join(t.baseDir, filename)

	// 安全检查：确保路径在 baseDir 内
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absBaseDir, err := filepath.Abs(t.baseDir)
	if err != nil {
		return "", err
	}
	if !filepath.HasPrefix(absPath, absBaseDir) {
		return "", fmt.Errorf("access denied: path outside base directory")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// WriteFile 写入文件
// @tool Write content to a file
func (t *FileTool) WriteFile(ctx context.Context, filename string, content string) (string, error) {
	path := filepath.Join(t.baseDir, filename)

	// 安全检查
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absBaseDir, err := filepath.Abs(t.baseDir)
	if err != nil {
		return "", err
	}
	if !filepath.HasPrefix(absPath, absBaseDir) {
		return "", fmt.Errorf("access denied: path outside base directory")
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("File written successfully: %s", filename), nil
}

// ListFiles 列出文件
// @tool List files in a directory
func (t *FileTool) ListFiles(ctx context.Context, dir string) (string, error) {
	path := filepath.Join(t.baseDir, dir)

	// 安全检查
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absBaseDir, err := filepath.Abs(t.baseDir)
	if err != nil {
		return "", err
	}
	if !filepath.HasPrefix(absPath, absBaseDir) {
		return "", fmt.Errorf("access denied: path outside base directory")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	result := fmt.Sprintf("Files in %s:\n", dir)
	for _, entry := range entries {
		if entry.IsDir() {
			result += fmt.Sprintf("  [DIR]  %s\n", entry.Name())
		} else {
			info, _ := entry.Info()
			result += fmt.Sprintf("  [FILE] %s (%d bytes)\n", entry.Name(), info.Size())
		}
	}

	return result, nil
}
