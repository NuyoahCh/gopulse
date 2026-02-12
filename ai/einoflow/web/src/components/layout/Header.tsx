import { Brain, Github, Menu, PanelsTopLeft } from 'lucide-react';
import { Button } from '../ui/button';

const Header = () => {
  return (
    <header className="sticky top-0 z-20 border-b border-slate-200 bg-white/90 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <div className="flex items-center gap-3">
          <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-900 text-white">
            <Brain size={22} />
          </div>
          <div>
            <p className="text-sm font-semibold uppercase tracking-widest text-slate-500">Einoflow</p>
            <h1 className="text-xl font-semibold text-slate-900">Workflow Console</h1>
          </div>
        </div>
        <div className="hidden items-center gap-3 md:flex">
          <Button variant="ghost" className="gap-2 text-slate-600">
            <PanelsTopLeft size={18} />
            Workbench
          </Button>
          <Button variant="ghost" className="gap-2 text-slate-600">
            <Github size={18} />
            GitHub
          </Button>
          <Button className="shadow-lg">Launch</Button>
        </div>
        <Button variant="ghost" className="md:hidden">
          <Menu size={20} />
        </Button>
      </div>
    </header>
  );
};

export default Header;
