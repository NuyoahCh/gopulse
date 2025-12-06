import { useState, useEffect, useRef } from 'react';
import { Send, Bot, User, Loader2 } from 'lucide-react';
import { Button } from '../components/ui/button';
import { Card } from '../components/ui/card';
import { MarkdownRenderer } from '../components/MarkdownRenderer';
import { chat, chatStream, listModels, type ChatMessage, type Model } from '../api/llm';

const ChatPage = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [streaming, setStreaming] = useState(false);
  const [models, setModels] = useState<Model[]>([]);
  const [selectedModel, setSelectedModel] = useState<string>('');
  const [selectedProvider, setSelectedProvider] = useState<string>('');
  const [loadingModels, setLoadingModels] = useState(true);
  const [error, setError] = useState<string>('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // 加载可用模型
    setLoadingModels(true);
    listModels()
      .then((response) => {
        setModels(response.models);
        if (response.models.length > 0) {
          setSelectedProvider(response.models[0].provider);
          setSelectedModel(response.models[0].id);
        }
        setError('');
      })
      .catch((err) => {
        console.error('Failed to load models:', err);
        setError('无法加载模型列表，请确保后端服务正在运行');
      })
      .finally(() => {
        setLoadingModels(false);
      });
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim() || loading) return;

    const userMessage: ChatMessage = {
      role: 'user',
      content: input.trim(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      if (streaming) {
        // 流式响应 - 立即关闭 loading 状态，避免显示两个消息框
        setLoading(false);
        
        let assistantContent = '';
        const assistantMessage: ChatMessage = {
          role: 'assistant',
          content: '',
        };
        setMessages((prev) => [...prev, assistantMessage]);

        await chatStream(
          {
            provider: selectedProvider,
            model: selectedModel,
            messages: [...messages, userMessage],
            stream: true,
          },
          (content) => {
            assistantContent += content;
            setMessages((prev) => {
              const newMessages = [...prev];
              newMessages[newMessages.length - 1] = {
                role: 'assistant',
                content: assistantContent,
              };
              return newMessages;
            });
          },
          () => {
            // 流式完成，无需额外操作
          },
          (error) => {
            console.error('Stream error:', error);
            setMessages((prev) => [
              ...prev.slice(0, -1),
              {
                role: 'assistant',
                content: '抱歉，发生了错误。请重试。',
              },
            ]);
          }
        );
      } else {
        // 普通响应
        const response = await chat({
          provider: selectedProvider,
          model: selectedModel,
          messages: [...messages, userMessage],
        });

        setMessages((prev) => [
          ...prev,
          {
            role: 'assistant',
            content: response.content,
          },
        ]);
        setLoading(false);
      }
    } catch (error) {
      console.error('Chat error:', error);
      setLoading(false);
      setMessages((prev) => [
        ...prev,
        {
          role: 'assistant',
          content: '抱歉，发生了错误。请重试。',
        },
      ]);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  if (loadingModels) {
    return (
      <div className="flex h-screen items-center justify-center bg-slate-50">
        <div className="text-center">
          <Loader2 className="mx-auto h-12 w-12 animate-spin text-blue-600" />
          <p className="mt-4 text-slate-600">加载模型列表...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-screen items-center justify-center bg-slate-50">
        <Card className="max-w-md p-8 text-center">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-100">
            <Bot className="h-8 w-8 text-red-600" />
          </div>
          <h2 className="text-xl font-semibold text-slate-900">连接失败</h2>
          <p className="mt-2 text-sm text-slate-600">{error}</p>
          <Button
            onClick={() => window.location.reload()}
            className="mt-6"
          >
            重新加载
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex h-screen flex-col bg-slate-50">
      {/* Header */}
      <div className="border-b bg-white px-4 py-3">
        <div className="mx-auto flex max-w-6xl items-center justify-between">
          <div>
            <h1 className="text-xl font-semibold text-slate-900">AI 对话</h1>
            <p className="text-xs text-slate-500">与 AI 模型进行智能对话</p>
          </div>
          <div className="flex items-center gap-4">
            <select
              value={selectedProvider}
              onChange={(e) => {
                setSelectedProvider(e.target.value);
                const providerModels = models.filter((m) => m.provider === e.target.value);
                if (providerModels.length > 0) {
                  setSelectedModel(providerModels[0].id);
                }
              }}
              className="rounded-lg border border-slate-200 px-3 py-2 text-sm"
            >
              {Array.from(new Set(models.map((m) => m.provider))).map((provider) => (
                <option key={provider} value={provider}>
                  {provider}
                </option>
              ))}
            </select>
            <select
              value={selectedModel}
              onChange={(e) => setSelectedModel(e.target.value)}
              className="rounded-lg border border-slate-200 px-3 py-2 text-sm"
            >
              {models
                .filter((m) => m.provider === selectedProvider)
                .map((model) => (
                  <option key={model.id} value={model.id}>
                    {model.name}
                  </option>
                ))}
            </select>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={streaming}
                onChange={(e) => setStreaming(e.target.checked)}
                className="rounded"
              />
              流式输出
            </label>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-4 py-4">
        <div className="mx-auto max-w-6xl space-y-4">
          {messages.length === 0 && (
            <Card className="p-10 text-center">
              <Bot className="mx-auto h-12 w-12 text-slate-400" />
              <h3 className="mt-4 text-lg font-semibold text-slate-900">开始对话</h3>
              <p className="mt-2 text-sm text-slate-500">
                输入你的问题，AI 会为你提供帮助
              </p>
            </Card>
          )}

          {messages.map((message, index) => (
            <div
              key={index}
              className={`flex gap-4 ${
                message.role === 'user' ? 'justify-end' : 'justify-start'
              }`}
            >
              {message.role === 'assistant' && (
                <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-blue-100">
                  <Bot className="h-5 w-5 text-blue-600" />
                </div>
              )}
              <Card
                className={`max-w-[80%] p-4 ${
                  message.role === 'user'
                    ? 'bg-blue-600 text-white'
                    : 'bg-white'
                }`}
              >
                {message.role === 'user' ? (
                  <p className="whitespace-pre-wrap text-sm">{message.content}</p>
                ) : (
                  <MarkdownRenderer content={message.content} className="text-sm" />
                )}
              </Card>
              {message.role === 'user' && (
                <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-slate-200">
                  <User className="h-5 w-5 text-slate-600" />
                </div>
              )}
            </div>
          ))}

          {loading && (
            <div className="flex gap-4">
              <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-blue-100">
                <Bot className="h-5 w-5 text-blue-600" />
              </div>
              <Card className="p-4">
                <Loader2 className="h-5 w-5 animate-spin text-slate-400" />
              </Card>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input */}
      <div className="border-t bg-white px-4 py-3">
        <div className="mx-auto flex max-w-6xl gap-3">
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="输入你的问题..."
            className="flex-1 resize-none rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            rows={2}
            disabled={loading}
          />
          <Button
            onClick={handleSend}
            disabled={loading || !input.trim()}
            size="lg"
            className="self-end"
          >
            {loading ? (
              <Loader2 className="h-5 w-5 animate-spin" />
            ) : (
              <Send className="h-5 w-5" />
            )}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ChatPage;
