-- Migration: Create Savings Goals and Contributions schema
CREATE TABLE savings_goals (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    target_amount NUMERIC(15, 2) NOT NULL,
    current_amount NUMERIC(15, 2) DEFAULT 0,
    target_date TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active', -- active, completed, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE goal_contributions (
    id BIGINT PRIMARY KEY,
    goal_id BIGINT NOT NULL REFERENCES savings_goals(id) ON DELETE CASCADE,
    fund_id BIGINT REFERENCES funds(id) ON DELETE SET NULL, -- Source fund
    amount NUMERIC(15, 2) NOT NULL,
    type VARCHAR(50) NOT NULL, -- deposit, withdrawal
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to update updated_at on savings_goals
CREATE TRIGGER update_savings_goals_updated_at
BEFORE UPDATE ON savings_goals
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
