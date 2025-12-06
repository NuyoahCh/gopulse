# 🔒 安全配置指南

## ✅ 已配置的安全措施

### 1. `.gitignore` 配置

以下敏感文件和目录已被排除，**不会上传到 GitHub**：

#### 🔑 密钥和环境变量
- ✅ `.env` - 包含所有 API Keys
- ✅ `.env.local`
- ✅ `.env.*.local`
- ✅ `web/.env`

#### 📦 依赖目录（体积大）
- ✅ `web/node_modules/` - 前端依赖（通常几百 MB）
- ✅ `node_modules/`
- ✅ `vendor/` - Go 依赖

#### 🗄️ 数据库文件
- ✅ `*.db`
- ✅ `*.sqlite`
- ✅ `data/*.db`

#### 🔨 编译产物
- ✅ `bin/` - 编译后的二进制文件
- ✅ `*.exe`, `*.dll`, `*.so`

#### 💻 IDE 配置
- ✅ `.idea/` - JetBrains IDE
- ✅ `.vscode/` - VS Code
- ✅ `.DS_Store` - macOS

---

## ⚠️ 重要提醒

### `internal/config/config.go` 文件

**注意**：这个文件**会被上传**到 GitHub，因为它是源代码的一部分。

**但是**：
- ✅ 这个文件**不包含真实密钥**
- ✅ 它只是读取环境变量的代码
- ✅ 真实密钥在 `.env` 文件中（已被 `.gitignore` 排除）

#### 文件内容示例
```go
// config.go - 这个文件是安全的，可以上传
type Config struct {
    ArkAPIKey string  // 从环境变量读取，不是硬编码
    // ...
}

func Load() (*Config, error) {
    return &Config{
        ArkAPIKey: os.Getenv("ARK_API_KEY"),  // ✅ 安全：从环境变量读取
        // ...
    }, nil
}
```

**❌ 危险做法**（不要这样写）：
```go
// ❌ 永远不要在代码中硬编码密钥！
ArkAPIKey: "feabe6d9-8244-4e30-aff4-e7ad167a2ae9"
```

---

## 📋 上传前检查清单

### 1. 检查是否有敏感信息
```bash
# 检查 .env 是否被忽略
git check-ignore .env
# 应该输出：.env

# 检查 node_modules 是否被忽略
git check-ignore web/node_modules/
# 应该输出：web/node_modules/

# 查看将要提交的文件
git status
# 确保没有 .env 文件
```

### 2. 搜索代码中的硬编码密钥
```bash
# 搜索可能的 API Key
grep -r "sk-" --include="*.go" --include="*.js" .
grep -r "api_key.*=" --include="*.go" --include="*.js" .

# 应该没有任何硬编码的密钥
```

### 3. 验证 .gitignore 生效
```bash
# 添加所有文件
git add .

# 查看暂存区
git status

# 确认以下文件不在列表中：
# - .env
# - web/node_modules/
# - bin/
# - *.db
```

---

## 🚀 首次设置指南

### 新用户如何配置

1. **克隆仓库**
   ```bash
   git clone https://github.com/your-username/einoflow.git
   cd einoflow
   ```

2. **复制环境变量模板**
   ```bash
   cp .env.example .env
   ```

3. **编辑 `.env` 文件，填入真实密钥**
   ```bash
   # 使用你喜欢的编辑器
   vim .env
   # 或
   code .env
   ```

4. **配置示例**
   ```bash
   # 字节豆包配置
   ARK_API_KEY="你的真实密钥"
   ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
   
   # OpenAI 配置
   OPENAI_API_KEY="你的真实密钥"
   OPENAI_BASE_URL=https://api.openai.com/v1
   ```

5. **安装依赖**
   ```bash
   # Go 依赖
   go mod download
   
   # 前端依赖
   cd web
   npm install
   ```

---

## 🔐 密钥管理最佳实践

### ✅ 推荐做法

1. **使用环境变量**
   ```go
   apiKey := os.Getenv("ARK_API_KEY")
   ```

2. **使用 .env 文件（开发环境）**
   ```bash
   # .env
   ARK_API_KEY=your_key_here
   ```

3. **使用密钥管理服务（生产环境）**
   - AWS Secrets Manager
   - HashiCorp Vault
   - Azure Key Vault

4. **定期轮换密钥**
   - 每 90 天更换一次
   - 发现泄露立即更换

### ❌ 危险做法

1. **❌ 硬编码在代码中**
   ```go
   apiKey := "sk-abc123..."  // 永远不要这样做！
   ```

2. **❌ 提交到 Git**
   ```bash
   git add .env  # 危险！
   ```

3. **❌ 在日志中打印**
   ```go
   log.Printf("API Key: %s", apiKey)  // 危险！
   ```

4. **❌ 在错误信息中暴露**
   ```go
   return fmt.Errorf("failed with key %s", apiKey)  // 危险！
   ```

---

## 🆘 如果密钥已经泄露

### 立即行动

1. **撤销泄露的密钥**
   - 登录 API 提供商控制台
   - 立即删除或禁用泄露的密钥

2. **生成新密钥**
   - 创建新的 API Key
   - 更新本地 `.env` 文件

3. **清理 Git 历史**（如果已提交）
   ```bash
   # 使用 BFG Repo-Cleaner 或 git-filter-repo
   # 警告：这会重写历史，需要强制推送
   
   # 安装 BFG
   brew install bfg
   
   # 删除敏感文件
   bfg --delete-files .env
   
   # 清理
   git reflog expire --expire=now --all
   git gc --prune=now --aggressive
   
   # 强制推送（警告：会影响所有协作者）
   git push --force
   ```

4. **通知团队**
   - 告知所有协作者
   - 要求更新本地仓库

---

## 📊 安全检查脚本

创建一个自动检查脚本：

```bash
#!/bin/bash
# scripts/security-check.sh

echo "🔍 安全检查开始..."

# 检查 .env 是否被忽略
if git check-ignore .env > /dev/null 2>&1; then
    echo "✅ .env 已被 .gitignore 排除"
else
    echo "❌ 警告：.env 未被忽略！"
    exit 1
fi

# 检查是否有硬编码的密钥
if grep -r "sk-[a-zA-Z0-9]" --include="*.go" --include="*.js" . > /dev/null 2>&1; then
    echo "❌ 警告：发现可能的硬编码密钥！"
    grep -r "sk-[a-zA-Z0-9]" --include="*.go" --include="*.js" .
    exit 1
else
    echo "✅ 未发现硬编码密钥"
fi

# 检查 node_modules 是否被忽略
if git check-ignore web/node_modules/ > /dev/null 2>&1; then
    echo "✅ node_modules 已被 .gitignore 排除"
else
    echo "⚠️  警告：node_modules 未被忽略"
fi

echo "✅ 安全检查完成！"
```

使用方法：
```bash
chmod +x scripts/security-check.sh
./scripts/security-check.sh
```

---

## 📝 总结

### 当前配置状态

| 项目 | 状态 | 说明 |
|------|------|------|
| `.env` | ✅ 已保护 | 不会上传到 GitHub |
| `web/node_modules/` | ✅ 已保护 | 不会上传到 GitHub |
| `internal/config/config.go` | ✅ 安全 | 可以上传（不含密钥） |
| `.env.example` | ✅ 已提供 | 模板文件，可以上传 |
| `.gitignore` | ✅ 已配置 | 完整的忽略规则 |

### 下一步

1. ✅ 提交 `.gitignore` 文件
   ```bash
   git add .gitignore
   git commit -m "Add .gitignore to protect sensitive files"
   ```

2. ✅ 提交 `.env.example` 文件
   ```bash
   git add .env.example
   git commit -m "Add .env.example template"
   ```

3. ✅ 验证配置
   ```bash
   git status
   # 确保 .env 不在列表中
   ```

4. ✅ 推送到 GitHub
   ```bash
   git push
   ```

---

**记住**：密钥安全是第一位的！永远不要在代码中硬编码密钥，永远不要提交 `.env` 文件到 Git！🔒
