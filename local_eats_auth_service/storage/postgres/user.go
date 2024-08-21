package postgres

import (
	pb "auth_serice/genproto"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"log"
	"time"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}

}

func (repo UserRepository) Register(request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = transaction.Rollback()
		} else if err != nil {
			_ = transaction.Rollback()
		} else {
			err = transaction.Commit()

		}
	}()

	id := uuid.New().String()
	query := `INSERT INTO users(id, user_name, email, password, full_name, user_type) VALUES ($1, $2, $3, $4, $5, $6)`
	result, err := transaction.Exec(query, id, request.UserName, request.Email, request.PasswordHash, request.FullName, request.UserType)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, fmt.Errorf("no rows were inserted")
	}

	response := pb.RegisterResponse{
		Id:           id,
		UserName:     request.UserName,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
		FullName:     request.FullName,
		UserType:     request.UserType,
		CreatedAt:    time.Now().String(),
	}
	return &response, nil
}

func (repo *UserRepository) Login(request *pb.LoginRequest) (*pb.Claims, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err := transaction.Rollback()
			if err != nil {
				return
			}
		}
	}()
	var claims pb.Claims
	transaction.QueryRow("select user_name,email,password ,user_type from users where  email=$1 and password=$2 and deleted_at is null", request.Email, request.PasswordHash).Scan(&claims.UserName, &claims.Email, &claims.PasswordHash, &claims.UserType)
	return &claims, err
}
func (repo *UserRepository) GetProfile(request *pb.IdRequest) (*pb.UserResponse, error) {
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

	var userProfile pb.UserResponse

	query := `
        SELECT id, user_name, email, full_name, user_type, address, phone_number, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

	var address sql.NullString // Use sql.NullString to handle NULL values

	err = transaction.QueryRow(query, request.Id).Scan(
		&userProfile.Id,
		&userProfile.UserName,
		&userProfile.Email,
		&userProfile.FullName,
		&userProfile.UserType,
		&address, // Scan into sql.NullString
		&userProfile.PhoneNumber,
		&userProfile.CreatedAt,
		&userProfile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if address.Valid {
		userProfile.Address = address.String
	} else {
		userProfile.Address = "" // or handle NULL case as needed
	}

	return &userProfile, nil
}

func (repo *UserRepository) EditProfile(request *pb.UpdateUserProfileRequest) (*pb.UserResponse, error) {
	transaction, err := repo.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := transaction.Rollback(); rollbackErr != nil {
				log.Printf("Rollback error: %v", rollbackErr)
			}
		} else {
			if commitErr := transaction.Commit(); commitErr != nil {
				log.Printf("Commit error: %v", commitErr)
			}
		}
	}()

	query := `UPDATE users SET full_name=$1, address=$2, phone_number=$3, updated_at=$4 WHERE id=$5 AND deleted_at IS NULL`
	_, err = transaction.Exec(query, request.FullName, request.Address, request.PhoneNumber, time.Now(), request.Id)
	if err != nil {
		return nil, err
	}

	var updatedUser pb.UserResponse
	query = `SELECT id, user_name, email, full_name, user_type, address, phone_number, created_at, updated_at FROM users WHERE id=$1 AND deleted_at IS NULL`
	err = transaction.QueryRow(query, request.Id).Scan(
		&updatedUser.Id, &updatedUser.UserName, &updatedUser.Email, &updatedUser.FullName,
		&updatedUser.UserType, &updatedUser.Address, &updatedUser.PhoneNumber,
		&updatedUser.CreatedAt, &updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

//func (repo *UserRepository) CreateRefreshToken() (*pb.Void, error) {
//	transaction, err := repo.Db.Begin()
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if err != nil {
//			err := transaction.Rollback()
//			if err != nil {
//				return
//			}
//		} else {
//			err := transaction.Commit()
//			if err != nil {
//				return
//			}
//		}
//	}()
//	result, err := transaction.Exec("update  users set token=$1 where email=$2")
//	if err != nil {
//		return nil, err
//	}
//	rowsAffected, err := result.RowsAffected()
//	if err != nil || rowsAffected == 0 {
//		return nil, err
//	}
//	return nil, nil
//}
