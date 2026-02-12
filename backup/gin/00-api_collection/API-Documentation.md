# Gin Framework Learning API 文档

本文档包含了完整的Gin框架学习项目的API接口说明。

## 导入方式

### Postman 导入
1. 打开 Postman
2. 点击左上角 "Import" 按钮
3. 选择 `Gin-API-Collection.postman_collection.json` 文件
4. 导入成功后即可在 Collections 中看到所有接口

### ApiPost 导入
1. 打开 ApiPost
2. 点击 "导入" → "Postman" 
3. 选择 `Gin-API-Collection.postman_collection.json` 文件
4. ApiPost 完全兼容 Postman Collection 格式

## 接口列表

### 1. QuickStart 模块 (端口: 8080)

#### GET `/`
- **描述**: 基础的Hello World示例
- **响应**: `hello world`

---

### 2. Compare 模块 (端口: 8080)

#### GET `/index`
- **描述**: 返回简单字符串响应
- **响应**: `Hello, World!`

#### GET `/get`
- **描述**: 返回JSON格式的标准响应
- **响应示例**:
```json
{
    "code": 0,
    "data": {},
    "msg": "成功"
}
```

#### POST `/post`
- **描述**: 接收JSON数据并返回
- **请求头**: `Content-Type: application/json`
- **请求体示例**:
```json
{
    "name": "张三",
    "age": 25,
    "city": "北京"
}
```
- **响应示例**:
```json
{
    "code": 0,
    "data": {
        "name": "张三",
        "age": 25,
        "city": "北京"
    },
    "msg": "成功"
}
```

---

### 3. Router 模块 (端口: 8080)

#### POST `/user/login`
- **描述**: 用户登录接口（Form表单格式）
- **请求头**: `Content-Type: application/x-www-form-urlencoded`
- **参数**:
  - `username`: 用户名（必填）
  - `password`: 密码（必填）
- **请求示例**:
```
username=admin&password=123456
```
- **响应**: `username=admin,password=123456`

---

### 4. Request 模块 (端口: 8080)

#### GET `/query`
- **描述**: 演示Query参数的多种获取方式
- **查询参数**:
  - `name`: 用户名称（可选，默认值: "shanyangsuanfa"）
  - `id`: 用户ID（必填）
- **请求示例**: `/query?name=张三&id=1001`
- **响应示例**:
```json
{
    "name": "张三",
    "name_with_default": "张三",
    "id": "1001"
}
```

---

### 5. Model 模块 (端口: 8080)

#### GET `/posts/index`
- **描述**: 渲染posts模板页面
- **响应**: HTML页面

#### GET `/users/index`
- **描述**: 渲染users模板页面
- **响应**: HTML页面

---

### 6. Bind 模块

#### POST `/login` (端口: 9090)
- **描述**: JSON格式的登录请求
- **请求头**: `Content-Type: application/json`
- **请求体**:
```json
{
    "username": "admin",
    "password": "123456"
}
```
- **验证规则**: username和password均为必填项

#### GET `/getb` (端口: 8080)
- **描述**: 嵌套结构体参数绑定
- **查询参数**:
  - `field_a`: 嵌套结构体字段
  - `field_b`: 主结构体字段
- **请求示例**: `/getb?field_a=valueA&field_b=valueB`
- **响应示例**:
```json
{
    "a": {
        "field_a": "valueA"
    },
    "b": "valueB"
}
```

#### GET `/getc` (端口: 8080)
- **描述**: 包含结构体指针的参数绑定
- **查询参数**:
  - `field_a`: 嵌套结构体指针字段
  - `field_c`: 主结构体字段
- **请求示例**: `/getc?field_a=valueA&field_c=valueC`

#### GET `/getd` (端口: 8080)
- **描述**: 匿名结构体的参数绑定
- **查询参数**:
  - `field_x`: 匿名结构体字段
  - `field_d`: 主结构体字段
- **请求示例**: `/getd?field_x=valueX&field_d=valueD`

---

### 7. Middleware 模块 (端口: 8080)

#### GET `/ping`
- **描述**: 带延迟监控中间件的Ping接口
- **功能**: 
  - 模拟100ms的处理时间
  - 中间件会在服务端日志输出请求耗时、状态码和路径
- **响应示例**:
```json
{
    "message": "pong"
}
```
- **服务端日志输出示例**:
```
[Metrics] | 200 |   100.123456ms | /ping
```

---

### 8. Sessions 模块 (端口: 8080)

#### GET `/login`
- **描述**: Session登录示例
- **依赖**: 需要Redis服务运行在 localhost:6379
- **响应示例**:
```json
{
    "status": "logged in"
}
```

---

## 运行不同模块

每个模块可以独立运行，进入对应目录执行：

```bash
# QuickStart
cd 01-quickstart && go run main.go

# Compare
cd 02-compare/gin && go run main.go

# Router
cd 03-router && go run login.go

# Request
cd 04-request && go run get.go

# Model
cd 05-model && go run main.go

# Bind (注意端口9090)
cd 06-bind && go run login.go

# Middleware
cd 07-middleware && go run monitor.go

# Sessions (需要先启动Redis)
cd 08-sessions && go run sessions.go
```

## 注意事项

1. **端口说明**:
   - 大部分接口运行在 `8080` 端口
   - `06-bind/login.go` 运行在 `9090` 端口

2. **依赖服务**:
   - Sessions模块需要Redis服务: `localhost:6379`

3. **模板文件**:
   - Model模块需要 `templates/` 目录下的模板文件

4. **测试建议**:
   - 建议按模块逐个启动测试
   - 注意不同模块使用的端口
   - 某些接口需要特定的请求格式（JSON/Form）

## 快速测试流程

1. 启动某个模块的服务
2. 在Postman/ApiPost中选择对应的文件夹
3. 确认请求的端口与运行的服务匹配
4. 发送请求进行测试

## 学习路径建议

1. **01-quickstart**: Gin基础入门
2. **02-compare**: 理解Gin相比原生Go的优势
3. **03-router**: 路由和POST请求
7. **07-middleware**: 中间件和请求监控
5. **05-model**: 模板渲染
6. **06-bind**: 数据绑定和验证
7. **07-middleware**: 中间件（仅代码示例）
8. **08-sessions**: Session管理

---

**版本**: v1.0
