import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';
import { MessageSquare, Bot, Database, GitBranch, Home } from 'lucide-react';
import HomePage from './pages/HomePage';
import ChatPage from './pages/ChatPage';
import AgentPage from './pages/AgentPage';
import RAGPage from './pages/RAGPage';
import GraphPage from './pages/GraphPage';

const Navigation = () => {
  const location = useLocation();

  const links = [
    { path: '/', icon: Home, label: '首页' },
    { path: '/chat', icon: MessageSquare, label: '对话' },
    { path: '/agent', icon: Bot, label: 'Agent' },
    { path: '/rag', icon: Database, label: 'RAG' },
    { path: '/graph', icon: GitBranch, label: 'Graph' },
  ];

  return (
    <nav className="border-b bg-white">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        <Link to="/" className="flex items-center gap-2">
          <div className="rounded-lg bg-gradient-to-br from-blue-600 to-indigo-600 p-2">
            <GitBranch className="h-6 w-6 text-white" />
          </div>
          <span className="text-xl font-semibold text-slate-900">EinoFlow</span>
        </Link>

        <div className="flex gap-2">
          {links.map((link) => {
            const Icon = link.icon;
            const isActive = location.pathname === link.path;
            return (
              <Link
                key={link.path}
                to={link.path}
                className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-blue-600 text-white'
                    : 'text-slate-700 hover:bg-slate-100'
                }`}
              >
                <Icon className="h-4 w-4" />
                {link.label}
              </Link>
            );
          })}
        </div>
      </div>
    </nav>
  );
};

const App = () => {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-slate-50">
        <Navigation />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/chat" element={<ChatPage />} />
          <Route path="/agent" element={<AgentPage />} />
          <Route path="/rag" element={<RAGPage />} />
          <Route path="/graph" element={<GraphPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
};

export default App;
