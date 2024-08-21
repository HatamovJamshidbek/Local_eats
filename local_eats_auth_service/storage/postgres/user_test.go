package postgres

import (
	pb "auth_serice/genproto"
	"database/sql"
	"testing"
)

func MockDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=username dbname=your_database sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestUserRepository_Register(t *testing.T) {
	db, err := MockDB()
	if err != nil {
		t.Fatalf("failed to connect to mock database: %v", err)
	}
	defer db.Close()
	repo := UserRepository{Db: db} // Assuming UserRepository is your repository struct

	request := &pb.RegisterRequest{
		UserName:     "testuser",
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		UserType:     "customer",
	}

	response, err := repo.Register(request)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response, got nil")
	}

	if response.Id != "" && response.Id != request.Email {
		t.Errorf("Expected ID to be %s, got %s", request.Email, response.Id)
	}
}

func TestUserRepository_Login(t *testing.T) {
	db, err := MockDB()
	if err != nil {
		t.Fatalf("failed to connect to mock database: %v", err)
	}
	defer db.Close()

	repo := &UserRepository{Db: db}
	request := &pb.LoginRequest{
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword",
	}

	claims, err := repo.Login(request)
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if claims == nil {
		t.Error("Expected non-nil claims, got nil")
	}

	if claims.Email != request.Email {
		t.Errorf("Expected UserName to be %s, got %s", request.Email, claims.UserName)
	}

}

func TestUserRepository_GetProfile(t *testing.T) {
	db, err := MockDB()
	if err != nil {
		t.Fatalf("failed to connect to mock database: %v", err)
	}
	defer db.Close()

	repo := &UserRepository{Db: db}

	request := &pb.IdRequest{
		Id: "user_id_here",
	}

	userProfile, err := repo.GetProfile(request)
	if err != nil {
		t.Fatalf("GetProfile returned error: %v", err)
	}

	if userProfile == nil {
		t.Error("Expected non-nil userProfile, got nil")
	}

	if userProfile.Id != request.Id {
		t.Errorf("Expected ID to be %s, got %s", request.Id, userProfile.Id)
	}

}
func TestUserRepository_EditProfile(t *testing.T) {
	db, err := MockDB()
	if err != nil {
		t.Fatalf("failed to connect to mock database: %v", err)
	}
	defer db.Close()

	repo := &UserRepository{Db: db}

	request := &pb.UpdateUserProfileRequest{
		Id:          "user_id_here",
		FullName:    "Updated Full Name",
		Address:     "Updated Address",
		PhoneNumber: "Updated Phone Number",
	}

	updatedUser, err := repo.EditProfile(request)
	if err != nil {
		t.Fatalf("EditProfile returned error: %v", err)
	}

	if updatedUser == nil {
		t.Error("Expected non-nil updatedUser, got nil")
	}

	if updatedUser.Id != request.Id {
		t.Errorf("Expected ID to be %s, got %s", request.Id, updatedUser.Id)
	}

}
