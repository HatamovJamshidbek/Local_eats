package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	pb "order_serive/genproto"
	"order_serive/help"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/spf13/cast"
)

type MealRepository struct {
	Db *sql.DB
}

func NewMealRepository(db *sql.DB) *MealRepository {
	return &MealRepository{Db: db}
}

func (repo *MealRepository) CreateMeal(request *pb.CreateMealRequest) (*pb.MealResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	request.Id = cast.ToString(uuid.New())
	query := `
        INSERT INTO meals (
            id, kitchen_id, name, description, price, category, available, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp)
    `
	_, err = transaction.Exec(query, request.Id, request.KitchenId, request.Name, request.Description, request.Price, request.Category, request.Available)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	response := pb.MealResponse{
		Id:          request.Id,
		KitchenId:   request.KitchenId,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Category:    request.Category,
		Available:   request.Available,
		CreatedAt:   cast.ToString(time.Now()),
	}
	return &response, nil
}

func (repo *MealRepository) UpdateMeal(request *pb.UpdateMealRequest) (*pb.MealResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	query := `
        UPDATE meals
        SET name = $1, price = $2, available = $3
        WHERE id = $4 AND deleted_at IS NULL
    `
	res, err := transaction.Exec(query, request.Name, request.Price, request.Available, request.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, err
	}

	var response pb.MealResponse
	selectQuery := `
        SELECT kitchen_id, name, description, price, category,
               ingredients, allergens, available, updated_at
        FROM meals
        WHERE id = $1 AND deleted_at IS NULL
    `
	err = transaction.QueryRow(selectQuery, request.Id).Scan(
		&response.KitchenId, &response.Name, &response.Description, &response.Price, &response.Category,
		pq.Array(&response.Ingredients), pq.Array(&response.Allergens), &response.Available, &response.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *MealRepository) DeleteMeal(request *pb.IdRequest) (*pb.Void, error) {
	result, err := repo.Db.Exec("UPDATE meals SET deleted_at = current_timestamp WHERE id = $1 AND deleted_at IS NULL", request.Id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (repo *MealRepository) GetAllMeal(request *pb.GetAllMealRequest) (*pb.MealsResponse, error) {
	var (
		params = make(map[string]interface{})
		arr    []interface{}
		limit  string
		offset string
	)
	if repo == nil || repo.Db == nil {
		fmt.Println("+++++++++++++++")
		return nil, fmt.Errorf("repository or database connection is nil")
	}
	filter := ""
	if len(request.Name) > 0 {
		params["name"] = request.Name
		filter += " AND name = :name"
	}
	if len(request.Category) > 0 {
		params["category"] = request.Category
		filter += " AND category = :category"
	}
	if request.Price > 0 {
		params["price"] = request.Price
		filter += " AND price = :price"
	}
	if len(request.KitchenId) > 0 {
		params["kitchen_id"] = request.KitchenId
		filter += " AND kitchen_id = :kitchen_id"
	}
	if request.Available {
		params["available"] = request.Available
		filter += " AND available = :available"
	}

	if request.LimitOffset.Limit > 0 {
		params["limit"] = request.LimitOffset.Limit
		limit = " LIMIT :limit"

	}
	if request.LimitOffset.Offset > 0 {
		params["offset"] = request.LimitOffset.Offset
		offset = " OFFSET :offset"
	}
	query := `SELECT id, name, price, category, available, description, ingredients, allergens, dietary_info, created_at, updated_at 
			  FROM meals WHERE deleted_at IS NULL ` + filter + limit + offset
	query, arr = help.ReplaceQueryParams(query, params)
	fmt.Println("Query:", query)

	rows, err := repo.Db.Query(query, arr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []*pb.MealResponse
	for rows.Next() {
		var meal pb.MealResponse
		err := rows.Scan(
			&meal.Id, &meal.Name, &meal.Price, &meal.Category,
			&meal.Available, &meal.Description, pq.Array(&meal.Ingredients),
			pq.Array(&meal.Allergens), pq.Array(&meal.DietaryInfo),
			&meal.CreatedAt, &meal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		meals = append(meals, &meal)
	}
	fmt.Println("++++++", meals)
	return &pb.MealsResponse{Meals: meals}, nil
}

func (repo *MealRepository) UpdateNutritionInfo(request *pb.UpdateNutritionInfoRequest) (*pb.Dish, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	query := `
        UPDATE meals 
        SET allergens = $1, 
            nutrition_info = jsonb_build_object('calories', CAST($2 AS NUMERIC), 'protein', CAST($3 AS NUMERIC), 'carbohydrates', CAST($4 AS NUMERIC), 'fat', CAST($5 AS NUMERIC)),
            dietary_info = $6, 
            updated_at = $7 
        WHERE id = $8 
        RETURNING id, name, allergens, nutrition_info, dietary_info, updated_at`

	updatedAt := time.Now().UTC()
	var dish pb.Dish

	// Log the types and values of the parameters
	fmt.Printf("Allergens: %v (type: %T)\n", request.Allergens, request.Allergens)
	fmt.Printf("Calories: %v (type: %T)\n", request.NutritionInfo.Calories, request.NutritionInfo.Calories)
	fmt.Printf("Protein: %v (type: %T)\n", request.NutritionInfo.Protein, request.NutritionInfo.Protein)
	fmt.Printf("Carbohydrates: %v (type: %T)\n", request.NutritionInfo.Carbohydrates, request.NutritionInfo.Carbohydrates)
	fmt.Printf("Fat: %v (type: %T)\n", request.NutritionInfo.Fat, request.NutritionInfo.Fat)
	fmt.Printf("DietaryInfo: %v (type: %T)\n", request.DietaryInfo, request.DietaryInfo)
	fmt.Printf("UpdatedAt: %v (type: %T)\n", updatedAt, updatedAt)
	fmt.Printf("DishId: %v (type: %T)\n", request.DishId, request.DishId)

	var nutritionInfoJSON []byte
	err = transaction.QueryRow(
		query,
		pq.Array(request.Allergens),
		request.NutritionInfo.Calories,
		request.NutritionInfo.Protein,
		request.NutritionInfo.Carbohydrates,
		request.NutritionInfo.Fat,
		pq.Array(request.DietaryInfo),
		updatedAt,
		request.DishId,
	).Scan(
		&dish.Id, &dish.Name, pq.Array(&dish.Allergens), &nutritionInfoJSON, pq.Array(&dish.DietaryInfo), &dish.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the nutritionInfoJSON into pb.Dish.NutritionInfo
	err = json.Unmarshal(nutritionInfoJSON, &dish.NutritionInfo)
	if err != nil {
		return nil, err
	}

	dish.UpdatedAt = updatedAt.String()

	return &dish, nil
}
