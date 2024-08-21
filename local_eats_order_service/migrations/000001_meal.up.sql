CREATE TABLE meals (
                        id UUID PRIMARY KEY default  gen_random_uuid(),
                        kitchen_id UUID ,
                        name VARCHAR(100) NOT NULL,
                        description TEXT,
                        price DECIMAL(10, 2) NOT NULL,
                        category VARCHAR(50),
                        ingredients TEXT[],
                        allergens TEXT[],
                        nutrition_info JSONB,
                        dietary_info TEXT[],
                        available BOOLEAN DEFAULT true,
                        created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP

);
--  1->insert into   meals(kitchen_id, name, description, price, category, ingredients, allergens, nutrition_info, dietary_info, available) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- 2->update meals
-- set  = name=$1,price=$2,available=$3,
-- where id=$4;

-- 3->update  meals set  deleted_at=current_time where id=$1

-- 4->select  meals.id,meals.name,meals.price,meals.category,meals.available,meals.description,meals.ingredients,meals.allergens,meals.dietary_info from meals

