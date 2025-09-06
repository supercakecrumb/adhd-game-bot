import React, { useState } from 'react';
import { Quest } from '../types/quest';

const Quests: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly' | 'adhoc'>('daily');
  
  // Mock quests data
  const mockQuests: Quest[] = [
    {
      id: 'quest-1',
      title: 'Morning Routine',
      description: 'Complete your morning routine checklist',
      category: 'daily',
      difficulty: 'easy',
      mode: 'BINARY',
      points_award: '30',
      cooldown_sec: 3600,
      streak_enabled: true,
      status: 'active',
    },
    {
      id: 'quest-2',
      title: 'Deep Work Block',
      description: 'Focus on a single task for 90 minutes',
      category: 'daily',
      difficulty: 'medium',
      mode: 'PARTIAL',
      points_award: '100',
      cooldown_sec: 7200,
      streak_enabled: true,
      status: 'active',
    },
    {
      id: 'quest-3',
      title: 'Weekly Review',
      description: 'Review your week and plan for next week',
      category: 'weekly',
      difficulty: 'hard',
      mode: 'BINARY',
      points_award: '150',
      cooldown_sec: 86400,
      streak_enabled: false,
      status: 'active',
    },
    {
      id: 'quest-4',
      title: 'Learn Something New',
      description: 'Spend time learning a new skill or topic',
      category: 'adhoc',
      difficulty: 'medium',
      mode: 'PER_MINUTE',
      points_award: '0',
      rate_points_per_min: '3',
      min_minutes: 15,
      max_minutes: 120,
      cooldown_sec: 3600,
      streak_enabled: false,
      status: 'active',
    },
  ];

  const filteredQuests = mockQuests.filter(quest => quest.category === activeTab);

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <div className="sticky top-0 z-10 border-b border-slate-700 bg-slate-900/80 backdrop-blur">
        <div className="max-w-4xl mx-auto px-4 py-3">
          <h1 className="text-2xl font-bold">Quests</h1>
        </div>
      </div>

      {/* Tabs */}
      <div className="max-w-4xl mx-auto px-4 py-6">
        <div className="flex space-x-4 border-b border-slate-700 mb-6">
          <button
            onClick={() => setActiveTab('daily')}
            className={`pb-3 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'daily'
                ? 'border-violet-500 text-violet-400'
                : 'border-transparent text-slate-400 hover:text-slate-300'
            }`}
          >
            Daily
          </button>
          <button
            onClick={() => setActiveTab('weekly')}
            className={`pb-3 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'weekly'
                ? 'border-violet-500 text-violet-400'
                : 'border-transparent text-slate-400 hover:text-slate-300'
            }`}
          >
            Weekly
          </button>
          <button
            onClick={() => setActiveTab('adhoc')}
            className={`pb-3 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'adhoc'
                ? 'border-violet-500 text-violet-400'
                : 'border-transparent text-slate-400 hover:text-slate-300'
            }`}
          >
            One-time
          </button>
        </div>

        {/* Quest List */}
        <div className="space-y-4">
          {filteredQuests.map((quest) => (
            <div key={quest.id} className="card">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <h3 className="font-medium">{quest.title}</h3>
                    <span className={`badge-difficulty-${quest.difficulty}`}>
                      {quest.difficulty.charAt(0).toUpperCase() + quest.difficulty.slice(1)}
                    </span>
                    <span className={`badge-mode-${quest.mode.toLowerCase().replace('_', '-')}`}>
                      {quest.mode === 'BINARY' ? 'Binary' : 
                       quest.mode === 'PARTIAL' ? 'Partial' : 'Per Minute'}
                    </span>
                  </div>
                  <p className="text-sm text-slate-400 mb-3">{quest.description}</p>
                  <div className="text-sm">
                    <span className="text-slate-300">Reward: </span>
                    <span className="font-medium text-violet-400">
                      {quest.mode === 'BINARY' ? `${quest.points_award} points` :
                       quest.mode === 'PARTIAL' ? `Up to ${quest.points_award} points` :
                       quest.rate_points_per_min ? `${quest.rate_points_per_min} points/minute` : '0 points'}
                    </span>
                  </div>
                </div>
                
                <button className="btn-primary whitespace-nowrap">
                  {quest.mode === 'BINARY' ? 'Complete' : 
                   quest.mode === 'PARTIAL' ? 'Log Progress' : 'Log Time'}
                </button>
              </div>
            </div>
          ))}
          
          {filteredQuests.length === 0 && (
            <div className="card text-center py-12">
              <p className="text-slate-400">No quests available in this category</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Quests;