import React from 'react';
import { getUserID } from '../services/auth';

const Profile: React.FC = () => {
  const userId = getUserID();

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <div className="sticky top-0 z-10 border-b border-slate-700 bg-slate-900/80 backdrop-blur">
        <div className="max-w-4xl mx-auto px-4 py-3">
          <h1 className="text-2xl font-bold">Profile</h1>
        </div>
      </div>

      {/* Profile Content */}
      <div className="max-w-4xl mx-auto px-4 py-6">
        <div className="card">
          <div className="flex items-center space-x-4 mb-6">
            <div className="h-16 w-16 rounded-full bg-violet-600 flex items-center justify-center">
              <span className="text-xl font-bold">U</span>
            </div>
            <div>
              <h2 className="text-xl font-bold">User {userId || 'Unknown'}</h2>
              <p className="text-slate-400">Balance: 1,250 points</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="card">
              <h3 className="font-medium mb-2">Timezone</h3>
              <p className="text-slate-300">America/New_York</p>
            </div>

            <div className="card">
              <h3 className="font-medium mb-2">Member Since</h3>
              <p className="text-slate-300">January 15, 2024</p>
            </div>
          </div>
        </div>

        <div className="card mt-6">
          <h2 className="text-xl font-bold mb-4">Achievements</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            {[1, 2, 3, 4, 5, 6].map((item) => (
              <div key={item} className="border border-slate-700 rounded-2xl p-4 text-center">
                <div className="h-12 w-12 rounded-full bg-slate-800 mx-auto flex items-center justify-center mb-2">
                  <span className="text-lg">🏆</span>
                </div>
                <h3 className="font-medium">Achievement {item}</h3>
                <p className="text-sm text-slate-400 mt-1">Unlock condition</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile;