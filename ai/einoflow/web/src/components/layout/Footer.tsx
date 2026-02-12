const Footer = () => {
  return (
    <footer className="border-t border-slate-200 bg-white/80">
      <div className="mx-auto flex max-w-6xl flex-col gap-2 px-6 py-6 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between">
        <p>© {new Date().getFullYear()} Einoflow. All rights reserved.</p>
        <div className="flex gap-4">
          <a className="hover:text-slate-900" href="https://github.com" target="_blank" rel="noreferrer">
            GitHub
          </a>
          <a className="hover:text-slate-900" href="https://shadcn.com" target="_blank" rel="noreferrer">
            UI System
          </a>
          <a className="hover:text-slate-900" href="mailto:hello@einoflow.io">
            Contact
          </a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
