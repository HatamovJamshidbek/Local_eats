CREATE TABLE orders (
                        id UUID PRIMARY KEY,
                        user_id UUID,
                        kitchen_id UUID,
                        items JSONB NOT NULL,
                        total_amount DECIMAL(10, 2) NOT NULL,
                        status VARCHAR(20) NOT NULL,
                        delivery_address TEXT NOT NULL,
                        delivery_time TIMESTAMP ,
                        created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP
);


1->insert into  orders( user_id, kitchen_id, items, total_amount, status, delivery_address, delivery_time) values ($1,$2,$3,$4,$5,$6,$7)
2->update orders set orders.status=$1 and  orders.updated_at where id=$2 and deleted_at is null ;
3-> select  orders.id,orders.total_amount,orders.status,orders.delivery_time from orders where deleted_at is null and user_id=$1;
4-> select orders.id,orders.user_id,orders.total_amount,orders.status,orders.delivery_time from orders where  deleted_at  is null  and kitchen_id=$1;
5->select orders.user_id,orders.kitchen_id,orders.total_amount,orders.status,orders.delivery_address,orders.delivery_time,orders.items,orders.created_at,orders.updated_at,orders.deleted_at from orders where id=$1 and deleted_at is null;

