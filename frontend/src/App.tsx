import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Quests from './pages/Quests';
import Shop from './pages/Shop';
import Profile from './pages/Profile';
import Admin from './pages/Admin';
import { getUserID, setUserID } from './services/auth';

const App: React.FC = () => {
  const [userId, setUserId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const storedUserId = getUserID();
    if (storedUserId) {
      setUserId(storedUserId);
      setLoading(false);
    } else {
      // Prompt for user ID if not set
      const inputUserId = prompt('Please enter your User ID:');
      if (inputUserId) {
        setUserID(inputUserId);
        setUserId(inputUserId);
      }
      setLoading(false);
    }
  }, []);

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-900">
        <div className="text-slate-100">Loading...</div>
      </div>
    );
  }

  if (!userId) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-900">
        <div className="text-slate-100">
          <h1 className="text-2xl font-bold mb-4">User ID Required</h1>
          <p>Please refresh the page and enter your User ID</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100">
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/quests" element={<Quests />} />
        <Route path="/shop" element={<Shop />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/admin" element={<Admin />} />
      </Routes>
    </div>
  );
};

export default App;