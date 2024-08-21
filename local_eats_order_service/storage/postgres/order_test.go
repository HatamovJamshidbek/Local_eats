package postgres

import (
	pb "order_serive/genproto"
	"testing"
)

func TestOrderRepository_CreateOrder(t *testing.T) {
	db, err := NewPostgresStorage()

	orderReq := pb.CreateOrderRequest{
		UserId:          "user123",
		KitchenId:       "kitchen123",
		TotalAmount:     30.00,
		Status:          "pending",
		DeliveryAddress: "123 Street",
		DeliveryTime:    "2024-07-20T15:00:00Z",
	}

	response, err := db.Order().CreateOrder(&orderReq)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	if response.UserId != orderReq.UserId || response.TotalAmount != orderReq.TotalAmount {
		t.Errorf("Response does not match the expected value.\nGot: %+v\nExpected: %+v", response, orderReq)
	}
}

func TestOrderRepository_UpdateOrder(t *testing.T) {
	db, _ := NewPostgresStorage()

	orderReq := pb.CreateOrderRequest{
		UserId:          "user123",
		KitchenId:       "kitchen123",
		TotalAmount:     30.00,
		Status:          "pending",
		DeliveryAddress: "123 Street",
		DeliveryTime:    "2024-07-20T15:00:00Z",
	}
	createResp, err := db.Order().CreateOrder(&orderReq)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	updateReq := pb.UpdateOrderStatusRequest{
		Id:     createResp.Id,
		Status: "completed",
	}

	updateResp, err := db.Order().UpdateOrderStatus(&updateReq)
	if err != nil {
		t.Fatalf("Failed to update order: %v", err)
	}

	if updateResp.Status != updateReq.Status {
		t.Errorf("Response does not match the expected value.\nGot: %+v\nExpected: %+v", updateResp.Status, updateReq.Status)
	}
}

func TestOrderRepository_GetOrderById(t *testing.T) {
	db, _ := NewPostgresStorage()

	orderReq := pb.CreateOrderRequest{}
	createResp, err := db.Order().GetOrderById(orderReq.Id)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	getReq := pb.GetOrderRequest{Id: createResp.Id}
	getResp, err := db.Order().GetOrderById(&getReq)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}

	if getResp.Id != createResp.Id {
		t.Errorf("Response does not match the expected value.\nGot: %+v\nExpected: %+v", getResp.Id, createResp.Id)
	}
}

func TestOrderRepository_GetOrdersForCustomer(t *testing.T) {
	db, err := NewPostgresStorage()

	orderReq := pb.GetOrdersRequest{
		UserId:    "user123",
		KitchenId: "kitchen123",
		Status:    "pending",
	}
	_, err = db.Order().GetOrdersForCustomer(&orderReq)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	getOrdersReq := pb.GetOrdersRequest{
		UserId: "user123",
		LimitOffset: &pb.LimitOffset{
			Limit:  10,
			Offset: 0,
		},
	}
	getOrdersResp, err := db.Order().GetOrdersForCustomer(&getOrdersReq)
	if err != nil {
		t.Fatalf("Failed to get orders for customer: %v", err)
	}

	if len(getOrdersResp.Orders) == 0 {
		t.Errorf("Expected to get orders but got none")
	}
}

func TestOrderRepository_GetOrdersForChef(t *testing.T) {
	db, _ := NewPostgresStorage()

	getOrdersReq := pb.GetOrdersRequest{
		KitchenId: "kitchen123",
		LimitOffset: &pb.LimitOffset{
			Limit:  10,
			Offset: 0,
		},
	}
	getOrdersResp, err := db.Order().GetOrdersForChef(&getOrdersReq)
	if err != nil {
		t.Fatalf("Failed to get orders for chef: %v", err)
	}

	if len(getOrdersResp.Orders) == 0 {
		t.Errorf("Expected to get orders but got none")
	}
}
