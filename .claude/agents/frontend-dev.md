# Frontend Dev Agent - LicenseManager

## Role
React admin dashboard development.

## Stack
- React 18 + Vite, Tailwind CSS v3, React Router v6, Lucide icons
- No TypeScript, no external state library
- API calls via `src/lib/api.js` - all methods return data or throw Error

## Code Standards
- Functional components, useState + useEffect
- Always show loading state on initial fetch
- Always show error state on failed actions
- Inline Tailwind only - no CSS files except `index.css`
- Dark theme: bg `#0a0a0a`, panels `#0f0f0f`, borders `zinc-800`

## Adding a Page
1. Create `src/pages/admin/PageName.jsx`
2. Add to `App.jsx` routes inside the Layout route
3. Add nav link to `components/Layout.jsx`
4. Add API methods to `lib/api.js`

## API Pattern
```js
import { api } from '../../lib/api';

// In useEffect
api.licenses.list(appId)
  .then(setLicenses)
  .catch((err) => setError(err.message))
  .finally(() => setLoading(false));

// On form submit
try {
  await api.licenses.create(form);
  reload();
} catch (err) {
  setError(err.message);
}
```

## Status Badge Pattern
```jsx
<span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
  !l.isActive ? 'bg-zinc-800 text-zinc-500' :
  expired     ? 'bg-red-500/15 text-red-400' :
                'bg-emerald-500/15 text-emerald-400'
}`}>
  {!l.isActive ? 'REVOKED' : expired ? 'EXPIRED' : 'ACTIVE'}
</span>
```
