BEGIN;
CREATE TABLE expense_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(120) NOT NULL CHECK (length(trim(name)) > 0),
    description VARCHAR(1000) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by TEXT NOT NULL REFERENCES users(id),
    updated_by TEXT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX expense_categories_name_unique ON expense_categories (lower(trim(name)));
CREATE SEQUENCE expense_reference_seq;
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference TEXT NOT NULL UNIQUE DEFAULT ('EXP-' || nextval('expense_reference_seq')::text),
    expense_date DATE NOT NULL,
    category_id UUID NOT NULL REFERENCES expense_categories(id) ON DELETE RESTRICT,
    category_name VARCHAR(120) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(5000) NOT NULL,
    amount NUMERIC(14,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'UGX' CHECK (currency = 'UGX'),
    payee VARCHAR(200) NOT NULL,
    payment_method VARCHAR(32) NOT NULL CHECK (payment_method IN ('CASH','BANK_TRANSFER','MOBILE_MONEY','CARD','OTHER')),
    payment_reference VARCHAR(200) NOT NULL DEFAULT '',
    department_id TEXT NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    department_name VARCHAR(120) NOT NULL,
    recorded_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    recorded_by_name TEXT NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX expenses_date_idx ON expenses (expense_date DESC, recorded_at DESC);
CREATE INDEX expenses_category_idx ON expenses(category_id);
CREATE INDEX expenses_department_idx ON expenses(department_id);
-- Small supporting receipts are stored atomically with the record, never as public media URLs.
CREATE TABLE expense_attachments (
    expense_id UUID PRIMARY KEY REFERENCES expenses(id) ON DELETE RESTRICT,
    original_name VARCHAR(255) NOT NULL,
    mime_type TEXT NOT NULL CHECK (mime_type IN ('application/pdf','image/jpeg','image/png','image/webp')),
    size_bytes INTEGER NOT NULL CHECK (size_bytes > 0 AND size_bytes <= 5242880),
    content BYTEA NOT NULL,
    CHECK (octet_length(content) = size_bytes)
);
COMMIT;
