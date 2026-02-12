import { Link } from 'react-router-dom';
import { MessageSquare, Bot, Database, GitBranch, ArrowRight, Sparkles } from 'lucide-react';
import { Card } from '../components/ui/card';
import { Button } from '../components/ui/button';

const HomePage = () => {
  const features = [
    {
      icon: MessageSquare,
      title: 'AI 对话',
      description: '与多个 AI 模型进行智能对话，支持流式输出和模型切换',
      path: '/chat',
      color: 'from-blue-500 to-cyan-500',
      features: ['多模型支持', '流式响应', '对话历史'],
    },
    {
      icon: Bot,
      title: 'AI Agent',
      description: '智能 Agent 帮助你完成复杂任务，包括写作、分析、代码生成',
      path: '/agent',
      color: 'from-purple-500 to-pink-500',
      features: ['任务规划', '自动执行', '结果总结'],
    },
    {
      icon: Database,
      title: 'RAG 检索',
      description: '索引你的文档，通过智能检索获取精准答案',
      path: '/rag',
      color: 'from-green-500 to-emerald-500',
      features: ['文档索引', '语义检索', '来源追溯'],
    },
    {
      icon: GitBranch,
      title: 'Graph 编排',
      description: '多步骤图编排处理复杂问题，自动分析和执行',
      path: '/graph',
      color: 'from-orange-500 to-red-500',
      features: ['多步骤', '自动编排', '可视化流程'],
    },
  ];

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="bg-gradient-to-br from-slate-900 via-slate-800 to-blue-900 px-6 py-20 text-white">
        <div className="mx-auto max-w-6xl">
          <div className="flex items-center gap-2 text-sm">
            <Sparkles className="h-4 w-4" />
            <span className="uppercase tracking-wider text-slate-300">
              Powered by Eino Framework
            </span>
          </div>
          <h1 className="mt-6 text-5xl font-bold leading-tight lg:text-6xl">
            智能 AI 编排平台
            <br />
            <span className="bg-gradient-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">
              EinoFlow
            </span>
          </h1>
          <p className="mt-6 max-w-2xl text-xl text-slate-300">
            基于 Eino 框架构建的 AI 应用平台，提供对话、Agent、RAG、Graph 等完整功能
          </p>
          <div className="mt-8 flex flex-wrap gap-4">
            <Link to="/chat">
              <Button size="lg" className="shadow-xl">
                开始使用
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </Link>
            <a
              href="https://github.com/cloudwego/eino"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Button variant="outline" size="lg" className="border-white/30 text-white hover:bg-white/10">
                了解 Eino
              </Button>
            </a>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="px-6 py-20">
        <div className="mx-auto max-w-6xl">
          <div className="mb-12 text-center">
            <h2 className="text-3xl font-bold text-slate-900">核心功能</h2>
            <p className="mt-3 text-slate-600">
              探索 EinoFlow 提供的强大 AI 能力
            </p>
          </div>

          <div className="grid gap-8 md:grid-cols-2">
            {features.map((feature) => {
              const Icon = feature.icon;
              return (
                <Link key={feature.path} to={feature.path}>
                  <Card className="group h-full p-6 transition-all hover:shadow-xl">
                    <div className="flex items-start gap-4">
                      <div
                        className={`flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-gradient-to-br ${feature.color} shadow-lg`}
                      >
                        <Icon className="h-6 w-6 text-white" />
                      </div>
                      <div className="flex-1">
                        <h3 className="text-xl font-semibold text-slate-900">
                          {feature.title}
                        </h3>
                        <p className="mt-2 text-sm text-slate-600">
                          {feature.description}
                        </p>
                        <div className="mt-4 flex flex-wrap gap-2">
                          {feature.features.map((f) => (
                            <span
                              key={f}
                              className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700"
                            >
                              {f}
                            </span>
                          ))}
                        </div>
                        <div className="mt-4 flex items-center gap-2 text-sm font-medium text-blue-600 transition-all group-hover:gap-3">
                          立即体验
                          <ArrowRight className="h-4 w-4" />
                        </div>
                      </div>
                    </div>
                  </Card>
                </Link>
              );
            })}
          </div>
        </div>
      </section>

      {/* Tech Stack Section */}
      <section className="bg-slate-50 px-6 py-20">
        <div className="mx-auto max-w-6xl">
          <div className="mb-12 text-center">
            <h2 className="text-3xl font-bold text-slate-900">技术栈</h2>
            <p className="mt-3 text-slate-600">
              使用现代化技术构建
            </p>
          </div>

          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            <Card className="p-6 text-center">
              <h3 className="font-semibold text-slate-900">后端</h3>
              <p className="mt-2 text-sm text-slate-600">
                Go + Eino Framework
              </p>
            </Card>
            <Card className="p-6 text-center">
              <h3 className="font-semibold text-slate-900">前端</h3>
              <p className="mt-2 text-sm text-slate-600">
                React 18 + TypeScript
              </p>
            </Card>
            <Card className="p-6 text-center">
              <h3 className="font-semibold text-slate-900">UI</h3>
              <p className="mt-2 text-sm text-slate-600">
                TailwindCSS + shadcn/ui
              </p>
            </Card>
            <Card className="p-6 text-center">
              <h3 className="font-semibold text-slate-900">AI 模型</h3>
              <p className="mt-2 text-sm text-slate-600">
                豆包 + OpenAI
              </p>
            </Card>
          </div>
        </div>
      </section>
    </div>
  );
};

export default HomePage;
