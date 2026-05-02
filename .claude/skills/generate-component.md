# Skill: Generate React Component

## Usage
`/generate-component [PageName] [description]`

## Example
`/generate-component Users Admin page to list and manage customer user accounts`

## What It Does
1. Creates `frontend/src/pages/admin/[PageName].jsx`
2. Adds route to `frontend/src/App.jsx`
3. Adds nav link to `frontend/src/components/Layout.jsx`
4. Adds API methods to `frontend/src/lib/api.js`

## Template
```jsx
import { useState, useEffect } from 'react';
import { api } from '../../lib/api';

export default function PageName() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const load = () =>
    api.resource.list()
      .then(setData)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));

  useEffect(() => { load(); }, []);

  if (loading) return <div className="p-8 text-zinc-500 text-sm">Loading...</div>;

  return (
    <div className="p-8">
      <h1 className="text-lg font-bold text-white mb-6">Page Title</h1>
      {error && <p className="text-red-400 text-sm mb-4">{error}</p>}
      {/* table or card list */}
    </div>
  );
}
```
