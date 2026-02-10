CREATE TABLE IF NOT EXISTS customers (
    customer_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(user_id) NULL,
    name VARCHAR(100), 
    email VARCHAR(150), 
    phone VARCHAR(20), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (user_id IS NULL AND name IS NOT NULL AND email IS NOT NULL AND phone IS NOT NULL)
        OR
        (user_id IS NOT NULL AND name IS NULL AND email IS NULL AND phone IS NULL)
    )
);