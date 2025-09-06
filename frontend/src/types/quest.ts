export interface Quest {
  id: string;
  title: string;
  description: string;
  category: 'daily' | 'weekly' | 'adhoc';
  difficulty: 'easy' | 'medium' | 'hard';
  mode: 'BINARY' | 'PARTIAL' | 'PER_MINUTE';
  points_award: string;
  rate_points_per_min?: string;
  min_minutes?: number;
  max_minutes?: number;
  daily_points_cap?: string;
  cooldown_sec: number;
  streak_enabled: boolean;
  status: 'active' | 'paused' | 'archived';
}

export interface CreateQuestRequest {
  title: string;
  description: string;
  category: 'daily' | 'weekly' | 'adhoc';
  difficulty: 'easy' | 'medium' | 'hard';
  mode: 'BINARY' | 'PARTIAL' | 'PER_MINUTE';
  points_award: string;
  rate_points_per_min?: string;
  min_minutes?: number;
  max_minutes?: number;
  daily_points_cap?: string;
  cooldown_sec?: number;
  streak_enabled?: boolean;
  status?: 'active' | 'paused' | 'archived';
}

export interface CompleteQuestRequest {
  idempotency_key: string;
  completion_ratio?: number;
  minutes?: number;
}

export interface CompleteQuestResponse {
  awarded_points: string;
  submitted_at: string;
  streak_count?: number;
}