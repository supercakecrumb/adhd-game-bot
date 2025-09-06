export interface Dungeon {
  id: string;
  title: string;
  admin_user_id: number;
  telegram_chat_id?: number;
  created_at: string;
}

export interface CreateDungeonRequest {
  title: string;
  telegram_chat_id?: number;
}

export interface DungeonSummary {
  balance: string;
  todays_quests: any[];
  recent_purchases: any[];
}