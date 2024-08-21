package handlers

import (
	"api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateCommentHandler handles the creation of a new comment.
// @Summary Create Comment
// @Description Create a new comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param Create body genproto.CreateReviewRequest true "Create Comment"
// @Success 200 {object} genproto.ReviewResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/comment/create [post]
func (h *Handler) CreateCommentHandler(ctx *gin.Context) {
	request := genproto.CreateReviewRequest{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		fmt.Println("_________", err)
		BadRequest(ctx, fmt.Errorf("sholud not blind json "))
		return
	}
	_, err = h.UserService.UserProfile(ctx, &genproto.IdRequest{Id: request.UserId})
	if err != nil {
		fmt.Println("++++++++", err)
		BadRequest(ctx, fmt.Errorf("user_id not foound in database"))
		return
	}
	_, err = h.UserService.GetByIdKitchen(ctx, &genproto.IdRequest{
		Id: request.KitchenId,
	})
	if err != nil {
		fmt.Println("_________", err)
		BadRequest(ctx, fmt.Errorf("kitchen id not found kitchen table"))
	}

	order, err := h.OrderService.CreateReview(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)

		InternalServerError(ctx, fmt.Errorf("Failed to retrieve comment"))
		return
	}
	ctx.JSON(http.StatusOK, order)
}

// GetReviewsForKitchenHandler retrieves reviews for a kitchen based on query parameters.
// @Summary Get Reviews for Kitchen
// @Description Get reviews for a kitchen based on query parameters
// @Tags Comment
// @Produce json
// @Param kitchen_id path string true "Kitchen ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} genproto.GetReviewsResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/comment/{kitchen_id} [get]
func (h *Handler) GetReviewsForKitchenHandler(ctx *gin.Context) {
	request := genproto.GetReviewsRequest{
		KitchenId: ctx.Param("kitchen_id"),
	}

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid limit parameter"))
			return
		}
		request.LimitOffset.Limit = limit
	}

	offsetStr := ctx.Query("offset")
	if offsetStr != "" {
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid offset parameter"))
			return
		}
		request.LimitOffset.Offset = offset
	}

	response, err := h.OrderService.GetReviews(ctx, &request)
	if err != nil {
		fmt.Println("++++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to retrieve reviews"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
