# Phase 3: Web Interface Integration

## Overview
Connect the React frontend with the backend APIs to create a fully functional web interface for quest management. Users will authenticate via Telegram and manage their quests through a rich web UI.

## Goals
- Integrate frontend authentication with Telegram auth flow
- Connect all frontend pages to real backend APIs
- Implement quest creation, editing, and completion via web
- Add dungeon management and member invitation
- Ensure responsive design works on mobile devices

## Tasks

### 3.1 Frontend Authentication Integration
**File**: `frontend/src/services/auth.ts` (modify existing)

```typescript
interface AuthUser {
  id: number;
  telegram_user_id?: number;
  username: string;
  email?: string;
  balance: string;
  timezone: string;
}

interface AuthResponse {
  user: AuthUser;
  session_token: string;
  expires_at: string;
}

class AuthService {
  private baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  
  async validateSession(): Promise<AuthUser | null> {
    try {
      const response = await fetch(`${this.baseURL}/auth/validate`, {
        credentials: 'include', // Include cookies
      });
      
      if (!response.ok) {
        return null;
      }
      
      return await response.json();
    } catch (error) {
      console.error('Session validation failed:', error);
      return null;
    }
  }
  
  async logout(): Promise<void> {
    try {
      await fetch(`${this.baseURL}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Logout failed:', error);
    }
    
    // Clear local storage
    localStorage.removeItem('user_id');
    
    // Redirect to login page
    window.location.href = '/login';
  }
  
  getTelegramLoginURL(): string {
    return `${this.baseURL}/auth/telegram/login`;
  }
}

export const authService = new AuthService();
```

**Testing**:
- Test session validation on app load
- Test logout functionality
- Test Telegram login URL generation
- Test cookie-based authentication

### 3.2 Update App.tsx with Real Authentication
**File**: `frontend/src/App.tsx` (modify existing)

```typescript
import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Quests from './pages/Quests';
import Shop from './pages/Shop';
import Profile from './pages/Profile';
import Admin from './pages/Admin';
import Login from './pages/Login';
import { authService } from './services/auth';

interface AuthUser {
  id: number;
  username: string;
  balance: string;
  timezone: string;
}

const App: React.FC = () => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [authChecked, setAuthChecked] = useState(false);

  useEffect(() => {
    checkAuthentication();
  }, []);

  const checkAuthentication = async () => {
    try {
      const validatedUser = await authService.validateSession();
      setUser(validatedUser);
    } catch (error) {
      console.error('Auth check failed:', error);
      setUser(null);
    } finally {
      setLoading(false);
      setAuthChecked(true);
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-900">
        <div className="text-slate-100">Loading...</div>
      </div>
    );
  }

  if (!user && authChecked) {
    return <Login />;
  }

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100">
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<Dashboard user={user} />} />
        <Route path="/quests" element={<Quests user={user} />} />
        <Route path="/shop" element={<Shop user={user} />} />
        <Route path="/profile" element={<Profile user={user} />} />
        <Route path="/admin" element={<Admin user={user} />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </div>
  );
};

export default App;
```

**Testing**:
- Test authentication flow on app load
- Test redirect to login when not authenticated
- Test user prop passing to components
- Test navigation after authentication

### 3.3 Create Login Page
**File**: `frontend/src/pages/Login.tsx` (new)

```typescript
import React from 'react';
import { authService } from '../services/auth';

const Login: React.FC = () => {
  const handleTelegramLogin = () => {
    window.location.href = authService.getTelegramLoginURL();
  };

  return (
    <div className="min-h-screen bg-slate-900 flex items-center justify-center">
      <div className="bg-slate-800 p-8 rounded-lg shadow-lg max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-slate-100 mb-2">
            🎮 ADHD Game Bot
          </h1>
          <p className="text-slate-400">
            Transform your tasks into an engaging game
          </p>
        </div>
        
        <div className="space-y-4">
          <button
            onClick={handleTelegramLogin}
            className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 rounded-lg transition-colors flex items-center justify-center space-x-2"
          >
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.568 8.16l-1.61 7.59c-.12.54-.44.67-.89.42l-2.46-1.82-1.19 1.14c-.13.13-.24.24-.49.24l.17-2.45 4.51-4.08c.2-.17-.04-.27-.3-.1l-5.57 3.51-2.4-.75c-.52-.16-.53-.52.11-.77l9.39-3.62c.43-.16.81.1.67.73z"/>
            </svg>
            <span>Login with Telegram</span>
          </button>
          
          <div className="text-center text-sm text-slate-400">
            <p>You need to have started the bot first:</p>
            <a 
              href="https://t.me/YourBotUsername" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-blue-400 hover:text-blue-300 underline"
            >
              @YourBotUsername
            </a>
          </div>
        </div>
        
        <div className="mt-8 text-center text-xs text-slate-500">
          <p>By logging in, you agree to our terms of service</p>
        </div>
      </div>
    </div>
  );
};

export default Login;
```

**Testing**:
- Test login page renders correctly
- Test Telegram login button redirects properly
- Test responsive design on mobile
- Test bot link opens correctly

### 3.4 Update API Service with Real Endpoints
**File**: `frontend/src/services/api.ts` (modify existing)

```typescript
import { Quest, CreateQuestRequest, CompleteQuestRequest } from '../types/quest';
import { Dungeon, CreateDungeonRequest } from '../types/dungeon';
import { User } from '../types/user';

class ApiService {
  private baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const response = await fetch(url, {
      ...options,
      credentials: 'include', // Include cookies for auth
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });
    
    if (!response.ok) {
      const error = await response.text();
      throw new Error(`API Error: ${response.status} - ${error}`);
    }
    
    return response.json();
  }
  
  // Quest endpoints
  async getQuests(dungeonId: string): Promise<Quest[]> {
    return this.request<Quest[]>(`/api/dungeons/${dungeonId}/quests`);
  }
  
  async createQuest(dungeonId: string, quest: CreateQuestRequest): Promise<Quest> {
    return this.request<Quest>(`/api/dungeons/${dungeonId}/quests`, {
      method: 'POST',
      body: JSON.stringify(quest),
    });
  }
  
  async updateQuest(questId: string, quest: Partial<Quest>): Promise<Quest> {
    return this.request<Quest>(`/api/quests/${questId}`, {
      method: 'PUT',
      body: JSON.stringify(quest),
    });
  }
  
  async completeQuest(questId: string, completion: CompleteQuestRequest): Promise<void> {
    return this.request<void>(`/api/quests/${questId}/complete`, {
      method: 'POST',
      body: JSON.stringify(completion),
    });
  }
  
  async deleteQuest(questId: string): Promise<void> {
    return this.request<void>(`/api/quests/${questId}`, {
      method: 'DELETE',
    });
  }
  
  // Dungeon endpoints
  async getUserDungeons(): Promise<Dungeon[]> {
    return this.request<Dungeon[]>('/api/dungeons');
  }
  
  async createDungeon(dungeon: CreateDungeonRequest): Promise<Dungeon> {
    return this.request<Dungeon>('/api/dungeons', {
      method: 'POST',
      body: JSON.stringify(dungeon),
    });
  }
  
  async getDungeon(dungeonId: string): Promise<Dungeon> {
    return this.request<Dungeon>(`/api/dungeons/${dungeonId}`);
  }
  
  async joinDungeon(inviteCode: string): Promise<Dungeon> {
    return this.request<Dungeon>('/api/dungeons/join', {
      method: 'POST',
      body: JSON.stringify({ invite_code: inviteCode }),
    });
  }
  
  // User endpoints
  async getCurrentUser(): Promise<User> {
    return this.request<User>('/api/user');
  }
  
  async updateUser(user: Partial<User>): Promise<User> {
    return this.request<User>('/api/user', {
      method: 'PUT',
      body: JSON.stringify(user),
    });
  }
  
  // Stats endpoints
  async getUserStats(): Promise<any> {
    return this.request<any>('/api/user/stats');
  }
  
  async getDungeonStats(dungeonId: string): Promise<any> {
    return this.request<any>(`/api/dungeons/${dungeonId}/stats`);
  }
}

export const apiService = new ApiService();
```

**Testing**:
- Test all API endpoints with real backend
- Test error handling for failed requests
- Test authentication headers are sent
- Test request/response data formats

### 3.5 Update Quest Management Page
**File**: `frontend/src/pages/Quests.tsx` (modify existing)

```typescript
import React, { useState, useEffect } from 'react';
import { Plus, Target, Clock, Award } from 'lucide-react';
import { Quest, CreateQuestRequest } from '../types/quest';
import { Dungeon } from '../types/dungeon';
import { apiService } from '../services/api';
import QuestRow from '../components/quest/QuestRow';
import CompletionDialog from '../components/quest/CompletionDialog';
import DungeonSelector from '../components/layout/DungeonSelector';

interface QuestsProps {
  user: any;
}

const Quests: React.FC<QuestsProps> = ({ user }) => {
  const [quests, setQuests] = useState<Quest[]>([]);
  const [dungeons, setDungeons] = useState<Dungeon[]>([]);
  const [selectedDungeon, setSelectedDungeon] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [completionQuest, setCompletionQuest] = useState<Quest | null>(null);
  
  const [newQuest, setNewQuest] = useState<CreateQuestRequest>({
    title: '',
    description: '',
    category: 'daily',
    difficulty: 'medium',
    mode: 'BINARY',
    points_award: 10,
    streak_enabled: true,
  });

  useEffect(() => {
    loadDungeons();
  }, []);

  useEffect(() => {
    if (selectedDungeon) {
      loadQuests();
    }
  }, [selectedDungeon]);

  const loadDungeons = async () => {
    try {
      const userDungeons = await apiService.getUserDungeons();
      setDungeons(userDungeons);
      if (userDungeons.length > 0) {
        setSelectedDungeon(userDungeons[0].id);
      }
    } catch (error) {
      console.error('Failed to load dungeons:', error);
    }
  };

  const loadQuests = async () => {
    if (!selectedDungeon) return;
    
    setLoading(true);
    try {
      const dungeonQuests = await apiService.getQuests(selectedDungeon);
      setQuests(dungeonQuests);
    } catch (error) {
      console.error('Failed to load quests:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateQuest = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedDungeon) return;

    try {
      const createdQuest = await apiService.createQuest(selectedDungeon, newQuest);
      setQuests([...quests, createdQuest]);
      setShowCreateForm(false);
      setNewQuest({
        title: '',
        description: '',
        category: 'daily',
        difficulty: 'medium',
        mode: 'BINARY',
        points_award: 10,
        streak_enabled: true,
      });
    } catch (error) {
      console.error('Failed to create quest:', error);
      alert('Failed to create quest. Please try again.');
    }
  };

  const handleCompleteQuest = (quest: Quest) => {
    setCompletionQuest(quest);
  };

  const handleQuestCompleted = async (questId: string, completion: any) => {
    try {
      await apiService.completeQuest(questId, completion);
      
      // Refresh quests to show updated state
      await loadQuests();
      
      setCompletionQuest(null);
      
      // Show success message
      alert('Quest completed successfully! Points have been awarded.');
    } catch (error) {
      console.error('Failed to complete quest:', error);
      alert('Failed to complete quest. Please try again.');
    }
  };

  if (loading && quests.length === 0) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-slate-400">Loading quests...</div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-slate-100">Quests</h1>
        <button
          onClick={() => setShowCreateForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center space-x-2"
        >
          <Plus className="w-4 h-4" />
          <span>Create Quest</span>
        </button>
      </div>

      <DungeonSelector
        dungeons={dungeons}
        selectedDungeon={selectedDungeon}
        onDungeonChange={setSelectedDungeon}
      />

      {/* Quest Creation Form */}
      {showCreateForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-slate-800 p-6 rounded-lg max-w-md w-full mx-4">
            <h2 className="text-xl font-bold mb-4">Create New Quest</h2>
            <form onSubmit={handleCreateQuest} className="space-y-4">
              <input
                type="text"
                placeholder="Quest title"
                value={newQuest.title}
                onChange={(e) => setNewQuest({...newQuest, title: e.target.value})}
                className="w-full p-2 bg-slate-700 rounded border border-slate-600"
                required
              />
              <textarea
                placeholder="Description (optional)"
                value={newQuest.description}
                onChange={(e) => setNewQuest({...newQuest, description: e.target.value})}
                className="w-full p-2 bg-slate-700 rounded border border-slate-600"
                rows={3}
              />
              <div className="grid grid-cols-2 gap-4">
                <select
                  value={newQuest.category}
                  onChange={(e) => setNewQuest({...newQuest, category: e.target.value as any})}
                  className="p-2 bg-slate-700 rounded border border-slate-600"
                >
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="adhoc">One-time</option>
                </select>
                <select
                  value={newQuest.difficulty}
                  onChange={(e) => setNewQuest({...newQuest, difficulty: e.target.value as any})}
                  className="p-2 bg-slate-700 rounded border border-slate-600"
                >
                  <option value="easy">Easy</option>
                  <option value="medium">Medium</option>
                  <option value="hard">Hard</option>
                </select>
              </div>
              <input
                type="number"
                placeholder="Points reward"
                value={newQuest.points_award}
                onChange={(e) => setNewQuest({...newQuest, points_award: parseInt(e.target.value)})}
                className="w-full p-2 bg-slate-700 rounded border border-slate-600"
                min="1"
                required
              />
              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="px-4 py-2 text-slate-400 hover:text-slate-200"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded"
                >
                  Create Quest
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Quest List */}
      <div className="space-y-4">
        {quests.length === 0 ? (
          <div className="text-center py-12">
            <Target className="w-16 h-16 text-slate-600 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-slate-300 mb-2">No quests yet</h3>
            <p className="text-slate-400 mb-4">Create your first quest to get started!</p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg"
            >
              Create Quest
            </button>
          </div>
        ) : (
          quests.map((quest) => (
            <QuestRow
              key={quest.id}
              quest={quest}
              onComplete={() => handleCompleteQuest(quest)}
            />
          ))
        )}
      </div>

      {/* Completion Dialog */}
      {completionQuest && (
        <CompletionDialog
          quest={completionQuest}
          onComplete={handleQuestCompleted}
          onCancel={() => setCompletionQuest(null)}
        />
      )}
    </div>
  );
};

export default Quests;
```

**Testing**:
- Test quest creation form
- Test quest listing from API
- Test quest completion flow
- Test dungeon switching
- Test error handling

### 3.6 Add Environment Configuration
**File**: `frontend/.env.example` (modify existing)

```env
# API Configuration
VITE_API_URL=http://localhost:8080

# Telegram Bot Configuration
VITE_BOT_USERNAME=YourBotUsername

# Development settings
VITE_DEV_MODE=true
```

**File**: `frontend/vite.config.ts` (modify existing)

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})
```

**Testing**:
- Test environment variables are loaded
- Test API proxy works in development
- Test production build works

## Testing Strategy

### 3.7 Frontend Integration Tests
**File**: `frontend/src/tests/integration.test.tsx` (new)

```typescript
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import App from '../App';

// Mock API service
jest.mock('../services/api');

describe('Frontend Integration', () => {
  test('shows login page when not authenticated', async () => {
    render(
      <BrowserRouter>
        <App />
      </BrowserRouter>
    );
    
    await waitFor(() => {
      expect(screen.getByText('Login with Telegram')).toBeInTheDocument();
    });
  });
  
  test('shows dashboard when authenticated', async () => {
    // Mock authenticated user
    // Test dashboard renders
  });
  
  test('quest creation flow works', async () => {
    // Test complete quest creation flow
  });
});
```

### 3.8 API Integration Tests
**File**: `test/integration/web_api_test.go`

```go
func TestWebAPIIntegration(t *testing.T) {
    // Test all web API endpoints
    // Test authentication flow
    // Test CORS headers
    // Test error responses
}
```

## Verification Checklist

- [ ] Frontend authenticates users via Telegram
- [ ] All pages connect to real backend APIs
- [ ] Quest creation, editing, and completion work
- [ ] Dungeon management works correctly
- [ ] Error handling provides good user experience
- [ ] Responsive design works on mobile
- [ ] Environment configuration works
- [ ] Integration tests pass

## Files to Create/Modify

1. `frontend/src/services/auth.ts` - Update with real auth
2. `frontend/src/services/api.ts` - Connect to real APIs
3. `frontend/src/App.tsx` - Add real authentication
4. `frontend/src/pages/Login.tsx` - New login page
5. `frontend/src/pages/Quests.tsx` - Connect to real data
6. `frontend/.env.example` - Environment configuration
7. `frontend/vite.config.ts` - Development proxy
8. `frontend/src/tests/integration.test.tsx` - Frontend tests
9. `test/integration/web_api_test.go` - Backend tests

## Success Criteria

✅ Users can authenticate via Telegram and access web interface
✅ All quest management features work through web UI
✅ Real-time data synchronization between bot and web
✅ Responsive design works on all devices
✅ Error handling provides clear feedback
✅ Performance is acceptable for typical usage
✅ Integration tests validate full user workflows

## Next Phase
Once web integration is complete, move to **Phase 4: Quest Reward System** to implement the point calculation and awarding logic.