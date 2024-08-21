package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	pb "order_serive/genproto"
	"order_serive/help"
	"time"
)

type ReviewRepository struct {
	Db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{Db: db}
}

func (repo *ReviewRepository) CreateReview(request *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	query := `INSERT INTO reviews (id, order_id, user_id, kitchen_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	id := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)

	_, err := repo.Db.Exec(query, id, request.OrderId, request.UserId, request.KitchenId, request.Rating, request.Comment, createdAt)
	if err != nil {
		return nil, err
	}

	response := &pb.ReviewResponse{
		Id:        id,
		OrderId:   request.OrderId,
		UserId:    request.UserId,
		KitchenId: request.KitchenId,
		Rating:    request.Rating,
		Comment:   request.Comment,
		CreatedAt: createdAt,
	}

	return response, nil
}

func (repo *ReviewRepository) GetReviews(request *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	var (
		params = make(map[string]interface{})
		arr    []interface{}
		limit  string
		offset string
	)

	filter := " WHERE deleted_at IS NULL AND kitchen_id=$1"
	arr = append(arr, request.KitchenId)

	if request.LimitOffset.Limit > 0 {
		params["limit"] = request.LimitOffset.Limit
		limit = " LIMIT :limit"
	}
	if request.LimitOffset.Offset > 0 {
		params["offset"] = request.LimitOffset.Offset
		offset = " OFFSET :offset"
	}

	query := `SELECT id, user_id, order_id, rating, comment, created_at FROM reviews` + filter + limit + offset
	query, arr = help.ReplaceQueryParams(query, params)
	fmt.Println("Query:", query)

	rows, err := repo.Db.Query(query, arr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*pb.ReviewResponse
	var totalRating float64
	var count int

	for rows.Next() {
		var review pb.ReviewResponse
		err := rows.Scan(&review.Id, &review.UserId, &review.OrderId, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
		totalRating += review.Rating
		count++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := &pb.GetReviewsResponse{
		Reviews:       reviews,
		Total:         int32(count),
		AverageRating: totalRating / float64(count),
		LimitOffset: &pb.LimitOffset{
			Limit:  request.LimitOffset.Limit,
			Offset: request.LimitOffset.Offset,
		},
	}

	return response, nil
}
