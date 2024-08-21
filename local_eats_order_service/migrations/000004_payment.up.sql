CREATE TABLE payments (
                          id UUID PRIMARY KEY,
                          order_id UUID REFERENCES orders(id),
                          amount DECIMAL(10, 2) NOT NULL,
                          status VARCHAR(20) NOT NULL,
                          payment_method VARCHAR(50) NOT NULL,
                          transaction_id VARCHAR(100),
                          created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);


insert into payments( order_id, amount, status, payment_method, transaction_id) values ($1,$2,$3,$4,$5)
