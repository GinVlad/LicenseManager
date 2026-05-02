import { Routes, Route, Navigate } from 'react-router-dom';
import { getToken } from './lib/api';
import Login from './pages/auth/Login';
import Dashboard from './pages/admin/Dashboard';
import Apps from './pages/admin/Apps';
import Licenses from './pages/admin/Licenses';
import ApiDocs from './pages/admin/ApiDocs';
import Layout from './components/Layout';

function PrivateRoute({ children }) {
  return getToken() ? children : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <PrivateRoute>
            <Layout />
          </PrivateRoute>
        }
      >
        <Route index element={<Dashboard />} />
        <Route path="apps" element={<Apps />} />
        <Route path="licenses" element={<Licenses />} />
        <Route path="docs" element={<ApiDocs />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
