# EinoVerse API 文档

企业员工内部知识库和请假申请审批系统 API 文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: 统一使用 JSON 格式

### 统一响应格式

#### 成功响应
```json
{
    "code": 0,
    "msg": "success",
    "data": {}
}
```

#### 错误响应
```json
{
    "code": "ERROR_CODE",
    "msg": "错误描述信息"
}
```

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| `0` | 成功 |
| `DOC_NOT_FOUND` | 文档未找到 |
| `INVALID_INPUT` | 输入参数无效 |
| `LLM_ERROR` | LLM服务调用失败 |
| `APP_NOT_FOUND` | 申请未找到 |
| `INTERNAL_ERROR` | 内部服务器错误 |

---

## 1. 健康检查

### GET /health

检查服务是否正常运行。

#### 请求示例
```bash
GET http://localhost:8080/health
```

#### 响应示例
```json
{
    "status": "ok",
    "model_backend": "eino" // 或 "mock"
}
```

---

## 2. 知识库管理

### 2.1 创建文档

#### POST /api/v1/knowledgebase/documents

创建新的知识库文档。

#### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `title` | string | 是 | 文档标题 |
| `content` | string | 是 | 文档内容 |
| `tags` | string[] | 否 | 标签列表 |
| `author` | string | 否 | 作者 |

#### 请求示例
```json
{
    "title": "员工手册",
    "content": "这是员工手册的内容，包含公司规章制度、福利待遇等信息。\n\n1. 工作时间：周一至周五 9:00-18:00\n2. 请假制度：需要提前申请，经主管批准后方可生效。\n3. 年假政策：入职满一年后可享受5天年假。",
    "tags": ["员工手册", "规章制度"],
    "author": "HR部门"
}
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "doc_123456789"
    }
}
```

---

### 2.2 获取文档

#### GET /api/v1/knowledgebase/documents/:id

根据ID获取文档详情。

#### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 文档ID |

#### 请求示例
```bash
GET http://localhost:8080/api/v1/knowledgebase/documents/doc_123456789
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "doc_123456789",
        "title": "员工手册",
        "content": "这是员工手册的内容...",
        "tags": ["员工手册", "规章制度"],
        "author": "HR部门",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
    }
}
```

---

### 2.3 搜索文档

#### GET /api/v1/knowledgebase/documents/search

根据关键词搜索文档。

#### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `q` | string | 是 | 搜索关键词 |
| `limit` | int | 否 | 返回结果数量，默认5，最大100 |

#### 请求示例
```bash
GET http://localhost:8080/api/v1/knowledgebase/documents/search?q=请假&limit=5
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": [
        {
            "document": {
                "id": "doc_123456789",
                "title": "员工手册",
                "content": "...",
                "tags": ["员工手册", "规章制度"],
                "created_at": "2024-01-15T10:30:00Z",
                "updated_at": "2024-01-15T10:30:00Z"
            },
            "score": 10.0
        }
    ]
}
```

---

### 2.4 问答

#### POST /api/v1/knowledgebase/ask

基于知识库内容进行智能问答。**需要配置 Eino API 才能使用**。

#### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `question` | string | 是 | 问题 |

#### 请求示例
```json
{
    "question": "请问公司年假政策是什么？"
}
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "answer": "根据员工手册，公司年假政策为：入职满一年后可享受5天年假。",
        "source_docs": ["doc_123456789"],
        "confidence": 0.95
    }
}
```

---

## 3. 请假管理

### 3.1 创建请假申请

#### POST /api/v1/leave/applications

通过自然语言文本自动创建结构化请假申请。**需要配置 Eino API 才能使用**。

#### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `employee_id` | string | 是 | 员工ID |
| `employee_name` | string | 是 | 员工姓名 |
| `supervisor` | string | 是 | 直属主管 |
| `text` | string | 是 | 请假文本描述 |

#### 请求示例
```json
{
    "employee_id": "E001",
    "employee_name": "张三",
    "supervisor": "李主管",
    "text": "本周五因家中有事需要请假一天，工作已交接给同事王五处理，预计下周一返回工作岗位。"
}
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "app_123456789",
        "employee_id": "E001",
        "employee_name": "张三",
        "supervisor": "李主管",
        "leave_type": "personal",
        "start_date": "2024-01-19T00:00:00Z",
        "end_date": "2024-01-19T23:59:59Z",
        "days": 1.0,
        "reason": "家中有事",
        "work_handover": "工作已交接给同事王五处理",
        "status": "pending",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
    }
}
```

#### 请假类型说明

| 类型值 | 说明 |
|--------|------|
| `sick` | 病假 |
| `personal` | 事假 |
| `annual` | 年假 |
| `adjust` | 调休 |

#### 申请状态说明

| 状态值 | 说明 |
|--------|------|
| `pending` | 待审批 |
| `approved` | 已批准 |
| `rejected` | 已拒绝 |

---

### 3.2 获取请假申请

#### GET /api/v1/leave/applications/:id

根据ID获取请假申请详情。

#### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 申请ID |

#### 请求示例
```bash
GET http://localhost:8080/api/v1/leave/applications/app_123456789
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "id": "app_123456789",
        "employee_id": "E001",
        "employee_name": "张三",
        "supervisor": "李主管",
        "leave_type": "personal",
        "start_date": "2024-01-19T00:00:00Z",
        "end_date": "2024-01-19T23:59:59Z",
        "days": 1.0,
        "reason": "家中有事",
        "work_handover": "工作已交接给同事王五处理",
        "status": "pending",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
    }
}
```

---

### 3.3 审批请假申请

#### POST /api/v1/leave/applications/:id/approve

对请假申请进行审批并获取AI建议。**需要配置 Eino API 才能使用**。

#### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 申请ID |

#### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `approver` | string | 是 | 审批人姓名 |
| `comments` | string | 否 | 审批备注 |

**注意**: 申请ID从URL路径参数中获取，不需要在请求体中提供。

#### 请求示例
```json
{
    "approver": "李主管",
    "comments": "同意"
}
```

**完整URL示例**:
```
POST http://localhost:8080/api/v1/leave/applications/app_123456789/approve
```

#### 响应示例
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "decision": "approved",
        "confidence": 0.85,
        "reason": "请假时间合理，工作已交接，建议批准",
        "suggestions": [
            "建议批准本次请假申请",
            "确认工作交接已完成"
        ]
    }
}
```

#### 审批决策说明

| 决策值 | 说明 |
|--------|------|
| `approved` | 已批准 |
| `rejected` | 已拒绝 |
| `conditional` | 条件批准（待补充材料等） |

---

## 环境配置

### 必需环境变量

```bash
# Eino API 配置（可选，不配置将使用 Mock 模式）
export EINO_API_KEY=your_api_key
export EINO_API_BASE=https://api.eino.bytedance.com/v1
export EINO_MODEL=doubao-pro-128k

# 服务器配置
export PORT=8080
export HOST=0.0.0.0

# 日志级别
export LOG_LEVEL=info  # debug, info, warn, error
```

### 功能说明

- **有 Eino API Key**: 所有功能正常工作，包括智能问答和请假申请生成
- **无 Eino API Key**: 使用 Mock 模式，LLM 相关功能受限，其他功能正常

---

## 测试示例

### 完整流程示例

#### 1. 健康检查
```bash
curl -X GET http://localhost:8080/health
```

#### 2. 创建知识库文档
```bash
curl -X POST http://localhost:8080/api/v1/knowledgebase/documents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "请假制度",
    "content": "员工请假需提前申请，事假需提前3天申请，病假可当天申请。",
    "tags": ["请假", "制度"],
    "author": "HR"
  }'
```

#### 3. 搜索文档
```bash
curl -X GET "http://localhost:8080/api/v1/knowledgebase/documents/search?q=请假&limit=5"
```

#### 4. 问答
```bash
curl -X POST http://localhost:8080/api/v1/knowledgebase/ask \
  -H "Content-Type: application/json" \
  -d '{
    "question": "请假需要提前多久申请？"
  }'
```

#### 5. 创建请假申请
```bash
curl -X POST http://localhost:8080/api/v1/leave/applications \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "E001",
    "employee_name": "张三",
    "supervisor": "李主管",
    "text": "明天因感冒需要请假一天，工作已交接给王五。"
  }'
```

#### 6. 获取请假申请（使用上一步返回的ID）
```bash
curl -X GET http://localhost:8080/api/v1/leave/applications/app_123456789
```

#### 7. 审批请假申请
```bash
curl -X POST http://localhost:8080/api/v1/leave/applications/app_123456789/approve \
  -H "Content-Type: application/json" \
  -d '{
    "approver": "李主管",
    "comments": "同意"
  }'
```

---

## 注意事项

1. **日期格式**: 所有日期字段使用 ISO 8601 格式（如：`2024-01-19T00:00:00Z`）
2. **时间计算**: 请假天数按自然日计算，包含起始和结束日期
3. **LLM 功能**: 需要配置 Eino API Key 才能使用智能问答和请假申请生成功能
4. **数据持久化**: 当前版本使用内存存储，服务重启后数据会丢失
5. **并发安全**: 所有操作都进行了并发安全处理

---

## 更新日志

- **v1.0.0** (2024-01-15)
  - 初始版本发布
  - 支持知识库文档管理
  - 支持智能问答
  - 支持请假申请创建和审批

