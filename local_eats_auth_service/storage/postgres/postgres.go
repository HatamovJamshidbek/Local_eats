package postgres

import (
	"auth_serice/config"
	"auth_serice/storage"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db          *sql.DB
	userRepo    storage.IUserStorage
	kitchenRepo storage.IKitchenStorage
}

func (s *PostgresStorage) Users() storage.IUserStorage {
	return s.userRepo
}

func (s *PostgresStorage) Kitchens() storage.IKitchenStorage {
	return s.kitchenRepo
}

func NewPostgresStorage() (storage.IStorage, error) {
	cnf := config.Load()

	conDb := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cnf.PostgresHost, cnf.PostgresPort, cnf.PostgresUser, cnf.PostgresDB, cnf.PostgresPassword)

	db, err := sql.Open("postgres", conDb)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	userRepo := NewUserRepository(db)
	kitchenRepo := NewKitchenRepository(db)

	return &PostgresStorage{
		db:          db,
		userRepo:    userRepo,
		kitchenRepo: kitchenRepo,
	}, nil
}

func (s *PostgresStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("error closing PostgreSQL connection: %w", err)
	}
	return nil
}
