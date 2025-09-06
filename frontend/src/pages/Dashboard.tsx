import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getUserID } from '../services/auth';
import { Dungeon, DungeonSummary } from '../types/dungeon';
import QuestRow from '../components/quest/QuestRow';
import BalanceChip from '../components/layout/BalanceChip';
import DungeonSelector from '../components/layout/DungeonSelector';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [userId] = useState<string | null>(getUserID());
  const [dungeons, setDungeons] = useState<Dungeon[]>([]);
  const [selectedDungeon, setSelectedDungeon] = useState<Dungeon | null>(null);
  const [summary, setSummary] = useState<DungeonSummary | null>(null);
  const [loading, setLoading] = useState(true);

  // Mock data for now
  useEffect(() => {
    if (!userId) {
      navigate('/');
      return;
    }

    // Mock dungeons data
    const mockDungeons: Dungeon[] = [
      {
        id: 'dungeon-1',
        title: 'Productivity Palace',
        admin_user_id: 123,
        created_at: new Date().toISOString(),
      },
      {
        id: 'dungeon-2',
        title: 'Focus Fortress',
        admin_user_id: 123,
        created_at: new Date().toISOString(),
      },
    ];

    setDungeons(mockDungeons);
    setSelectedDungeon(mockDungeons[0]);

    // Mock summary data
    const mockSummary: DungeonSummary = {
      balance: '1250',
      todays_quests: [
        {
          id: 'quest-1',
          title: 'Morning Excercise',
          description: 'Start your day with 10 minutes of stretches',
          category: 'daily',
          difficulty: 'easy',
          mode: 'BINARY',
          points_award: '50',
          cooldown_sec: 3600,
          streak_enabled: true,
          status: 'active',
        },
        {
          id: 'quest-2',
          title: 'Deep Work Session',
          description: 'Complete a 90-minute focused work session',
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
          title: 'Wash Dishes',
          description: 'Wash as many as you can',
          category: 'daily',
          difficulty: 'medium',
          mode: 'PER_MINUTE',
          points_award: '0',
          rate_points_per_min: '2.5',
          min_minutes: 10,
          max_minutes: 120,
          cooldown_sec: 3600,
          streak_enabled: true,
          status: 'active',
        },
      ],
      recent_purchases: [
        {
          id: 1,
          name: 'Extra Life',
          price: '200',
          purchased_at: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
        },
        {
          id: 2,
          name: 'Focus Boost',
          price: '150',
          purchased_at: new Date(Date.now() - 172800000).toISOString(), // 2 days ago
        },
        {
          id: 3,
          name: 'Time Extension',
          price: '300',
          purchased_at: new Date(Date.now() - 259200000).toISOString(), // 3 days ago
        },
      ],
    };

    setSummary(mockSummary);
    setLoading(false);
  }, [userId, navigate]);

  const handleDungeonChange = (dungeonId: string) => {
    const dungeon = dungeons.find(d => d.id === dungeonId);
    if (dungeon) {
      setSelectedDungeon(dungeon);
      // In a real app, we would refetch data here
    }
  };

  const handleQuestComplete = (questId: string, completion: any) => {
    // In a real app, we would call the API here
    console.log('Completing quest:', questId, completion);
    
    // Optimistically remove from list
    if (summary) {
      setSummary({
        ...summary,
        todays_quests: summary.todays_quests.filter(q => q.id !== questId),
      });
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-900">
        <div className="text-slate-100">Loading dashboard...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <div className="sticky top-0 z-10 border-b border-slate-700 bg-slate-900/80 backdrop-blur">
        <div className="flex items-center justify-between px-4 py-3">
          <div className="flex items-center space-x-2">
            <div className="h-8 w-8 rounded-lg bg-violet-600"></div>
            <h1 className="text-xl font-bold">ADHD Game Bot</h1>
          </div>
          
          <div className="flex items-center space-x-4">
            <DungeonSelector 
              dungeons={dungeons} 
              selectedDungeon={selectedDungeon} 
              onDungeonChange={handleDungeonChange} 
            />
            <BalanceChip balance={summary?.balance || '0'} />
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-4xl mx-auto px-4 py-6">
        <div className="mb-8">
          <h2 className="text-2xl font-bold mb-6">Today's Quests</h2>
          <div className="space-y-4">
            {summary?.todays_quests.map((quest) => (
              <QuestRow 
                key={quest.id} 
                quest={quest} 
                onComplete={handleQuestComplete} 
              />
            ))}
            
            {summary?.todays_quests.length === 0 && (
              <div className="card text-center py-12">
                <p className="text-slate-400">No quests available for today</p>
              </div>
            )}
          </div>
        </div>

        <div className="mb-8">
          <h2 className="text-2xl font-bold mb-6">Recent Rewards</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {summary?.recent_purchases.map((purchase: any) => (
              <div key={purchase.id} className="card">
                <h3 className="font-medium">{purchase.name}</h3>
                <p className="text-sm text-slate-400 mt-1">
                  {new Date(purchase.purchased_at).toLocaleDateString()}
                </p>
                <div className="mt-2 text-violet-400 font-medium">
                  -{purchase.price} points
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="card">
          <h2 className="text-xl font-bold mb-4">Keep Going!</h2>
          <p className="text-slate-300">
            You're doing great! Completing quests consistently helps build positive habits. 
            Remember to take breaks and celebrate your progress.
          </p>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;