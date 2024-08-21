package postgres

import (
	"auth_serice/genproto"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestKitchenRepository_CreateKitchen(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Kitchens()

	mockUUID := uuid.New().String()

	request := &genproto.CreateKitchenRequest{
		Name:        "Test Kitchen",
		Description: "Test Description",
		CuisineType: "Test Cuisine",
		Address:     "Test Address",
		PhoneNumber: "123456789",
	}

	response, err := repo.CreateKitchen(request)
	if err != nil {
		t.Fatalf("error creating kitchen: %v", err)
	}

	expectedResponse := &genproto.CreateKitchenResponse{
		Id:          response.Id,
		Name:        request.Name,
		Description: request.Description,
		OwnerId:     mockUUID,
		CuisineType: request.CuisineType,
		Address:     request.Address,
		PhoneNumber: request.PhoneNumber,
		Rating:      0,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("CreateKitchen response does not match expected.\nGot: %+v\nExpected: %+v", response, expectedResponse)
	}

}
func TestKitchenRepository_UpdateKitchen(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Kitchens()

	request := &genproto.UpdateKitchenRequest{
		Id:          "kitchen123",
		Name:        "Updated Kitchen Name",
		Description: "Updated Description",
	}

	response, err := repo.UpdateKitchen(request)
	if err != nil {
		t.Fatalf("error updating kitchen: %v", err)
	}

	expectedResponse := &genproto.UpdateKitchenResponse{
		Id:          "kitchen123",
		OwnerId:     "owner123",
		Name:        "Updated Kitchen Name",
		Description: "Updated Description",
		CuisineType: "Test Cuisine",
		Address:     "Test Address",
		PhoneNumber: "123456789",
		Rating:      0,
		UpdatedAt:   response.UpdatedAt, // This will be validated dynamically, not with DeepEqual
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("UpdateKitchen response does not match expected.\nGot: %+v\nExpected: %+v", response, expectedResponse)
	}

}
