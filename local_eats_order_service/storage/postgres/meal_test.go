package postgres

import (
	"reflect"
	"testing"

	pb "order_serive/genproto"
)

func TestMealRepository_CreateMeal(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Meal()

	request := &pb.CreateMealRequest{
		KitchenId:   "kitchen123",
		Name:        "Spaghetti Carbonara",
		Description: "Italian pasta dish with eggs, cheese, bacon, and black pepper.",
		Price:       12.99,
		Category:    "Pasta",
		Available:   true,
	}

	response, err := repo.CreateMeal(request)
	if err != nil {
		t.Fatalf("Error creating meal: %v", err)
	}

	expectedResponse := &pb.MealResponse{
		Id:          response.Id,
		KitchenId:   request.KitchenId,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Category:    request.Category,
		Available:   request.Available,
		CreatedAt:   response.CreatedAt,
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("CreateMeal response does not match expected.\nGot: %+v\nExpected: %+v", response, expectedResponse)
	}
}

func TestMealRepository_UpdateMeal(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Meal()

	createRequest := &pb.CreateMealRequest{
		KitchenId:   "kitchen123",
		Name:        "Pizza Margherita",
		Description: "Classic Italian pizza with tomato, mozzarella, and basil.",
		Price:       10.99,
		Category:    "Pizza",
		Available:   true,
	}
	createdMeal, err := repo.CreateMeal(createRequest)
	if err != nil {
		t.Fatalf("Error creating meal for update: %v", err)
	}

	updateRequest := &pb.UpdateMealRequest{
		Id:        createdMeal.Id,
		Name:      "Pizza Margherita Updated",
		Price:     12.99,
		Available: false,
	}

	response, err := repo.UpdateMeal(updateRequest)
	if err != nil {
		t.Fatalf("Error updating meal: %v", err)
	}

	expectedResponse := &pb.MealResponse{
		Id:        response.Id,
		Name:      updateRequest.Name,
		Price:     updateRequest.Price,
		Available: updateRequest.Available,
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("UpdateMeal response does not match expected.\nGot: %+v\nExpected: %+v", response, expectedResponse)
	}
}

func TestMealRepository_DeleteMeal(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Meal()

	createRequest := &pb.CreateMealRequest{
		KitchenId:   "kitchen123",
		Name:        "Sushi",
		Description: "Japanese dish of specially prepared vinegared rice combined with a variety of ingredients.",
		Price:       15.99,
		Category:    "Sushi",
		Available:   true,
	}
	createdMeal, err := repo.CreateMeal(createRequest)
	if err != nil {
		t.Fatalf("Error creating meal for delete: %v", err)
	}

	deleteRequest := &pb.IdRequest{
		Id: createdMeal.Id,
	}

	_, err = repo.DeleteMeal(deleteRequest)
	if err != nil {
		t.Fatalf("Error deleting meal: %v", err)
	}

	getAllRequest := &pb.GetAllMealRequest{}
	response, err := repo.GetAllMeal(getAllRequest)
	if err != nil {
		t.Fatalf("Error getting all meals after delete: %v", err)
	}

	for _, meal := range response.Meals {
		if meal.Id == deleteRequest.Id {
			t.Errorf("Meal with ID %s still exists after delete", deleteRequest.Id)
		}
	}
}

func TestMealRepository_GetAllMeal(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Error initializing storage: %v", err)
	}

	repo := storage.Meal()

	getAllRequest := &pb.GetAllMealRequest{}
	response, err := repo.GetAllMeal(getAllRequest)
	if err != nil {
		t.Fatalf("Error getting all meals: %v", err)
	}

	if response == nil || len(response.Meals) == 0 {
		t.Errorf("GetAllMeal returned nil or empty response")
	}
}
