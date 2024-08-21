package handlers

import (
	"api_get_way/genproto"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

// CreateOrderHandler handles the creation of a new order.
// @Summary Create Order
// @Description Create a new order
// @Tags Order
// @Accept json
// @Produce json
// @Param Create body genproto.CreateOrderRequest true "Create Order"
// @Success 200 {object} genproto.OrderResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/order/create [post]
func (h *Handler) CreateOrderHandler(ctx *gin.Context) {
	var request genproto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, err)
		return
	}
	parse, err := uuid.Parse(request.UserId)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("user id error"))
		return
	}
	_, err = h.UserService.UserProfile(ctx, &genproto.IdRequest{Id: parse.String()})
	if err != nil {
		fmt.Println("+++++++", err)
		BadRequest(ctx, fmt.Errorf("user_id not found in database"))
		return
	}

	if _, err := uuid.Parse(request.KitchenId); err != nil {
		BadRequest(ctx, errors.New("invalid KitchenId format"))
		return
	}

	_, err = h.UserService.GetByIdKitchen(ctx, &genproto.IdRequest{Id: request.KitchenId})
	if err != nil {
		BadRequest(ctx, errors.New("kitchen_id not found"))
		return
	}
	//2023-07-17T10:30:00Z
	deliveryTime, err := time.Parse(time.RFC3339, request.DeliveryTime)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("invalid DeliveryTime format: %v", err))
		return
	}

	request.DeliveryTime = deliveryTime.Format(time.RFC3339)
	response, err := h.OrderService.CreateOrder(ctx, &request)
	if err != nil {
		InternalServerError(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

// UpdateOrderHandler handles the update of an order's status.
// @Summary Update Order Status
// @Description Update an order's status
// @Tags Order
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param Update body genproto.UpdateOrderStatusRequest true "Update Order Status"
// @Success 200 {object} genproto.OrderStatusResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/order/{order_id}/status [put]
func (h *Handler) UpdateOrderHandler(ctx *gin.Context) {
	var request genproto.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		fmt.Println("_________", err)
		BadRequest(ctx, err)
		return
	}

	id := ctx.Param("order_id")
	parse, err := uuid.Parse(id)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("order_id xatto kiritildi"))
		return
	}
	request.Id = parse.String()
	response, err := h.OrderService.UpdateOrderStatus(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)

		InternalServerError(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

// GetOrdersForChefHandler retrieves orders for a chef based on query parameters.
// @Summary Get Orders for Chef
// @Description Get orders for a chef based on query parameters
// @Tags Order
// @Produce json
// @Param kitchen_id path string true "Kitchen ID"
// @Param status query string false "Status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} genproto.GetOrdersResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/order/chef/{kitchen_id} [get]
func (h *Handler) GetOrdersForChefHandler(ctx *gin.Context) {
	request := genproto.GetOrdersRequest{
		Status: ctx.Query("status"),
	}

	kitchenId := ctx.Param("kitchen_id")
	fmt.Println("++++++++", kitchenId)

	// Validate kitchen_id as UUID
	if _, err := uuid.Parse(kitchenId); err != nil {
		fmt.Println("+++++++", err)
		BadRequest(ctx, fmt.Errorf("kitchen_id is not a valid UUID"))
		return
	}
	request.KitchenId = kitchenId

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid limit parameter"))
			return
		}
		request.LimitOffset.Limit = int64(limit)
	}

	offsetStr := ctx.Query("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid offset parameter"))
			return
		}
		request.LimitOffset.Offset = int64(offset)
	}

	// Optionally handle user_id if needed
	userId := ctx.Query("user_id")
	if userId != "" {
		if _, err := uuid.Parse(userId); err != nil {
			BadRequest(ctx, fmt.Errorf("user_id is not a valid UUID"))
			return
		}
		request.UserId = userId
	}

	orders, err := h.OrderService.GetOrdersForChef(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to retrieve orders"))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// GetOrdersForCustomerHandler retrieves orders for a customer based on query parameters.
// @Summary Get Orders for Customer
// @Description Get orders for a customer based on query parameters
// @Tags Order
// @Produce json
// @Param user_id path string true "User ID"
// @Param status query string false "Status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param kitchen_id query string false "Kitchen ID"
// @Success 200 {object} genproto.GetOrdersResponse "Successfully retrieved orders"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/order_service/order/customer/{user_id} [get]
func (h *Handler) GetOrdersForCustomerHandler(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	fmt.Println("+++++++ User id ", userId) // Log user_id to check its value
	if userId == "" {
		BadRequest(ctx, fmt.Errorf("user_id is empty"))
		return
	}

	if _, err := uuid.Parse(userId); err != nil {
		BadRequest(ctx, fmt.Errorf("user_id is not a valid UUID: %v", err))
		return
	}

	request := genproto.GetOrdersRequest{
		Status: ctx.Query("status"),
	}
	request.UserId = userId

	// Check if user exists
	if _, err := h.UserService.UserProfile(ctx, &genproto.IdRequest{Id: request.UserId}); err != nil {
		BadRequest(ctx, fmt.Errorf("user_id not found in database: %v", err))
		return
	}

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid limit parameter: %v", err))
			return
		}
		request.LimitOffset.Limit = int64(limit)
	}

	offsetStr := ctx.Query("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid offset parameter: %v", err))
			return
		}
		request.LimitOffset.Offset = int64(offset)
	}

	kitchenId := ctx.Query("kitchen_id")
	if kitchenId != "" {
		if _, err := uuid.Parse(kitchenId); err != nil {
			BadRequest(ctx, fmt.Errorf("kitchen_id is not a valid UUID: %v", err))
			return
		}
		request.KitchenId = kitchenId
	}

	// Retrieve orders
	orders, err := h.OrderService.GetOrdersForCustomer(ctx, &request)
	if err != nil {
		InternalServerError(ctx, fmt.Errorf("Failed to retrieve orders: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// GetOrderByIdHandler retrieves an order by its ID.
// @Summary Get Order by ID
// @Description Get an order by its ID
// @Tags Order
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} genproto.OrderResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/order/{id} [get]
func (h *Handler) GetOrderByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(id)
	parse, err := uuid.Parse(id)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("order_id xatto kiritildi"))
		return
	}

	request := &genproto.GetOrderRequest{Id: parse.String()}
	order, err := h.OrderService.GetOrderById(ctx, request)
	if err != nil {
		fmt.Println("++++++++++++", err)

		InternalServerError(ctx, fmt.Errorf("Failed to retrieve orders"))

		return
	}
	ctx.JSON(http.StatusOK, order)
}
