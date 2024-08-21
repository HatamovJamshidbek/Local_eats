package postgres

import (
	"database/sql"
	"fmt"
	"order_serive/config"
	st "order_serive/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db       *sql.DB
	Comments st.Reviews
	Meals    st.Meals
	Orders   st.Orders
	Payments st.Payments
}

func NewPostgresStorage() (st.InitRoot, error) {
	cnf := config.Load()
	conDb := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", cnf.PostgresHost, cnf.PostgresPort, cnf.PostgresUser, cnf.PostgresDatabase, cnf.PostgresPassword)
	db, err := sql.Open("postgres", conDb)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Storage{
		Db: db,
	}, nil
}

func (s *Storage) Review() st.Reviews {
	if s.Comments == nil {
		s.Comments = &ReviewRepository{s.Db}
	}
	return s.Comments
}

func (s *Storage) Meal() st.Meals {
	if s.Meals == nil {
		s.Meals = &MealRepository{s.Db}
	}
	return s.Meals
}

func (s *Storage) Order() st.Orders {
	if s.Orders == nil {
		s.Orders = &OrderRepository{s.Db}
	}
	return s.Orders
}

func (s *Storage) Payment() st.Payments {
	if s.Payments == nil {
		s.Payments = &PaymentRepository{s.Db}
	}
	return s.Payments
}
