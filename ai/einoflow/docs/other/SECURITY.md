# 安全配置指南

## 🔐 API Key 管理最佳实践

### ✅ 正确做法

1. **使用环境变量**
   ```bash
   # 在 .env 文件中配置
   ARK_API_KEY=your_real_api_key_here
   ```

2. **确保 .env 文件被 gitignore**
   ```bash
   # 检查 .gitignore 中是否包含
   grep "^\.env$" .gitignore
   ```

3. **使用 .env.example 作为模板**
   ```bash
   # 复制示例文件
   cp .env.example .env
   # 然后编辑 .env 填入真实密钥
   ```

### ❌ 错误做法

1. **永远不要在代码中硬编码密钥**
   ```go
   // ❌ 错误 - 硬编码密钥
   ArkAPIKey: "feabe6d9-8244-4e30-aff4-e7ad167a2ae9"
   
   // ✅ 正确 - 从环境变量读取
   ArkAPIKey: getEnv("ARK_API_KEY", "")
   ```

2. **不要将 .env 文件提交到 Git**
   ```bash
   # 如果不小心提交了，立即移除
   git rm --cached .env
   git commit -m "Remove .env file"
   ```

3. **不要在日志中打印密钥**
   ```go
   // ❌ 错误
   logger.Info("API Key: " + cfg.ArkAPIKey)
   
   // ✅ 正确
   logger.Info("API Key configured: " + (cfg.ArkAPIKey != ""))
   ```

## 🛡️ 配置文件安全检查

### 检查清单

- [ ] `.env` 文件在 `.gitignore` 中
- [ ] `.env` 文件未被 Git 跟踪
- [ ] 代码中没有硬编码的 API Key
- [ ] `.env.example` 中只有占位符，没有真实密钥
- [ ] 日志中不打印敏感信息

### 自动检查

运行以下命令检查是否有安全问题：

```bash
# 检查 .env 是否被跟踪
git ls-files | grep "^\.env$"
# 如果有输出，说明 .env 被跟踪了，需要移除

# 检查代码中是否有可疑的硬编码
grep -r "API_KEY.*=.*\"[^y]" internal/ cmd/ pkg/
# 如果有输出，需要检查是否是硬编码的密钥
```

## 🔒 生产环境部署

### Docker 部署

使用环境变量或 secrets 管理：

```yaml
# docker-compose.yml
services:
  einoflow:
    image: einoflow:latest
    environment:
      - ARK_API_KEY=${ARK_API_KEY}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    # 或使用 secrets
    secrets:
      - ark_api_key
```

### Kubernetes 部署

使用 Kubernetes Secrets：

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: einoflow-secrets
type: Opaque
stringData:
  ark-api-key: your_real_api_key_here
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: einoflow
spec:
  template:
    spec:
      containers:
      - name: einoflow
        env:
        - name: ARK_API_KEY
          valueFrom:
            secretKeyRef:
              name: einoflow-secrets
              key: ark-api-key
```

## 📋 密钥轮换

定期轮换 API Key：

1. 在 LLM 提供商控制台生成新密钥
2. 更新 `.env` 文件
3. 重启服务
4. 删除旧密钥

## 🚨 密钥泄漏应急处理

如果不小心将密钥提交到 Git：

1. **立即撤销密钥**
   - 登录 LLM 提供商控制台
   - 撤销泄漏的 API Key

2. **从 Git 历史中移除**
   ```bash
   # 使用 git filter-branch 或 BFG Repo-Cleaner
   git filter-branch --force --index-filter \
     "git rm --cached --ignore-unmatch .env" \
     --prune-empty --tag-name-filter cat -- --all
   
   # 强制推送
   git push origin --force --all
   ```

3. **生成新密钥**
   - 生成新的 API Key
   - 更新 `.env` 文件
   - 通知团队成员更新

## 📚 相关资源

- [字节豆包 API 文档](https://www.volcengine.com/docs/82379)
- [OpenAI API 安全最佳实践](https://platform.openai.com/docs/guides/safety-best-practices)
- [OWASP API 安全指南](https://owasp.org/www-project-api-security/)
