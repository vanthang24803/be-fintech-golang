export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

export interface ApiError {
  code: number;
  message: string;
  error?: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  google_id: string | null;
  created_at: string;
  updated_at: string;
}

export interface Profile {
  id: string;
  user_id: string;
  full_name: string | null;
  avatar_url: string | null;
  phone_number: string | null;
  date_of_birth: string | null;
  created_at: string;
  updated_at: string;
}

export interface GetProfileResponse {
  user: User;
  profile: Profile;
}

export interface Category {
  id: string;
  user_id: string;
  name: string;
  type: "income" | "expense";
  icon: string | null;
  created_at: string;
  updated_at: string;
}

export interface SourcePayment {
  id: string;
  user_id: string;
  name: string;
  type: string;
  balance: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: string;
  user_id: string;
  source_payment_id: string;
  category_id: string;
  amount: number;
  type: "income" | "expense";
  description: string | null;
  transaction_date: string;
  source_name: string;
  category_name: string;
  created_at: string;
  updated_at: string;
}

export interface MonthlyTrendItem {
  month: string;
  income: number;
  expense: number;
  net_profit: number;
}

export interface CategorySummaryItem {
  category_id: string;
  category_name: string;
  category_icon: string | null;
  total_amount: number;
  percentage: number;
}
