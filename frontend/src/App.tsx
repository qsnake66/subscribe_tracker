import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        {/* Placeholder for Register */}
        <Route path="/register" element={<div className="text-white flex items-center justify-center min-h-screen">Registration Page Coming Soon</div>} />
      </Routes>
    </Router>
  );
}

export default App;
