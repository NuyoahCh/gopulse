# 快速测试指南

## 导入步骤

### Postman
```bash
File → Import → 选择 Gin-API-Collection.postman_collection.json
```

### ApiPost/Other  
```bash
导入 → Postman → 选择 Gin-API-Collection.postman_collection.json
```

## 测试检查清单

### 基础测试 (8080端口)
- [ ] GET `/` - Hello World
- [ ] GET `/index` - 字符串响应
- [ ] GET `/get` - JSON响应
- [ ] POST `/post` - JSON请求与响应

### 表单提交 (8080端口)
- [ ] POST `/user/login` - Form表单登录

### 参数处理 (8080端口)
- [ ] GET `/query?name=test&id=123` - Query参数
- [ ] GET `/query?id=123` - 默认参数测试

### 数据绑定 (8080/9090端口)
- [ ] POST `/login` (9090端口) - JSON绑定
- [ ] GET `/getb?field_a=A&field_b=B` - 嵌套结构体
- [ ] GET `/getc?field_a=A&field_c=C` - 指针绑定
- [ ] GET `/getd?field_x=X&field_d=D` - 匿名结构体

### 中间件测试 (8080端口)
- [ ] GET `/ping` - 延迟监控中间件（查看服务端日志）

### Session测试 (8080端口，需要Redis)
- [ ] GET `/login` - Session存储

## 快速启动命令

```bash
# 测试基础功能
cd 01-quickstart && go run main.go

# 测试JSON和POST
cd 02-compare/gin && go run main.go

# 测试Form登录
cd 03-router && go run login.go

# 测试Query参数
cd 04-request && go run get.go

# 测试JSON绑定（端口9090）
cd 06-bind && go run login.go

# 测试中间件（查看日志输出）
cd 07-middleware && go run monitor.go

# 测试Session（需要Redis）
cd 08-sessions && go run sessions.go
```

## 常见问题

1. **端口冲突**: 确保8080/9090端口未被占用
2. **Redis未启动**: Session模块需要先启动Redis
3. **Import失败**: 确认Postman版本支持v2.1格式
4. **请求失败**: 检查服务是否已启动，端口是否正确

## 推荐测试顺序

1. 启动 `02-compare/gin` → 测试基础CRUD
2. 启动 `04-request` → 测试参数获取
3. 启动 `06-bind/login` → 测试JSON绑定
4. 启动 `07-middleware` → 测试中间件（查看日志输出）
5. （可选）启动Redis → 测试Session
