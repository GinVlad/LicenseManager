import { useState, useEffect } from 'react';
import { Plus, Copy, Pencil, Trash2 } from 'lucide-react';
import { api } from '../../lib/api';

export default function Apps() {
  const [apps, setApps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [form, setForm] = useState({ name: '', slug: '', maxHwidPerLicense: 2 });
  const [error, setError] = useState('');

  const [editing, setEditing] = useState(null);

  const load = () => api.apps.list().then((a) => setApps(a || []));

  useEffect(() => { load().finally(() => setLoading(false)); }, []);

  const handleCreate = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await api.apps.create({ ...form, maxHwidPerLicense: Number(form.maxHwidPerLicense) });
      setCreating(false);
      setForm({ name: '', slug: '', maxHwidPerLicense: 2 });
      load();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleUpdate = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await api.apps.update(editing, { ...form, maxHwidPerLicense: Number(form.maxHwidPerLicense) });
      setEditing(null);
      setForm({ name: '', slug: '', maxHwidPerLicense: 2 });
      load();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Delete this app and ALL its licenses? This cannot be undone.')) return;
    try {
      await api.apps.delete(id);
      load();
    } catch (err) {
      alert(err.message);
    }
  };

  const copyKey = (key) => navigator.clipboard.writeText(key);

  const startEdit = (a) => {
    setCreating(false);
    setEditing(a.id);
    setForm({ name: a.name, slug: a.slug, maxHwidPerLicense: a.maxHwidPerLicense });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-lg font-bold text-white">Apps</h1>
        <button
          onClick={() => setCreating(true)}
          className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-4 py-2 rounded-lg transition-colors"
        >
          <Plus className="w-4 h-4" /> New App
        </button>
      </div>

      {(creating || editing) && (
        <form
          onSubmit={editing ? handleUpdate : handleCreate}
          className="bg-[#0f0f0f] border border-zinc-800 rounded-xl p-6 mb-6 space-y-4"
        >
          <h2 className="text-sm font-semibold text-white">{editing ? 'Edit App' : 'Create App'}</h2>
          {error && <p className="text-red-400 text-sm">{error}</p>}
          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Name</label>
              <input
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="eBay Creator"
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500/50"
                required
              />
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Slug</label>
              <input
                value={form.slug}
                onChange={(e) => setForm({ ...form, slug: e.target.value.toLowerCase().replace(/\s+/g, '-') })}
                placeholder="ebay-creator"
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500/50"
                required
              />
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Max HWIDs per license</label>
              <input
                type="number"
                min="1"
                max="10"
                value={form.maxHwidPerLicense}
                onChange={(e) => setForm({ ...form, maxHwidPerLicense: e.target.value })}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500/50"
              />
            </div>
          </div>
          <div className="flex gap-3">
            <button type="submit" className="bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-4 py-2 rounded-lg">
              {editing ? 'Save Changes' : 'Create'}
            </button>
            <button
              type="button"
              onClick={() => { setCreating(false); setEditing(null); setError(''); setForm({ name: '', slug: '', maxHwidPerLicense: 2 }); }}
              className="bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-sm px-4 py-2 rounded-lg"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      <div className="bg-[#0f0f0f] border border-zinc-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-zinc-800 text-xs text-zinc-500">
              <th className="px-5 py-3 text-left">Name</th>
              <th className="px-5 py-3 text-left">Slug</th>
              <th className="px-5 py-3 text-left">Max HWIDs</th>
              <th className="px-5 py-3 text-left">API Key</th>
              <th className="px-5 py-3 text-left">Created</th>
              <th className="px-5 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {apps.map((a) => (
              <tr key={a.id} className={`border-b border-zinc-800/50 hover:bg-zinc-900/30 ${editing === a.id ? 'bg-zinc-900/30' : ''}`}>
                <td className="px-5 py-3 font-medium text-white">{a.name}</td>
                <td className="px-5 py-3 font-mono text-xs text-zinc-400">{a.slug}</td>
                <td className="px-5 py-3 text-zinc-400">{a.maxHwidPerLicense}</td>
                <td className="px-5 py-3">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-[11px] text-zinc-500">{a.apiKey.slice(0, 16)}...</span>
                    <button onClick={() => copyKey(a.apiKey)} className="p-1 hover:bg-zinc-800 rounded text-zinc-400 hover:text-white transition-colors" title="Copy API Key">
                      <Copy className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
                <td className="px-5 py-3 text-zinc-500 text-xs">{new Date(a.createdAt).toLocaleDateString()}</td>
                <td className="px-5 py-3 text-right">
                  <div className="flex items-center justify-end gap-1">
                    <button onClick={() => startEdit(a)} title="Edit App"
                      className="p-1.5 hover:bg-zinc-800 rounded transition-colors">
                      <Pencil className="w-4 h-4 text-zinc-400 hover:text-blue-400" />
                    </button>
                    <button onClick={() => handleDelete(a.id)} title="Delete App"
                      className="p-1.5 hover:bg-zinc-800 rounded transition-colors">
                      <Trash2 className="w-4 h-4 text-zinc-400 hover:text-red-500" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {apps.length === 0 && !loading && (
          <div className="p-8 text-center text-zinc-600 text-sm italic">No apps yet — create your first one</div>
        )}
      </div>
    </div>
  );
}
