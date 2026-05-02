import { useState, useEffect } from 'react';
import { Plus, Trash2, Monitor, ChevronDown, ChevronUp, ToggleLeft, ToggleRight } from 'lucide-react';
import { api } from '../../lib/api';

const statusLabel = (l) => {
  if (!l.isActive) return { label: 'REVOKED', cls: 'bg-zinc-800 text-zinc-500' };
  if (new Date(l.expiresAt) <= new Date()) return { label: 'EXPIRED', cls: 'bg-red-500/15 text-red-400' };
  return { label: 'ACTIVE', cls: 'bg-emerald-500/15 text-emerald-400' };
};

export default function Licenses() {
  const [apps, setApps] = useState([]);
  const [licenses, setLicenses] = useState([]);
  const [selectedApp, setSelectedApp] = useState('');
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [expandedHWIDs, setExpandedHWIDs] = useState({});
  const [hwids, setHWIDs] = useState({});
  const [addingHwidTo, setAddingHwidTo] = useState(null);
  const [newHwid, setNewHwid] = useState({ hwid: '', deviceName: '' });
  const [editingHwidFor, setEditingHwidFor] = useState(null);
  const [editHwidValue, setEditHwidValue] = useState('');
  const [error, setError] = useState('');
  const [form, setForm] = useState({ appId: '', plan: 'basic', maxThreads: 5, maxHwid: 2, days: 30, note: '' });

  useEffect(() => {
    api.apps.list().then((a) => {
      setApps(a || []);
      if (a && a.length > 0) setForm((f) => ({ ...f, appId: a[0].id, maxHwid: a[0].maxHwidPerLicense ?? 2 }));
    });
  }, []);

  const loadLicenses = (appId) => {
    setLoading(true);
    api.licenses.list(appId).then((l) => setLicenses(l || [])).finally(() => setLoading(false));
  };

  useEffect(() => { loadLicenses(selectedApp); }, [selectedApp]);

  const handleCreate = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await api.licenses.create({ ...form, maxThreads: Number(form.maxThreads), maxHwid: Number(form.maxHwid), days: Number(form.days) });
      setCreating(false);
      loadLicenses(selectedApp);
    } catch (err) {
      setError(err.message);
    }
  };

  const toggleActive = async (l) => {
    await api.licenses.update(l.id, { isActive: !l.isActive });
    loadLicenses(selectedApp);
  };

  const deleteLicense = async (id) => {
    if (!confirm('Delete this license?')) return;
    await api.licenses.delete(id);
    loadLicenses(selectedApp);
  };

  const toggleHWIDs = async (licenseId) => {
    if (expandedHWIDs[licenseId]) {
      setExpandedHWIDs((p) => ({ ...p, [licenseId]: false }));
      return;
    }
    const data = await api.licenses.hwids(licenseId);
    setHWIDs((p) => ({ ...p, [licenseId]: data || [] }));
    setExpandedHWIDs((p) => ({ ...p, [licenseId]: true }));
  };

  const deleteHWID = async (hwidId, licenseId) => {
    await api.hwids.delete(hwidId);
    const data = await api.licenses.hwids(licenseId);
    setHWIDs((p) => ({ ...p, [licenseId]: data || [] }));
  };

  const handleSaveMaxHwid = async (l) => {
    const val = Number(editHwidValue);
    if (!isNaN(val) && val > 0) {
      await api.licenses.update(l.id, { maxHwid: val });
      loadLicenses(selectedApp);
    }
    setEditingHwidFor(null);
  };

  const handleAddHwid = async (licenseId) => {
    if (!newHwid.hwid) return;
    try {
      await api.licenses.addHwid(licenseId, newHwid);
      const data = await api.licenses.hwids(licenseId);
      setHWIDs((p) => ({ ...p, [licenseId]: data || [] }));
      setAddingHwidTo(null);
      setNewHwid({ hwid: '', deviceName: '' });
    } catch (err) {
      alert(err.message);
    }
  };

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-4">
          <h1 className="text-lg font-bold text-white">Licenses</h1>
          <select
            value={selectedApp}
            onChange={(e) => setSelectedApp(e.target.value)}
            className="bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-1.5 text-sm text-zinc-300 focus:outline-none"
          >
            <option value="">All Apps</option>
            {apps.map((a) => <option key={a.id} value={a.id}>{a.name}</option>)}
          </select>
        </div>
        <button
          onClick={() => setCreating(true)}
          className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-4 py-2 rounded-lg"
        >
          <Plus className="w-4 h-4" /> Issue License
        </button>
      </div>

      {creating && (
        <form
          onSubmit={handleCreate}
          className="bg-[#0f0f0f] border border-zinc-800 rounded-xl p-6 mb-6 space-y-4"
        >
          <h2 className="text-sm font-semibold text-white">Issue New License</h2>
          {error && <p className="text-red-400 text-sm">{error}</p>}
          <div className="grid grid-cols-6 gap-4">
            <div>
              <label className="text-xs text-zinc-400 block mb-1">App</label>
              <select
                value={form.appId}
                onChange={(e) => setForm({ ...form, appId: e.target.value })}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none"
                required
              >
                {apps.map((a) => <option key={a.id} value={a.id}>{a.name}</option>)}
              </select>
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Plan</label>
              <select
                value={form.plan}
                onChange={(e) => setForm({ ...form, plan: e.target.value })}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none"
              >
                <option value="basic">Basic</option>
                <option value="pro">Pro</option>
                <option value="enterprise">Enterprise</option>
              </select>
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Max Threads</label>
              <input
                type="number" min="1" max="50"
                value={form.maxThreads}
                onChange={(e) => setForm({ ...form, maxThreads: e.target.value })}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none"
              />
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Max Machines</label>
              <input
                type="number" min="1" max="20"
                value={form.maxHwid}
                onChange={(e) => setForm({ ...form, maxHwid: e.target.value })}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none"
              />
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Duration (days)</label>
              <input
                type="number" min="1"
                value={form.days}
                onChange={(e) => setForm({ ...form, days: e.target.value })}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none"
              />
            </div>
            <div>
              <label className="text-xs text-zinc-400 block mb-1">Note (optional)</label>
              <input
                value={form.note}
                onChange={(e) => setForm({ ...form, note: e.target.value })}
                placeholder="Customer name..."
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-3 py-2 text-sm focus:outline-none"
              />
            </div>
          </div>
          <div className="flex gap-3">
            <button type="submit" className="bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-4 py-2 rounded-lg">
              Issue
            </button>
            <button type="button" onClick={() => { setCreating(false); setError(''); }}
              className="bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-sm px-4 py-2 rounded-lg">
              Cancel
            </button>
          </div>
        </form>
      )}

      <div className="bg-[#0f0f0f] border border-zinc-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-zinc-800 text-xs text-zinc-500">
              <th className="px-5 py-3 text-left">Key</th>
              <th className="px-5 py-3 text-left">Plan</th>
              <th className="px-5 py-3 text-left">Threads</th>
              <th className="px-5 py-3 text-left">Machines</th>
              <th className="px-5 py-3 text-left">Expires</th>
              <th className="px-5 py-3 text-left">Note</th>
              <th className="px-5 py-3 text-left">Status</th>
              <th className="px-5 py-3 text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            {licenses.map((l) => {
              const { label, cls } = statusLabel(l);
              const expanded = expandedHWIDs[l.id];
              return (
                <>
                  <tr key={l.id} className="border-b border-zinc-800/50 hover:bg-zinc-900/30">
                    <td className="px-5 py-3 font-mono text-xs text-white">{l.key}</td>
                    <td className="px-5 py-3 text-zinc-400">{l.plan}</td>
                    <td className="px-5 py-3 text-zinc-400">{l.maxThreads}</td>
                    <td className="px-5 py-3 text-zinc-400">
                      {editingHwidFor === l.id ? (
                        <div className="flex items-center gap-1">
                          <input
                            type="number"
                            min="1"
                            value={editHwidValue}
                            onChange={(e) => setEditHwidValue(e.target.value)}
                            className="w-16 bg-zinc-900 border border-zinc-700 rounded px-2 py-0.5 text-xs text-white outline-none"
                            autoFocus
                            onBlur={() => handleSaveMaxHwid(l)}
                            onKeyDown={(e) => e.key === 'Enter' && handleSaveMaxHwid(l)}
                          />
                        </div>
                      ) : (
                        <span
                          className="cursor-pointer hover:text-white border-b border-dashed border-zinc-600 pb-0.5"
                          onClick={() => {
                            setEditingHwidFor(l.id);
                            setEditHwidValue(l.maxHwid.toString());
                          }}
                          title="Click to edit max machines"
                        >
                          {l.maxHwid}
                        </span>
                      )}
                    </td>
                    <td className="px-5 py-3 text-zinc-400 text-xs">{new Date(l.expiresAt).toLocaleDateString()}</td>
                    <td className="px-5 py-3 text-zinc-500 text-xs max-w-[120px] truncate">{l.note}</td>
                    <td className="px-5 py-3">
                      <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${cls}`}>{label}</span>
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-1">
                        <button onClick={() => toggleActive(l)} title={l.isActive ? 'Revoke' : 'Activate'}
                          className="p-1.5 hover:bg-zinc-800 rounded transition-colors">
                          {l.isActive
                            ? <ToggleRight className="w-4 h-4 text-emerald-500" />
                            : <ToggleLeft className="w-4 h-4 text-zinc-500" />}
                        </button>
                        <button onClick={() => toggleHWIDs(l.id)} title="Devices"
                          className="p-1.5 hover:bg-zinc-800 rounded transition-colors">
                          <Monitor className="w-4 h-4 text-zinc-500" />
                          {expanded ? <ChevronUp className="w-3 h-3 inline text-zinc-500" /> : <ChevronDown className="w-3 h-3 inline text-zinc-500" />}
                        </button>
                        <button onClick={() => deleteLicense(l.id)} title="Delete"
                          className="p-1.5 hover:bg-zinc-800 rounded transition-colors">
                          <Trash2 className="w-4 h-4 text-red-500" />
                        </button>
                      </div>
                    </td>
                  </tr>
                  {expanded && (
                    <tr key={`hwid-${l.id}`} className="bg-zinc-900/20">
                      <td colSpan={8} className="px-8 py-3">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-xs font-semibold text-zinc-400 uppercase tracking-wider">Registered Devices</span>
                          <button
                            onClick={() => setAddingHwidTo(addingHwidTo === l.id ? null : l.id)}
                            className="text-xs text-blue-400 hover:text-blue-300 flex items-center gap-1"
                          >
                            <Plus className="w-3 h-3" /> Manually Add Device
                          </button>
                        </div>

                        {addingHwidTo === l.id && (
                          <div className="flex items-center gap-2 mb-3 bg-zinc-900 p-2 rounded border border-zinc-800">
                            <input
                              value={newHwid.hwid}
                              onChange={(e) => setNewHwid({ ...newHwid, hwid: e.target.value })}
                              placeholder="Hardware ID (e.g. mac address)"
                              className="bg-[#0f0f0f] border border-zinc-700 rounded px-2 py-1 text-xs text-white focus:outline-none flex-1 font-mono"
                            />
                            <input
                              value={newHwid.deviceName}
                              onChange={(e) => setNewHwid({ ...newHwid, deviceName: e.target.value })}
                              placeholder="Device Name (optional)"
                              className="bg-[#0f0f0f] border border-zinc-700 rounded px-2 py-1 text-xs text-white focus:outline-none w-48"
                            />
                            <button
                              onClick={() => handleAddHwid(l.id)}
                              className="bg-blue-600 hover:bg-blue-500 text-white text-xs px-3 py-1 rounded font-medium"
                            >
                              Add
                            </button>
                          </div>
                        )}

                        {(hwids[l.id] || []).length === 0 ? (
                          <span className="text-xs text-zinc-600 italic">No devices registered yet</span>
                        ) : (
                          <div className="space-y-1">
                            {hwids[l.id].map((h) => (
                              <div key={h.id} className="flex items-center justify-between text-xs">
                                <div className="flex items-center gap-3">
                                  <Monitor className="w-3.5 h-3.5 text-zinc-500" />
                                  <span className="text-zinc-300">{h.deviceName}</span>
                                  <span className="font-mono text-zinc-600">{h.hwid.slice(0, 16)}...</span>
                                  <span className="text-zinc-600">Last seen: {new Date(h.lastSeenAt).toLocaleDateString()}</span>
                                </div>
                                <button onClick={() => deleteHWID(h.id, l.id)}
                                  className="p-1 hover:bg-zinc-800 rounded text-red-500">
                                  <Trash2 className="w-3.5 h-3.5" />
                                </button>
                              </div>
                            ))}
                          </div>
                        )}
                      </td>
                    </tr>
                  )}
                </>
              );
            })}
          </tbody>
        </table>
        {licenses.length === 0 && !loading && (
          <div className="p-8 text-center text-zinc-600 text-sm italic">No licenses yet</div>
        )}
      </div>
    </div>
  );
}
