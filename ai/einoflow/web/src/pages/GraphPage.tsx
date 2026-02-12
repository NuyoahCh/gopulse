import { useState } from 'react';
import { Play, Loader2, GitBranch, CheckCircle2 } from 'lucide-react';
import { Button } from '../components/ui/button';
import { Card } from '../components/ui/card';
import { runGraph, type GraphResponse } from '../api/graph';

const GraphPage = () => {
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<GraphResponse | null>(null);

  const handleRun = async () => {
    if (!input.trim() || loading) return;

    setLoading(true);
    setResult(null);

    try {
      const response = await runGraph({ query: input.trim() });
      setResult(response);
    } catch (error) {
      console.error('Graph execution failed:', error);
    } finally {
      setLoading(false);
    }
  };

  const examples = [
    '如何设计一个高可用的微服务架构？',
    '分析电商系统的核心业务流程',
    '制定一个产品从0到1的完整计划',
    '设计一个分布式缓存系统',
  ];

  return (
    <div className="h-screen overflow-y-auto bg-slate-50 px-4 py-4">
      <div className="mx-auto max-w-6xl">
        {/* Header */}
        <div className="mb-4">
          <h1 className="text-2xl font-semibold text-slate-900">Graph 多步骤处理</h1>
          <p className="mt-1 text-sm text-slate-600">
            通过图编排实现复杂任务的多步骤分析和处理
          </p>
        </div>

        {/* Input Card */}
        <Card className="p-4">
          <div className="space-y-3">
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                复杂问题
              </label>
              <textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="输入需要多步骤分析的复杂问题..."
                className="w-full resize-none rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                rows={4}
                disabled={loading}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-sm text-slate-500">
                <GitBranch className="h-4 w-4" />
                Graph 会自动进行：问题分析 → 计划制定 → 执行总结
              </div>
              <Button onClick={handleRun} disabled={loading || !input.trim()} size="lg">
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                    执行中...
                  </>
                ) : (
                  <>
                    <Play className="mr-2 h-5 w-5" />
                    运行 Graph
                  </>
                )}
              </Button>
            </div>
          </div>
        </Card>

        {/* Examples */}
        <div className="mt-4">
          <p className="mb-2 text-sm font-medium text-slate-700">示例问题：</p>
          <div className="grid gap-2 sm:grid-cols-2">
            {examples.map((example, index) => (
              <button
                key={index}
                onClick={() => setInput(example)}
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
          <div className="mt-4 space-y-3">
            {/* Steps */}
            {result.steps && result.steps.length > 0 && (
              <Card className="p-4">
                <h3 className="mb-3 flex items-center gap-2 text-base font-semibold text-slate-900">
                  <GitBranch className="h-5 w-5 text-blue-600" />
                  执行步骤
                </h3>
                <div className="space-y-4">
                  {result.steps.map((step, index) => (
                    <div key={index} className="relative">
                      {index < result.steps!.length - 1 && (
                        <div className="absolute left-4 top-10 h-full w-0.5 bg-slate-200" />
                      )}
                      <div className="flex gap-4">
                        <div className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-blue-100">
                          <CheckCircle2 className="h-4 w-4 text-blue-600" />
                        </div>
                        <div className="flex-1">
                          <h4 className="mb-2 font-semibold text-slate-900">
                            {step.node === 'analyze' && '步骤 1: 问题分析'}
                            {step.node === 'plan' && '步骤 2: 计划制定'}
                            {step.node === 'execute' && '步骤 3: 执行总结'}
                            {!['analyze', 'plan', 'execute'].includes(step.node) && `步骤: ${step.node}`}
                          </h4>
                          <div className="rounded-lg bg-slate-50 p-4">
                            <p className="whitespace-pre-wrap text-sm text-slate-700">
                              {step.output}
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </Card>
            )}

            {/* Final Result */}
            <Card className="p-4">
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-base font-semibold text-slate-900">最终结果</h3>
                {result.execution_time && (
                  <span className="text-sm text-slate-500">
                    耗时: {result.execution_time.toFixed(2)}s
                  </span>
                )}
              </div>
              <div className="rounded-lg bg-gradient-to-br from-blue-50 to-indigo-50 p-6">
                <p className="whitespace-pre-wrap text-sm leading-relaxed text-slate-900">
                  {result.result}
                </p>
              </div>
            </Card>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <Card className="mt-8 p-10 text-center">
            <Loader2 className="mx-auto h-12 w-12 animate-spin text-blue-600" />
            <p className="mt-4 text-slate-600">Graph 正在执行多步骤处理...</p>
            <p className="mt-2 text-sm text-slate-500">
              这可能需要 10-30 秒，请耐心等待
            </p>
          </Card>
        )}
      </div>
    </div>
  );
};

export default GraphPage;
