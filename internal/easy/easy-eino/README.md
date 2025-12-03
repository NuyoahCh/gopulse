# eino-tutorial

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go)](https://go.dev/doc/devel/release)
[![CloudWeGo Eino](https://img.shields.io/badge/CloudWeGo-Eino-0063D1)](https://www.cloudwego.io/zh/docs/eino/)

> 一组基于 [CloudWeGo Eino](https://github.com/cloudwego/eino) 的 Go 语言示例，演示如何构建链路、模板、工具、Embedding、回调等多种 AI 应用场景。

## 目录

- [项目概览](#项目概览)
- [功能亮点](#功能亮点)
- [目录结构](#目录结构)
- [准备工作](#准备工作)
- [快速上手](#快速上手)
- [运行示例](#运行示例)
- [环境变量](#环境变量)
- [常用命令](#常用命令)
- [贡献指南](#贡献指南)
- [常见问题](#常见问题)
- [参考资料](#参考资料)
- [许可证](#许可证)

## 项目概览

`eino-tutorial` 旨在帮助开发者快速理解并实践 CloudWeGo Eino 生态。仓库中的 Go 示例覆盖了从最基础的对话式调用，到链路编排、模板管理、向量检索、工具调用等进阶特性，并配套中文说明文档，适合作为入门或团队培训材料。

## 功能亮点

- **链路编排**：通过 `compose.Chain` 串联 ChatTemplate 与 ChatModel，示范多步推理、分支、并行与流式处理。参见 `example-chain/`。
- **模型调用**：展示如何使用 DeepSeek、ARK（豆包）、Ollama 等模型提供商，包含单轮、多轮、参数化与流式输出。参见 `example-chatmodel/`、`example-doubao/`、`example-ollama/`。
- **提示模板**：演示使用 `prompt.FromMessages` 构建复杂模板，包括角色设定、变量插值和多模型协同。参见 `example-chattemplate/`。
- **工具与代理**：结合函数调用、数据库检索、本地文件等工具，构建可扩展的智能体。参见 `example-tool/` 与 `example-react/`。
- **回调与监控**：通过回调函数实时观测链路执行、收集指标。参见 `example-callback/`。
- **Embedding 与向量检索**：完整流程包含文档加载、向量化、内存存储及检索调用。参见 `example-embedding/`。

## 目录结构

```
.
├── example-chain/         # 链式调用与组合式工作流示例
├── example-chatmodel/     # 不同对话模型的调用姿势
├── example-chattemplate/  # ChatTemplate 构建与复用
├── example-callback/      # 回调与可观测性实践
├── example-deepseek/      # DeepSeek 模型最小可运行示例
├── example-doubao/        # 字节豆包（ARK）模型调用示例
├── example-embedding/     # Embedding、向量检索与文档代理
├── example-ollama/        # 本地 Ollama 模型的集成
├── example-react/         # ReAct 风格的智能体与工具链
└── example-tool/          # 自定义工具与 ToolSet 管理
```

每个目录同时包含 Go 源码与运行输出的 Markdown 记录，便于比对代码与实际结果。

## 准备工作

1. **安装 Go**：确保本地 Go 版本 ≥ 1.24。可以通过 `go version` 验证。
2. **克隆仓库**：
   ```bash
   git clone https://github.com/your-org/eino-tutorial.git
   cd eino-tutorial
   ```
3. **依赖管理**：首次进入仓库时运行 `go mod tidy`（可选），以拉取缺失依赖。
4. **外部服务凭据**：部分示例需要可用的云端 API Key（例如 DeepSeek、ARK/豆包），具体见下文 [环境变量](#环境变量)。
5. **本地模型服务**：若要运行 Ollama 示例，请提前安装并启动 [Ollama](https://ollama.com/) 服务，并确保模型（如 `llama2`）已下载。

## 快速上手

以下步骤演示如何运行最简单的链式对话示例：

```bash
# 1. 设置必要的环境变量（以 DeepSeek 为例）
export DEEPSEEK_API_KEY="sk-xxxxx"

# 2. 运行示例
go run ./example-chain/easy.go
```

程序将构建一个包含系统/用户信息的 ChatTemplate，调用 DeepSeek ChatModel，并打印模型回复。

## 运行示例

- **单个示例**：使用 `go run` 指向具体文件，例如：
  ```bash
  go run ./example-react/multiple.go
  go run ./example-embedding/vectoring.go
  ```
- **批量体验**：某些目录下的 Markdown 文件（如 `example-chain/*.md`）记录了预期输出，可作为运行后的参考。
- **自定义参数**：修改源码中的模板内容、工具定义或链路结构，即可快速验证新的想法。

> 提示：由于示例之间共享模块依赖，可在仓库根目录统一运行命令，无需额外配置 GOPATH。

## 环境变量

下表列出了常见示例所需的环境变量：

| 变量名 | 用途 | 适用示例 |
| ------ | ---- | -------- |
| `DEEPSEEK_API_KEY` | 调用 DeepSeek ChatModel | `example-chain`, `example-chatmodel`, `example-react`, `example-callback`, `example-ollama` 等 |
| `ARK_API_KEY` | 访问字节豆包（ARK）平台 | `example-doubao`, `example-embedding` |
| `ARK_MODEL_NAME` | 指定 Doubao ChatModel 名称（如 `ep-xxxx`） | `example-doubao` |
| `ARK_EMBEDDING_MODEL` | 指定 Doubao Embedding 模型（如 `doubao-embedding-large`） | `example-embedding` |

运行前请确保 API Key 已配置到环境变量或通过 `.env`/配置管理工具注入。

## 常用命令

```bash
# 检查依赖是否可用
go list ./...

# 运行全部（如有测试）
go test ./...

# 静态检查（可选）
GOEXPERIMENT=boringcrypto go vet ./...
```

虽然仓库目前以示例为主，建议在自定义扩展时结合 `go fmt`, `go vet`, `staticcheck` 等工具保持代码质量。

## 贡献指南

欢迎通过 Issue 或 PR 分享改进意见、补充示例或修正文档：

1. Fork 仓库并基于最新 `main` 分支创建特性分支。
2. 完成修改后运行必要的示例或检查，确保代码可执行。
3. 提交具有描述性的提交信息，并在 PR 中说明动机、主要变更与验证方式。
4. 若新增示例，请同步更新对应的 Markdown 输出，便于读者对照。

## 常见问题

- **Q: 没有云端 API Key 可以运行哪些示例？**
  - A: 你仍然可以运行 `example-tool/` 中的本地工具示例、`example-react/restaurant.go` 等无需外部模型的场景，以及基于 Ollama 的本地模型（需自备模型）。
- **Q: 为什么运行示例时提示 401/403？**
  - A: 请检查 API Key 是否填写正确，是否具备相应模型的调用权限，或是否超过调用额度。
- **Q: 是否提供英文文档？**
  - A: 目前仓库内 Markdown 说明以中文为主，欢迎贡献翻译版本。

## 参考资料

- [CloudWeGo Eino 官方文档](https://www.cloudwego.io/zh/docs/eino/)
- [CloudWeGo Eino GitHub 仓库](https://github.com/cloudwego/eino)
- [CloudWeGo Eino 扩展组件](https://github.com/cloudwego/eino-ext)
- [Ollama 官方网站](https://ollama.com/)

## 许可证

当前仓库尚未提供显式的开源许可证。若需在生产环境或商业场景中使用，请联系仓库维护者确认授权条款。
