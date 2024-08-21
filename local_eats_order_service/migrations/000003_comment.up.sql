CREATE TABLE reviews (
                         id UUID PRIMARY KEY default gen_random_uuid(),
                         order_id UUID REFERENCES orders(id),
                         user_id UUID ,
                         kitchen_id UUID ,
                         rating DECIMAL(2, 1) NOT NULL,
                         comment TEXT,
                         created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
                         deleted_at TIMESTAMP
);

insert into reviews( order_id, user_id, kitchen_id, rating, comment) values ($1,$2,$3,$4,$5);
select  reviews.id,reviews.user_id,reviews.rating,reviews.comment from reviews where kitchen_id=$1 and deleted_at is null ;
