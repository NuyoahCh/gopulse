# 🚀 安全提交指南

## ✅ 当前状态

### 已完成的安全配置

1. ✅ **`.gitignore` 已创建** - 保护敏感文件
2. ✅ **`.env` 已从 git 删除** - 密钥不会上传
3. ✅ **硬编码密钥已移除** - 示例代码使用环境变量
4. ✅ **安全检查脚本已创建** - 自动检测安全问题
5. ✅ **安全检查通过** ✨

---

## 📋 待提交的文件

### 当前暂存区
```bash
要提交的变更：
  删除：     .env              ✅ 正确：从 git 中移除密钥文件
  新文件：   .gitignore        ✅ 正确：添加忽略规则
```

### 未暂存的修改
```bash
  修改：     examples/agent/weather_agent.go    ✅ 已移除硬编码密钥
  修改：     examples/llm/basic_chat.go         ✅ 已移除硬编码密钥
  修改：     web/README.md                      ℹ️  其他修改
```

### 未跟踪的文件
```bash
  SECURITY.md                  ✅ 安全文档
  scripts/security-check.sh    ✅ 安全检查脚本
  web/package-lock.json        ℹ️  前端依赖锁文件
```

---

## 🎯 推荐的提交步骤

### 步骤 1: 提交安全配置（最重要）

```bash
# 提交 .gitignore 和删除 .env
git commit -m "security: Add .gitignore and remove .env from git

- Add comprehensive .gitignore rules
- Remove .env file from git tracking
- Protect sensitive files (API keys, node_modules, etc.)
- Ensure .env.example is available as template"
```

### 步骤 2: 提交代码修复

```bash
# 添加修改后的示例代码
git add examples/agent/weather_agent.go examples/llm/basic_chat.go

# 提交
git commit -m "security: Remove hardcoded API keys from examples

- Replace hardcoded API keys with environment variables
- Add validation for required environment variables
- Update examples to use os.Getenv()"
```

### 步骤 3: 提交安全文档和工具

```bash
# 添加安全相关文件
git add SECURITY.md scripts/security-check.sh

# 提交
git commit -m "docs: Add security documentation and check script

- Add SECURITY.md with best practices
- Add security-check.sh for pre-commit validation
- Include instructions for key management"
```

### 步骤 4: 提交其他文件（可选）

```bash
# 如果需要提交 web/README.md
git add web/README.md
git commit -m "docs: Update web README"

# 如果需要提交 package-lock.json
git add web/package-lock.json
git commit -m "chore: Update package-lock.json"
```

---

## 🔍 提交前最后检查

### 运行安全检查
```bash
./scripts/security-check.sh
```

**预期输出**：
```
🔍 开始安全检查...

📋 检查 1: .env 文件保护
✅ .env 已被 .gitignore 排除

📋 检查 2: node_modules 保护
✅ node_modules 已被 .gitignore 排除

📋 检查 3: 硬编码密钥检测
✅ 未发现硬编码密钥

📋 检查 4: Git 暂存区检查
✅ .env 正在从 git 中删除（正确操作）

📋 检查 5: 大文件检测
✅ 未发现过大文件

📋 检查 6: 编译产物检查
✅ 无编译产物在暂存区

================================
✅ 安全检查通过！可以安全提交。
```

### 手动检查清单

- [ ] `.env` 文件不在 `git status` 的"要提交的变更"中（除非是删除操作）
- [ ] `.gitignore` 文件包含 `.env` 规则
- [ ] 代码中没有硬编码的 API Key
- [ ] `web/node_modules/` 不在提交列表中
- [ ] 编译产物（`bin/`）不在提交列表中

---

## 🚀 一键提交（推荐）

如果你想一次性提交所有安全相关的修改：

```bash
# 1. 运行安全检查
./scripts/security-check.sh

# 2. 添加所有安全相关文件
git add .gitignore \
        examples/agent/weather_agent.go \
        examples/llm/basic_chat.go \
        SECURITY.md \
        scripts/security-check.sh

# 3. 查看将要提交的内容
git status

# 4. 提交
git commit -m "security: Comprehensive security improvements

- Add .gitignore to protect sensitive files
- Remove .env from git tracking
- Remove hardcoded API keys from examples
- Add security documentation and check script
- Ensure all secrets use environment variables

Fixes: Prevent accidental exposure of API keys"

# 5. 推送到 GitHub
git push origin main
```

---

## ⚠️ 重要提醒

### 如果 .env 已经被推送到 GitHub

如果你之前已经把 `.env` 推送到了 GitHub，**仅仅删除它是不够的**，因为它仍然在 git 历史中！

#### 解决方案 1: 撤销密钥（推荐）

1. **立即撤销所有泄露的 API Key**
   - 登录 API 提供商控制台
   - 删除或禁用泄露的密钥
   - 生成新的密钥

2. **更新本地 `.env`**
   ```bash
   # 编辑 .env，使用新密钥
   vim .env
   ```

3. **继续正常提交**
   ```bash
   git push origin main
   ```

#### 解决方案 2: 清理 Git 历史（高级）

⚠️ **警告**：这会重写 git 历史，影响所有协作者！

```bash
# 使用 BFG Repo-Cleaner
brew install bfg

# 删除 .env 文件的所有历史记录
bfg --delete-files .env

# 清理
git reflog expire --expire=now --all
git gc --prune=now --aggressive

# 强制推送（危险！）
git push --force origin main
```

**更好的做法**：直接撤销密钥，不要清理历史（除非必要）。

---

## 📊 验证提交结果

### 推送后验证

1. **访问 GitHub 仓库**
   ```
   https://github.com/your-username/einoflow
   ```

2. **检查文件列表**
   - ✅ `.gitignore` 应该存在
   - ❌ `.env` 不应该存在
   - ✅ `.env.example` 应该存在
   - ❌ `web/node_modules/` 不应该存在

3. **检查代码**
   - 打开 `examples/llm/basic_chat.go`
   - 确认使用 `os.Getenv("OPENAI_API_KEY")` 而不是硬编码

4. **搜索敏感信息**
   - 在 GitHub 仓库中搜索你的 API Key 的一部分
   - 应该找不到任何结果

---

## 🎉 完成后的状态

### 本地
```bash
# .env 文件仍然存在（但被 git 忽略）
$ ls -la .env
-rw-r--r--  1 user  staff  500 Nov 17 08:25 .env

# git 不再跟踪它
$ git status .env
位于分支 main
无文件要提交，干净的工作区
```

### GitHub
```
✅ .gitignore - 存在
❌ .env - 不存在（被忽略）
✅ .env.example - 存在（模板）
✅ SECURITY.md - 存在（文档）
✅ examples/*.go - 使用环境变量
```

---

## 🔄 后续维护

### 每次提交前运行检查
```bash
# 添加到 git hooks（可选）
echo '#!/bin/bash' > .git/hooks/pre-commit
echo './scripts/security-check.sh' >> .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### 定期审查
- 每月检查一次 `.gitignore` 规则
- 每季度轮换一次 API Key
- 发现问题立即修复

---

## 📞 需要帮助？

如果遇到问题：

1. **运行安全检查**
   ```bash
   ./scripts/security-check.sh
   ```

2. **查看文档**
   - `SECURITY.md` - 安全最佳实践
   - `README.md` - 项目说明

3. **常见问题**
   - Q: `.env` 还在 git 中？
   - A: 运行 `git rm --cached .env`
   
   - Q: 发现硬编码密钥？
   - A: 使用 `os.Getenv()` 替换
   
   - Q: 密钥已泄露？
   - A: 立即撤销并生成新密钥

---

## ✅ 准备好了吗？

现在你可以安全地提交代码了！

```bash
# 最后一次检查
./scripts/security-check.sh

# 如果通过，执行提交
git push origin main
```

🎉 **恭喜！你的代码现在是安全的！** 🔒
