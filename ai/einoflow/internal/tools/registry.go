package tools

// Registry 工具注册表
type Registry struct {
	tools map[string]interface{}
}

// NewRegistry 创建工具注册表
func NewRegistry(dbPath string, fileBaseDir string) *Registry {
	r := &Registry{
		tools: make(map[string]interface{}),
	}

	// 注册默认工具
	r.Register("weather", NewWeatherTool())
	r.Register("calculator", NewCalculatorTool())
	r.Register("search", NewSearchTool("")) // 使用 DuckDuckGo，不需要 API Key
	r.Register("file", NewFileTool(fileBaseDir))

	// 数据库工具（可选）
	if dbPath != "" {
		if dbTool, err := NewDatabaseTool(dbPath); err == nil {
			r.Register("database", dbTool)
		}
	}

	return r
}

// Register 注册工具
func (r *Registry) Register(name string, t interface{}) {
	r.tools[name] = t
}

// Get 获取工具
func (r *Registry) Get(name string) (interface{}, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// GetAll 获取所有工具
func (r *Registry) GetAll() []interface{} {
	tools := make([]interface{}, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// GetByNames 根据名称获取工具列表
func (r *Registry) GetByNames(names []string) []interface{} {
	tools := make([]interface{}, 0, len(names))
	for _, name := range names {
		if t, ok := r.tools[name]; ok {
			tools = append(tools, t)
		}
	}
	return tools
}

// ListToolNames 列出所有工具名称
func (r *Registry) ListToolNames() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}
