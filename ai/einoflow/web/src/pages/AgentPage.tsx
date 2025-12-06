import { useState } from 'react';
import { Bot, Loader2, Sparkles, Play, CheckCircle2, AlertCircle } from 'lucide-react';
import { Button } from '../components/ui/button';
import { Card } from '../components/ui/card';
import { MarkdownRenderer } from '../components/MarkdownRenderer';
import { runAgent, type AgentResponse } from '../api/agent';

const AgentPage = () => {
  const [task, setTask] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<AgentResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleRun = async () => {
    if (!task.trim() || loading) return;

    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await runAgent({ task: task.trim() });
      setResult(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : '发生未知错误');
    } finally {
      setLoading(false);
    }
  };

  const examples = [
    '给我写一篇关于 Go 语言的文章',
    '解释什么是 goroutine 和 channel',
    '用 Go 实现一个简单的 HTTP 服务器',
    '分析微服务架构的优缺点',
  ];

  return (
    <div className="h-screen overflow-y-auto bg-slate-50 px-4 py-4">
      <div className="mx-auto max-w-6xl">
        {/* Header */}
        <div className="mb-4">
          <h1 className="text-2xl font-semibold text-slate-900">AI Agent</h1>
          <p className="mt-1 text-sm text-slate-600">
            智能 Agent 可以帮助你完成复杂任务，包括写作、分析、代码生成等
          </p>
        </div>

        {/* Input Card */}
        <Card className="p-4">
          <div className="space-y-3">
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                任务描述
              </label>
              <textarea
                value={task}
                onChange={(e) => setTask(e.target.value)}
                placeholder="描述你想要 Agent 完成的任务..."
                className="w-full resize-none rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                rows={4}
                disabled={loading}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="text-sm text-slate-500">
                💡 提示：当前是简化版 Agent，适合知识问答、写作、分析等任务
              </div>
              <Button onClick={handleRun} disabled={loading || !task.trim()} size="lg">
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                    执行中...
                  </>
                ) : (
                  <>
                    <Play className="mr-2 h-5 w-5" />
                    运行 Agent
                  </>
                )}
              </Button>
            </div>
          </div>
        </Card>

        {/* Examples */}
        <div className="mt-4">
          <p className="mb-2 text-sm font-medium text-slate-700">示例任务：</p>
          <div className="grid gap-2 sm:grid-cols-2">
            {examples.map((example, index) => (
              <button
                key={index}
                onClick={() => setTask(example)}
                disabled={loading}
                className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-left text-sm text-slate-700 transition-colors hover:border-blue-300 hover:bg-blue-50 disabled:opacity-50"
              >
                {example}
              </button>
            ))}
          </div>
        </div>

        {/* Result */}
        {result && (
          <Card className="mt-4 p-4">
            <div className="mb-3 flex items-center gap-2">
              <CheckCircle2 className="h-5 w-5 text-green-600" />
              <h3 className="text-base font-semibold text-slate-900">执行结果</h3>
              {result.execution_time && (
                <span className="ml-auto text-sm text-slate-500">
                  耗时: {result.execution_time.toFixed(2)}s
                </span>
              )}
            </div>
            <div className="prose prose-slate max-w-none">
              <div className="rounded-lg bg-slate-50 p-4">
                <MarkdownRenderer content={result.answer} />
              </div>
            </div>
            {result.steps && result.steps.length > 0 && (
              <div className="mt-6">
                <h4 className="mb-3 text-sm font-semibold text-slate-700">执行步骤：</h4>
                <div className="space-y-2">
                  {result.steps.map((step, index) => (
                    <div
                      key={index}
                      className="flex items-start gap-3 rounded-lg bg-slate-50 p-3"
                    >
                      <div className="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-600">
                        {index + 1}
                      </div>
                      <p className="text-sm text-slate-700">{step}</p>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </Card>
        )}

        {/* Error */}
        {error && (
          <Card className="mt-8 border-red-200 bg-red-50 p-6">
            <div className="flex items-start gap-3">
              <AlertCircle className="h-5 w-5 flex-shrink-0 text-red-600" />
              <div>
                <h3 className="font-semibold text-red-900">执行失败</h3>
                <p className="mt-1 text-sm text-red-700">{error}</p>
              </div>
            </div>
          </Card>
        )}
      </div>
    </div>
  );
};

export default AgentPage;
