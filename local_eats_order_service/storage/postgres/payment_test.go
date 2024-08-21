package postgres

import (
	"reflect"
	"testing"

	pb "order_serive/genproto"
)

func TestPaymentRepository_CreatePayment(t *testing.T) {
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatalf("Storage yaratishda xatolik: %v", err)
	}

	repo := storage.Payment()
	request := &pb.CreatePaymentRequest{
		OrderId: "order123",
	}
	response, err := repo.CreatePayment(request)
	if err != nil {
		t.Fatalf("CreatePaymentda xatolik: %v", err)
	}

	expectedResponse := &pb.CreatePaymentResponse{
		Id:            response.Id,
		OrderId:       request.OrderId,
		Amount:        response.Amount,
		Status:        response.Status,
		TransactionId: response.TransactionId,
		CreatedAt:     response.CreatedAt,
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("CreatePayment natijasi kutilgan bilan mos kelmaydi.\nOldi: %+v\nKutilgan: %+v", response, expectedResponse)
	}
}
