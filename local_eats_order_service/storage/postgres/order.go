package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	pb "order_serive/genproto"
	"order_serive/help"
	"time"

	"github.com/google/uuid"
)

type OrderRepository struct {
	Db *sql.DB
}

func (s *OrderRepository) UpdateOrderStatus(request *pb.UpdateOrderStatusRequest) (*pb.OrderStatusResponse, error) {
	transaction, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id, status, updated_at`
	var updatedOrder pb.OrderStatusResponse
	err = transaction.QueryRow(query, request.Status, time.Now(), request.Id).Scan(
		&updatedOrder.Id,
		&updatedOrder.Status,
		&updatedOrder.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (repo *OrderRepository) CreateOrder(request *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback()
			panic(p)
		} else if err != nil {
			_ = transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()
	itemsJSON, err := json.Marshal(request.Items)
	if err != nil {
		return nil, err
	}

	requestId := uuid.New().String()

	deliveryTime, err := time.Parse(time.RFC3339, request.DeliveryTime)
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO orders (
            id, user_id, kitchen_id, items, total_amount, status,
            delivery_address, delivery_time, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err = transaction.Exec(query, requestId,
		request.UserId, request.KitchenId, itemsJSON,
		request.TotalAmount, request.Status, request.DeliveryAddress,
		deliveryTime, time.Now())
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	// Prepare and return the response
	response := &pb.OrderResponse{
		Id:              requestId,
		UserId:          request.UserId,
		KitchenId:       request.KitchenId,
		Items:           request.Items,
		TotalAmount:     request.TotalAmount,
		Status:          request.Status,
		DeliveryAddress: request.DeliveryAddress,
		DeliveryTime:    request.DeliveryTime,
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	return response, nil
}
func (repo *OrderRepository) GetOrdersForCustomer(request *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	var (
		params = make(map[string]interface{})
		arr    []interface{}
		limit  string
		offset string
	)
	filter := " WHERE deleted_at IS NULL AND user_id=$1"
	arr = append(arr, request.UserId)

	if request.Status != "" {
		params["status"] = request.Status
		filter += " AND status = :status"
	}
	if request.LimitOffset.Limit > 0 {
		params["limit"] = request.LimitOffset.Limit
		limit = " LIMIT :limit"
	}
	if request.LimitOffset.Offset > 0 {
		params["offset"] = request.LimitOffset.Offset
		offset = " OFFSET :offset"
	}

	query := `SELECT id, kitchen_id, total_amount, status, delivery_time FROM orders` + filter + limit + offset
	query, arr = help.ReplaceQueryParams(query, params)
	fmt.Println("Query:", query)

	rows, err := repo.Db.Query(query, arr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*pb.OrderResponse
	for rows.Next() {
		var order pb.OrderResponse
		err := rows.Scan(&order.Id, &order.KitchenId, &order.TotalAmount, &order.Status, &order.DeliveryTime)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	response := &pb.GetOrdersResponse{
		Orders: orders,
		Total:  int32(len(orders)),
		Page:   request.LimitOffset.Offset,
		Limit:  request.LimitOffset.Limit,
	}

	return response, nil
}

func (repo *OrderRepository) GetOrdersForChef(request *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	var (
		params = make(map[string]interface{})
		arr    []interface{}
		limit  string
		offset string
	)
	filter := " WHERE deleted_at IS NULL AND kitchen_id=$1"
	arr = append(arr, request.KitchenId)

	if request.Status != "" {
		params["status"] = request.Status
		filter += " AND status = :status"
	}
	if request.LimitOffset.Limit > 0 {
		params["limit"] = request.LimitOffset.Limit
		limit = " LIMIT :limit"
	}
	if request.LimitOffset.Offset > 0 {
		params["offset"] = request.LimitOffset.Offset
		offset = " OFFSET :offset"
	}

	query := `SELECT id, user_id, total_amount, status, delivery_time FROM orders` + filter + limit + offset
	query, arr = help.ReplaceQueryParams(query, params)
	fmt.Println("Query:", query)

	rows, err := repo.Db.Query(query, arr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*pb.OrderResponse
	for rows.Next() {
		var order pb.OrderResponse
		err := rows.Scan(&order.Id, &order.UserId, &order.TotalAmount, &order.Status, &order.DeliveryTime)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	response := &pb.GetOrdersResponse{
		Orders: orders,
		Total:  int32(len(orders)),
		Page:   request.LimitOffset.Offset,
		Limit:  request.LimitOffset.Limit,
	}

	return response, nil
}

func (repo *OrderRepository) GetOrderById(request *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	query := `SELECT id, user_id, kitchen_id, items, total_amount, status, delivery_address, delivery_time, created_at, updated_at 
              FROM orders WHERE id=$1`
	var (
		itemsJSON []byte
		response  pb.OrderResponse
	)

	err := repo.Db.QueryRow(query, request.Id).Scan(
		&response.Id,
		&response.UserId,
		&response.KitchenId,
		&itemsJSON,
		&response.TotalAmount,
		&response.Status,
		&response.DeliveryAddress,
		&response.DeliveryTime,
		&response.CreatedAt,
		&response.UpdatedAt,
	)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	if err := json.Unmarshal(itemsJSON, &response.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %v", err)
	}

	return &response, nil
}

func (repo *OrderRepository) GetUserActivity(request *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error) {
	queryTotalOrders := `
        SELECT COUNT(*), COALESCE(SUM(total_amount), 0)
        FROM orders
        WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
    `
	queryFavoriteCuisines := `
        SELECT COALESCE(items->>'cuisine_type', 'Unknown') AS cuisine_type, COUNT(*) AS orders_count
        FROM orders
        WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
        GROUP BY items->>'cuisine_type'
        ORDER BY orders_count DESC
    `
	queryFavoriteKitchens := `
        SELECT kitchen_id, COUNT(*) AS orders_count
        FROM orders
        WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
        GROUP BY kitchen_id
        ORDER BY orders_count DESC
    `

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

	var totalOrders int32
	var totalSpent float64
	err = transaction.QueryRow(queryTotalOrders, request.UserId, request.StartDate, request.EndDate).Scan(&totalOrders, &totalSpent)
	if err != nil {
		return nil, err
	}

	// Fetch favorite cuisines
	rows, err := transaction.Query(queryFavoriteCuisines, request.UserId, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favoriteCuisines []*pb.CuisineActivity
	for rows.Next() {
		var cuisine pb.CuisineActivity
		if err := rows.Scan(&cuisine.CuisineType, &cuisine.OrdersCount); err != nil {
			return nil, err
		}
		favoriteCuisines = append(favoriteCuisines, &cuisine)
	}

	// Fetch favorite kitchens
	rows, err = transaction.Query(queryFavoriteKitchens, request.UserId, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favoriteKitchens []*pb.KitchenActivity
	for rows.Next() {
		var kitchen pb.KitchenActivity
		var kitchenID string
		if err := rows.Scan(&kitchenID, &kitchen.OrdersCount); err != nil {
			return nil, err
		}

		// Fetch kitchen name from database based on kitchenID if needed
		kitchen.Id = kitchenID
		kitchen.Name = "Sample Kitchen Name" // Replace with actual database lookup

		favoriteKitchens = append(favoriteKitchens, &kitchen)
	}

	response := &pb.GetUserActivityResponse{
		TotalOrders:      totalOrders,
		TotalSpent:       float32(totalSpent),
		FavoriteCuisines: favoriteCuisines,
		FavoriteKitchens: favoriteKitchens,
	}

	return response, nil
}
func (repo *OrderRepository) GetKitchenStatistics(request *pb.GetKitchenStatisticRequest) (*pb.KitchenStatisticsResponse, error) {
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

	queryTotalOrders := `
        SELECT COUNT(*), COALESCE(SUM(total_amount), 0)
        FROM orders
        WHERE kitchen_id = $1 AND created_at BETWEEN $2 AND $3
    `
	queryAverageRating := `
        SELECT COALESCE(AVG(rating), 0)
        FROM reviews
        WHERE kitchen_id = $1 AND created_at BETWEEN $2 AND $3
    `
	queryTopDishes := `
        SELECT d.id, d.name, COUNT(*) AS orders_count, COALESCE(SUM(o.total_amount), 0) AS revenue
        FROM meals d
        JOIN orders o ON o.items @> jsonb_build_array(jsonb_build_object('id', d.id))::jsonb
        WHERE o.kitchen_id = $1 AND o.created_at BETWEEN $2 AND $3
        GROUP BY d.id, d.name
        ORDER BY orders_count DESC
        LIMIT 5
    `
	queryBusiestHours := `
        SELECT EXTRACT(HOUR FROM created_at) AS hour, COUNT(*) AS orders_count
        FROM orders
        WHERE kitchen_id = $1 AND created_at BETWEEN $2 AND $3
        GROUP BY EXTRACT(HOUR FROM created_at)
        ORDER BY orders_count DESC
        LIMIT 5
    `

	var totalOrders int32
	var totalRevenue float64
	var averageRating float32
	var topDishes []*pb.TopDish
	var busiestHours []*pb.BusiestHour

	err = transaction.QueryRow(queryTotalOrders, request.UserId, request.StartDate, request.EndDate).Scan(&totalOrders, &totalRevenue)
	if err != nil {
		return nil, err
	}

	err = transaction.QueryRow(queryAverageRating, request.UserId, request.StartDate, request.EndDate).Scan(&averageRating)
	if err != nil {
		return nil, err
	}

	rows, err := transaction.Query(queryTopDishes, request.UserId, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dish pb.TopDish
		if err := rows.Scan(&dish.Id, &dish.Name, &dish.OrdersCount, &dish.Revenue); err != nil {
			return nil, err
		}
		topDishes = append(topDishes, &dish)
	}

	rows, err = transaction.Query(queryBusiestHours, request.UserId, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hour int32
		var ordersCount int32
		if err := rows.Scan(&hour, &ordersCount); err != nil {
			return nil, err
		}
		busiestHours = append(busiestHours, &pb.BusiestHour{Hour: hour, OrdersCount: ordersCount})
	}

	response := &pb.KitchenStatisticsResponse{
		TotalOrders:   totalOrders,
		TotalRevenue:  float32(totalRevenue),
		AverageRating: averageRating,
		TopDishes:     topDishes,
		BusiestHours:  busiestHours,
	}

	return response, nil
}
