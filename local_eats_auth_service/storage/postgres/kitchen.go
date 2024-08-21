package postgres

import (
	pb "auth_serice/genproto"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type KitchenRepository struct {
	Db *sql.DB
}

func NewKitchenRepository(db *sql.DB) *KitchenRepository {
	return &KitchenRepository{Db: db}
}

func (repo *KitchenRepository) CreateKitchen(request *pb.CreateKitchenRequest) (*pb.CreateKitchenResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	query1 := `SELECT id FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT 1`
	var UserId string
	err = transaction.QueryRow(query1).Scan(&UserId)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	query := `
		INSERT INTO kitchens (id, owner_id, name, description, cuisine_type, address, phone_number, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	result, err := transaction.Exec(query, id, UserId, request.Name, request.Description, request.CuisineType, request.Address, request.PhoneNumber, 0)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	response := &pb.CreateKitchenResponse{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
		OwnerId:     UserId,
		CuisineType: request.CuisineType,
		Address:     request.Address,
		PhoneNumber: request.PhoneNumber,
		Rating:      0,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	return response, nil
}

func (repo *KitchenRepository) UpdateKitchen(request *pb.UpdateKitchenRequest) (*pb.UpdateKitchenResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	query := `UPDATE kitchens SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`
	result, err := transaction.Exec(query, request.Name, request.Description, request.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, err
	}

	var response pb.UpdateKitchenResponse
	query = `SELECT id, owner_id, name, description, cuisine_type, address, phone_number, rating, updated_at FROM kitchens WHERE id=$1`
	err = transaction.QueryRow(query, request.Id).Scan(&response.Id, &response.OwnerId, &response.Name, &response.Description, &response.CuisineType, &response.Address, &response.PhoneNumber, &response.Rating, &response.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *KitchenRepository) GetAllKitchens(request *pb.LimitOffset) (*pb.KitchensResponse, error) {
	query := `
		SELECT id, owner_id, name, description, cuisine_type, address, phone_number, rating, total_orders, created_at, updated_at
		FROM kitchens
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var params []interface{}

	if request.Limit > 0 {
		query += " LIMIT $1"
		params = append(params, request.Limit)
	}
	if request.Offset > 0 {
		query += " OFFSET $2"
		params = append(params, request.Offset)
	}

	fmt.Println("Query:", query)
	fmt.Println("Params:", params)

	rows, err := repo.Db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query kitchens: %w", err)
	}
	defer rows.Close()

	var kitchens []*pb.KitchenResponse
	for rows.Next() {
		var kitchen pb.KitchenResponse
		err := rows.Scan(
			&kitchen.Id,
			&kitchen.OwnerId,
			&kitchen.Name,
			&kitchen.Description,
			&kitchen.CuisineType,
			&kitchen.Address,
			&kitchen.PhoneNumber,
			&kitchen.Rating,
			&kitchen.TotalOrder,
			&kitchen.CreatedAt,
			&kitchen.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan kitchen row: %w", err)
		}
		kitchens = append(kitchens, &kitchen)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &pb.KitchensResponse{
		Kitchens: kitchens,
		Total:    float32(len(kitchens)),
	}, nil
}

func (repo *KitchenRepository) SearchKitchen(request *pb.SearchKitchenRequest) (*pb.KitchensResponse, error) {
	// Check if repo.Db is initialized
	if repo.Db == nil {
		return nil, errors.New("database connection is nil")
	}

	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	var params []interface{}
	query := `
		SELECT id, owner_id, name, description, cuisine_type, address, phone_number, rating, total_orders, created_at, updated_at
		FROM kitchens
		WHERE deleted_at IS NULL
	`

	// Build WHERE conditions based on request filters
	if len(request.Name) > 0 {
		params = append(params, request.Name)
		query += " AND name = $1"
	}
	if request.Rating > 0 {
		params = append(params, request.Rating)
		query += " AND rating = $2"
	}
	if request.LimitOffset.Limit > 0 {
		params = append(params, request.LimitOffset.Limit)
		query += " LIMIT $9"
	}
	if request.LimitOffset.Offset > 0 {
		params = append(params, request.LimitOffset.Offset)
		query += " OFFSET $10"
	}

	fmt.Println("Query:", query)
	fmt.Println("Params:", params)

	rows, err := transaction.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query kitchens: %w", err)
	}
	defer rows.Close()

	var kitchens []*pb.KitchenResponse
	var count int
	for rows.Next() {
		var kitchen pb.KitchenResponse
		err := rows.Scan(
			&kitchen.Id, &kitchen.OwnerId, &kitchen.Name, &kitchen.Description,
			&kitchen.CuisineType, &kitchen.Address, &kitchen.PhoneNumber, &kitchen.Rating,
			&kitchen.TotalOrder, &kitchen.CreatedAt, &kitchen.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan kitchen row: %w", err)
		}
		kitchens = append(kitchens, &kitchen)
		count++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &pb.KitchensResponse{
		Kitchens: kitchens,
		Total:    float32(count),
	}, nil
}

func (repo *KitchenRepository) UpdateWorkingHours(req *pb.UpdateWorkingHoursRequest) (*pb.UpdateWorkingHoursResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			transaction.Rollback()
			panic(p)
		} else if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	daysOfWeek := map[string]int{
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
		"sunday":    7,
	}

	workingHours := map[string]*pb.TimeRange{
		"monday":    req.WorkingHours.Monday,
		"tuesday":   req.WorkingHours.Tuesday,
		"wednesday": req.WorkingHours.Wednesday,
		"thursday":  req.WorkingHours.Thursday,
		"friday":    req.WorkingHours.Friday,
		"saturday":  req.WorkingHours.Saturday,
		"sunday":    req.WorkingHours.Sunday,
	}
	query := `
        INSERT INTO working_hours (kitchen_id, day_of_week, open_time, close_time)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (kitchen_id, day_of_week) DO UPDATE
        SET open_time = EXCLUDED.open_time,
            close_time = EXCLUDED.close_time
    `

	for day, timeRange := range workingHours {
		if timeRange != nil {
			_, err = transaction.Exec(query, req.KitchenId, daysOfWeek[day], timeRange.Open, timeRange.Close)
			if err != nil {
				return nil, err
			}
		}
	}

	res := &pb.UpdateWorkingHoursResponse{
		KitchenId:    req.KitchenId,
		WorkingHours: req.WorkingHours,
		UpdatedAt:    time.Now().String(),
	}

	return res, nil
}

func (repo *KitchenRepository) GetByIdKitchen(request *pb.IdRequest) (*pb.KitchenResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		fmt.Println("+++++", err)
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	query := `
        SELECT id, owner_id, name, description, cuisine_type, address, phone_number, rating, total_orders, created_at, updated_at
        FROM kitchens
        WHERE id = $1 AND deleted_at IS NULL
    `

	var kitchen pb.KitchenResponse
	err = repo.Db.QueryRow(query, id).Scan(
		&kitchen.Id,
		&kitchen.OwnerId,
		&kitchen.Name,
		&kitchen.Description,
		&kitchen.CuisineType,
		&kitchen.Address,
		&kitchen.PhoneNumber,
		&kitchen.Rating,
		&kitchen.TotalOrder,
		&kitchen.CreatedAt,
		&kitchen.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kitchen with ID %s not found", request.Id)
		}
		return nil, fmt.Errorf("error querying kitchen: %v", err)
	}

	return &kitchen, nil
}
