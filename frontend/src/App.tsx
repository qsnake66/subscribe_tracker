import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';
import { getAuthToken } from './lib/auth';

import type { ReactNode } from 'react';

function RequireAuth({ children }: { children: ReactNode }) {
  const token = getAuthToken();
  if (!token) {
    return <Navigate to="/" replace />;
  }
  return children;
}

function App() {
  const isAuthed = Boolean(getAuthToken());

  return (
    <Router>
      <Routes>
        <Route path="/" element={isAuthed ? <Navigate to="/app" replace /> : <LoginPage />} />
        <Route path="/register" element={isAuthed ? <Navigate to="/app" replace /> : <RegisterPage />} />
        <Route
          path="/app"
          element={
            <RequireAuth>
              <DashboardPage />
            </RequireAuth>
          }
        />
        <Route path="*" element={<Navigate to={isAuthed ? '/app' : '/'} replace />} />
      </Routes>
    </Router>
  );
}

export default App;
