# Frontend Rules - LicenseManager

## Stack
- React 18 + Vite
- Tailwind CSS v3
- React Router v6
- Lucide React (icons)
- No external state library (useState + fetch is enough)

## Layout
```
frontend/src/
|-- App.jsx              # Routes: /login, / (Layout), /apps, /licenses
|-- main.jsx             # ReactDOM root, BrowserRouter
|-- index.css            # Tailwind + scrollbar styles
|-- lib/
|   |-- api.js           # All fetch calls, token management
|-- components/
|   |-- Layout.jsx       # Sidebar nav + <Outlet />
|-- pages/
    |-- auth/Login.jsx
    |-- admin/
        |-- Dashboard.jsx  # Stats + recent licenses
        |-- Apps.jsx       # App list + create form
        |-- Licenses.jsx   # License list + issue form + HWID expand
```

## API Client Pattern (`lib/api.js`)
```js
import { api } from '../lib/api';

// All calls throw Error on failure
const apps = await api.apps.list();
const license = await api.licenses.create({ appId, plan, days });
```

## Color Palette (dark theme)
```
bg:       #0a0a0a (page), #0f0f0f (panels)
border:   zinc-800
text:     white (headings), zinc-300 (body), zinc-500 (muted)
accent:   blue-600 (primary action)
status:
  active:  emerald-400 on emerald-500/15
  revoked: zinc-500 on zinc-800
  expired: red-400 on red-500/15
```

## Component Pattern
```jsx
export default function PageName() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    api.resource.list().then(setData).finally(() => setLoading(false));
  }, []);

  // render
}
```

## Adding a New Page
1. Create `src/pages/admin/NewPage.jsx`
2. Add route in `App.jsx` inside the Layout route
3. Add nav link in `components/Layout.jsx`
4. Add API method in `lib/api.js` if needed
5. Use `/generate-component` skill for scaffolding
