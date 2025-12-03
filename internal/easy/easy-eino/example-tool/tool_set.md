go run tool_set.go


\n========== 测试 1 ==========\n问题: 现在几点了？\nAI 的决策:\n  ✓ 使用工具: get_time\n  参数: {}\n  结果: [tool: 2025-11-14 09:41:05
tool_call_id: call_00_gWA1ouMcKkqvxB8JDqRtB6Yu
tool_call_name: get_time]\n\n========== 测试 2 ==========\n问题: 帮我计算 10 + 20 等于多少\nAI 的决策:\n  ✓ 使用工具: calculator\n  参数: {"expression": "10 + 20"}\n  结果: [tool: 30
tool_call_id: call_00_ma99WVmMhIwpp33COPuivIXk
tool_call_name: calculator]\n\n========== 测试 3 ==========\n问题: 你好，请介绍一下你自己\nAI 的决策:\n  ✓ 直接回答: 你好！我是一个AI助手，很高兴为你服务！我可以帮助你：

- 回答各种问题和提供信息
- 进行数学计算
- 获取当前时间
- 协助解决学习和工作中的问题
- 提供建议和指导

我目前具备计算器和时间查询的功能，可以帮你进行数学运算或者了解当前时间。

有什么我可以帮助你的吗？无论是学习、工作还是日常生活中的问题，我都很乐意协助你！\n%                                
