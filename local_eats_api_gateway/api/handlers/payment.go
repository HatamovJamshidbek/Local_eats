package handlers

import (
	"api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreatePaymentHandler handles the creation of payments.
// @Summary Create Payment
// @Description Create a new payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param payment body genproto.CreatePaymentRequest true "Payment Request"
// @Success 200 {object} genproto.CreatePaymentResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/payments [post]
func (h *Handler) CreatePaymentHandler(ctx *gin.Context) {
	var request genproto.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid request body"))
		return
	}
	fmt.Println("++++++", request.PaymentMethod)
	// Validate payment_method is not empty
	if len(request.PaymentMethod) == 0 {
		BadRequest(ctx, fmt.Errorf("Payment method is required"))
		return
	}

	response, err := h.OrderService.CreatePayment(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)
		InternalServerError(ctx, fmt.Errorf("payment not created"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
