-- 007_create_budgets.sql

CREATE TABLE IF NOT EXISTS budgets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL, -- NULL = Total budget
    amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    period VARCHAR(20) NOT NULL, -- 'monthly', 'weekly', 'custom'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Ensure no overlapping budgets for the same category/period (simplified)
    CONSTRAINT uk_user_category_period UNIQUE (user_id, category_id, period, start_date)
);

-- Index for fast lookup of active budgets for a user
CREATE INDEX IF NOT EXISTS idx_budgets_user_active_lookup ON budgets(user_id, is_active, start_date, end_date);
