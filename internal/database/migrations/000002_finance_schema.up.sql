-- SourcePayment (Ví / Ngân hàng)
CREATE TABLE IF NOT EXISTS sourcepayment (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'VND',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sourcepayment_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Categories (Danh mục Thu/Chi)
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT PRIMARY KEY,
    user_id BIGINT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    icon VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_category_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Transactions (Giao dịch)
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    sourcepayment_id BIGINT NOT NULL,
    category_id BIGINT,
    amount DECIMAL(15, 2) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transaction_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transaction_sourcepayment FOREIGN KEY(sourcepayment_id) REFERENCES sourcepayment(id) ON DELETE CASCADE,
    CONSTRAINT fk_transaction_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Budgets (Ngân sách)
CREATE TABLE IF NOT EXISTS budgets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_budget_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_budget_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
);
