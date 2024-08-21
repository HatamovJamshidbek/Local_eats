package postgres

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	pb "order_serive/genproto"
)

type PaymentRepository struct {
	Db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{Db: db}
}

func (repo *PaymentRepository) CreatePayment(request *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	paymentId := uuid.New().String()
	transactionId := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)

	query := `
	INSERT INTO payments (id, order_id, amount, payment_method, status, transaction_id, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	amount := 33.97
	status := "success"

	_, err = transaction.Exec(query, paymentId, request.OrderId, amount, request.PaymentMethod, status, transactionId, createdAt)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	response := &pb.CreatePaymentResponse{
		Id:            paymentId,
		OrderId:       request.OrderId,
		Amount:        amount,
		Status:        status,
		TransactionId: transactionId,
		CreatedAt:     createdAt,
	}

	return response, nil
}
