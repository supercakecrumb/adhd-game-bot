export interface ShopItem {
  id: number;
  chat_id: number;
  dungeon_id?: string;
  code: string;
  name: string;
  description: string;
  price: string;
  category: string;
  is_active: boolean;
  stock?: number;
  discount_tier_id?: number;
  created_at: string;
  updated_at: string;
}

export interface Purchase {
  id: number;
  user_id: number;
  item_id: number;
  dungeon_id: string;
  item_name: string;
  item_price: string;
  quantity: number;
  total_cost: string;
  status: 'pending' | 'completed' | 'refunded';
  discount_tier_id?: number;
  purchased_at: string;
}

export interface CreateShopItemRequest {
  name: string;
  code: string;
  description: string;
  price: string;
  category: string;
  is_active: boolean;
  stock?: number;
  discount_tier_id?: number;
}

export interface PurchaseRequest {
  quantity: number;
  idempotency_key: string;
}