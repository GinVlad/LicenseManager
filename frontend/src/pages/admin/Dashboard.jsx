import { useState, useEffect } from 'react';
import { Package, Key, CheckCircle2, XCircle } from 'lucide-react';
import { api } from '../../lib/api';

export default function Dashboard() {
  const [apps, setApps] = useState([]);
  const [licenses, setLicenses] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([api.apps.list(), api.licenses.list()])
      .then(([a, l]) => { setApps(a || []); setLicenses(l || []); })
      .finally(() => setLoading(false));
  }, []);

  const active = licenses.filter((l) => l.isActive && new Date(l.expiresAt) > new Date()).length;
  const expired = licenses.filter((l) => new Date(l.expiresAt) <= new Date()).length;

  const stats = [
    { label: 'Apps', value: apps.length, icon: Package, color: 'text-blue-400', bg: 'bg-blue-500/10 border-blue-500/20' },
    { label: 'Total Licenses', value: licenses.length, icon: Key, color: 'text-purple-400', bg: 'bg-purple-500/10 border-purple-500/20' },
    { label: 'Active', value: active, icon: CheckCircle2, color: 'text-emerald-400', bg: 'bg-emerald-500/10 border-emerald-500/20' },
    { label: 'Expired', value: expired, icon: XCircle, color: 'text-red-400', bg: 'bg-red-500/10 border-red-500/20' },
  ];

  if (loading) {
    return <div className="p-8 text-zinc-500 text-sm">Loading...</div>;
  }

  return (
    <div className="p-8">
      <h1 className="text-lg font-bold text-white mb-6">Dashboard</h1>

      <div className="grid grid-cols-4 gap-4 mb-8">
        {stats.map(({ label, value, icon: Icon, color, bg }) => (
          <div key={label} className={`border rounded-xl p-5 ${bg}`}>
            <div className={`text-xs font-bold uppercase tracking-widest mb-2 ${color}`}>{label}</div>
            <div className={`text-3xl font-mono font-bold ${color}`}>{value}</div>
          </div>
        ))}
      </div>

      <div className="bg-[#0f0f0f] border border-zinc-800 rounded-xl overflow-hidden">
        <div className="px-5 py-3 border-b border-zinc-800 text-xs font-bold text-zinc-400 uppercase tracking-widest">
          Recent Licenses
        </div>
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-zinc-800 text-xs text-zinc-500">
              <th className="px-5 py-3 text-left">Key</th>
              <th className="px-5 py-3 text-left">Plan</th>
              <th className="px-5 py-3 text-left">Threads</th>
              <th className="px-5 py-3 text-left">Expires</th>
              <th className="px-5 py-3 text-left">Status</th>
            </tr>
          </thead>
          <tbody>
            {licenses.slice(0, 10).map((l) => {
              const expired = new Date(l.expiresAt) <= new Date();
              return (
                <tr key={l.id} className="border-b border-zinc-800/50 hover:bg-zinc-900/30">
                  <td className="px-5 py-3 font-mono text-xs text-white">{l.key}</td>
                  <td className="px-5 py-3 text-zinc-400">{l.plan}</td>
                  <td className="px-5 py-3 text-zinc-400">{l.maxThreads}</td>
                  <td className="px-5 py-3 text-zinc-400">{new Date(l.expiresAt).toLocaleDateString()}</td>
                  <td className="px-5 py-3">
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                      !l.isActive ? 'bg-zinc-800 text-zinc-500' :
                      expired ? 'bg-red-500/15 text-red-400' :
                      'bg-emerald-500/15 text-emerald-400'
                    }`}>
                      {!l.isActive ? 'REVOKED' : expired ? 'EXPIRED' : 'ACTIVE'}
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {licenses.length === 0 && (
          <div className="p-8 text-center text-zinc-600 text-sm italic">No licenses yet</div>
        )}
      </div>
    </div>
  );
}
