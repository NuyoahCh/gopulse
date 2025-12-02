go run callback.go


\n[Callback] 开始执行: ReActAgent\n输入: [
{
"role": "user",
"content": "测试一下"
}
]\n\n[Callback] 开始执行: ChatModel\n输入: {
"Messages": [
{
"role": "user",
"content": "测试一下"
}
],
"Tools": [
{
"Name": "test_tool",
"Desc": "测试工具",
"Extra": null
}
],
"ToolChoice": "allowed",
"Config": {
"Model": "deepseek-chat",
"MaxTokens": 0,
"Temperature": 0,
"TopP": 0,
"Stop": null
},
"Extra": null
}\n\n[Callback] 完成执行: ChatModel\n\n[Callback] 完成执行: ReActAgent\n\n最终回答: 我目前无法执行测试操作。虽然有一个测试工具可用，但它的具体功能暂时无法调用。

如果您需要进行其他类型的测试或帮助，我很乐意为您提供支持。您可以告诉我具体想要测试什么内容，我会尽力通过其他方式来协助您。\n%                          