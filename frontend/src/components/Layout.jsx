import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { LayoutDashboard, Package, Key, LogOut, BookOpen } from 'lucide-react';
import { clearToken } from '../lib/api';

const nav = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard', end: true },
  { to: '/apps', icon: Package, label: 'Apps' },
  { to: '/licenses', icon: Key, label: 'Licenses' },
  { to: '/docs', icon: BookOpen, label: 'API Docs' },
];

export default function Layout() {
  const navigate = useNavigate();

  const handleLogout = () => {
    clearToken();
    navigate('/login');
  };

  return (
    <div className="flex h-screen bg-[#0a0a0a]">
      <aside className="w-56 border-r border-zinc-800 bg-[#0f0f0f] flex flex-col">
        <div className="h-14 flex items-center px-5 border-b border-zinc-800">
          <Key className="w-5 h-5 text-blue-500 mr-2" />
          <span className="font-bold text-white text-sm tracking-tight">LicenseManager</span>
        </div>
        <nav className="flex-1 p-3 space-y-1">
          {nav.map(({ to, icon: Icon, label, end }) => (
            <NavLink
              key={to}
              to={to}
              end={end}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                  isActive
                    ? 'bg-blue-600/10 text-blue-400'
                    : 'text-zinc-400 hover:bg-zinc-800 hover:text-white'
                }`
              }
            >
              <Icon className="w-4 h-4" />
              {label}
            </NavLink>
          ))}
        </nav>
        <div className="p-3 border-t border-zinc-800">
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2 w-full rounded-lg text-sm text-zinc-500 hover:bg-zinc-800 hover:text-white transition-colors"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </button>
        </div>
      </aside>
      <main className="flex-1 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  );
}
