CREATE TABLE kitchens (
                          id UUID PRIMARY KEY,
                          owner_id UUID REFERENCES users(id),
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          cuisine_type VARCHAR(50),
                          address TEXT NOT NULL,
                          phone_number VARCHAR(20),
                          rating DECIMAL(3, 2) DEFAULT 0,
                          total_orders INTEGER DEFAULT 0,
                          created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
                          deleted_at TIMESTAMP
);
