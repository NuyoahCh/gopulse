.PHONY: help build run test clean install demo

help:
	@echo "EinoFlow - 基于 Eino 的 AI 应用平台"
	@echo ""
	@echo "可用命令:"
	@echo "  make install    - 安装依赖"
	@echo "  make build      - 编译项目"
	@echo "  make run        - 运行服务器"
	@echo "  make demo       - 运行演示程序"
	@echo "  make test       - 运行测试"
	@echo "  make clean      - 清理构建文件"

install:
	@echo "安装依赖..."
	go mod download
	go mod tidy

build:
	@echo "编译项目..."
	go build -o bin/server cmd/server/main.go
	go build -o bin/demo examples/complete_demo.go

run:
	@echo "启动服务器..."
	go run cmd/server/main.go

demo:
	@echo "运行演示程序..."
	go run examples/complete_demo.go

test:
	@echo "运行测试..."
	go test -v ./...

clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f einoflow

# 快速测试各个功能
test-llm:
	@echo "测试 LLM 功能..."
	curl -X POST http://localhost:8080/api/v1/llm/chat \
		-H "Content-Type: application/json" \
		-d '{"model":"ep-20241116153014-gfmhp","messages":[{"role":"user","content":"你好"}]}'

test-stream:
	@echo "测试流式响应..."
	curl -X POST http://localhost:8080/api/v1/llm/chat/stream \
		-H "Content-Type: application/json" \
		-d '{"model":"ep-20241116153014-gfmhp","messages":[{"role":"user","content":"讲个故事"}]}'

test-models:
	@echo "获取模型列表..."
	curl http://localhost:8080/api/v1/llm/models

# 生成 Swagger 文档
swagger:
	@echo "生成 API 文档..."
	swag init -g cmd/server/main.go --output docs
	@echo "✅ 文档已生成，启动服务后访问 http://localhost:8080/swagger/index.html"