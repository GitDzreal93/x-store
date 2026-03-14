// ==================== 通用类型 ====================

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

export interface PageResult<T = any> {
  list: T[];
  total: number;
  page: number;
  size: number;
}

// ==================== 用户/认证 ====================

export interface User {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  role: string;
  status: number;
  created_at: string;
  updated_at: string;
}

export interface LoginReq {
  username: string;
  password: string;
}

export interface LoginResp {
  token: string;
  user: User;
}

// ==================== 分类 ====================

export interface Category {
  id: number;
  name: string;
  icon: string;
  sort: number;
  parent_id: number | null;
  status: number;
  created_at: string;
  updated_at: string;
}

// ==================== 商品 ====================

export interface ProductDetail {
  id: number;
  product_id: number;
  description: string;
  notice: string;
}

export interface Product {
  id: number;
  category_id: number;
  title: string;
  cover: string;
  price: number;
  original_price: number;
  stock: number;
  sales: number;
  delivery_type: string;
  tags: string;
  sort: number;
  status: number;
  is_new: boolean;
  is_hot: boolean;
  created_at: string;
  updated_at: string;
  category?: Category;
  detail?: ProductDetail;
}

export interface CreateProductReq {
  category_id: number;
  title: string;
  cover?: string;
  price: number;
  original_price?: number;
  delivery_type?: string;
  tags?: string;
  sort?: number;
  is_new?: boolean;
  is_hot?: boolean;
  description?: string;
  notice?: string;
}

export interface UpdateProductReq {
  category_id?: number;
  title?: string;
  cover?: string;
  price?: number;
  original_price?: number;
  delivery_type?: string;
  tags?: string;
  sort?: number;
  status?: number;
  is_new?: boolean;
  is_hot?: boolean;
  description?: string;
  notice?: string;
}

// ==================== 订单 ====================

export interface Order {
  id: number;
  order_no: string;
  user_id: number | null;
  product_id: number;
  email: string;
  amount: number;
  status: number;
  pay_method: string;
  paid_at: string | null;
  expire_at: string;
  created_at: string;
  updated_at: string;
  product?: Product;
  card_keys?: CardKey[];
}

// 订单状态
export const ORDER_STATUS_MAP: Record<number, { label: string; color: string }> = {
  0: { label: '待支付', color: 'gold' },
  1: { label: '已支付', color: 'blue' },
  2: { label: '已发货', color: 'green' },
  3: { label: '已完成', color: 'default' },
  4: { label: '已取消', color: 'red' },
  5: { label: '已退款', color: 'purple' },
};

// ==================== 卡密 ====================

export interface CardKey {
  id: number;
  product_id: number;
  order_id: number | null;
  content: string;
  status: number;
  created_at: string;
  updated_at: string;
}

export const CARD_KEY_STATUS_MAP: Record<number, { label: string; color: string }> = {
  0: { label: '可用', color: 'green' },
  1: { label: '已锁定', color: 'orange' },
  2: { label: '已售出', color: 'default' },
};

// ==================== 支付 ====================

export interface Payment {
  id: number;
  order_id: number;
  trade_no: string;
  pay_method: string;
  amount: number;
  status: number;
  raw_notify: string;
  completed_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface PaymentChannel {
  id: number;
  name: string;
  provider_type: string;
  channel_type: string;
  interaction_mode: string;
  config_json: string;
  fee_rate: number;
  fixed_fee: number;
  is_active: boolean;
  sort: number;
  created_at: string;
  updated_at: string;
}
