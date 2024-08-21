package handlers

import (
	pb "api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// GetKitchenStatisticsHandler retrieves statistics for a kitchen within a specified date range.
// @Summary Retrieve kitchen statistics
// @Description Retrieves statistics for a kitchen based on the provided kitchen_id, start_date, and end_date.
// @ID getKitchenStatistics
// @Produce json
// @Tags qo'shimcha Api
// @Param kitchen_id path string true "Kitchen ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} genproto.KitchenStatisticsResponse "Statistics for the kitchen"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/kitchen/statistics/{kitchen_id} [get]
func (h *Handler) GetKitchenStatisticsHandler(ctx *gin.Context) {
	kitchenID := ctx.Param("kitchen_id")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	parse, err := uuid.Parse(kitchenID)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("kitchen id is not uuuid"))
		return
	}
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid start_date format. Use YYYY-MM-DD"))
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

	request := &pb.GetKitchenStatisticRequest{
		UserId:    parse.String(),
		StartDate: startDate.Format(time.RFC3339),
		EndDate:   endDate.Format(time.RFC3339),
	}

	stats, err := h.OrderService.KitchenStatistic(ctx, request)
	if err != nil {
		fmt.Println("++++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to fetch kitchen statistics"))
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetUserActivityHandler retrieves activity statistics for a user within a specified date range.
// @Summary Retrieve user activity
// @Description Retrieves activity statistics for a user based on the provided start_date and end_date query parameters.
// @ID getUserActivity
// @Tags qo'shimcha Api
// @Produce json
// @Param user_id path string true "User ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} genproto.GetUserActivityResponse "Activity statistics for the user"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/users/activity/{user_id} [get]
func (h *Handler) GetUserActivityHandler(ctx *gin.Context) {
	// Parse request parameters
	req := pb.GetUserActivityRequest{
		StartDate: ctx.Query("start_date"),
		EndDate:   ctx.Query("end_date"),
	}

	userId := ctx.Param("user_id")
	parse, err := uuid.Parse(userId)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("user_id is not a valid UUID"))
		return
	}
	req.UserId = parse.String()

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid start_date format. Use YYYY-MM-DD"))
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

	req.StartDate = startDate.Format(time.RFC3339)
	req.EndDate = endDate.Format(time.RFC3339)

	stats, err := h.OrderService.ActivityUser(ctx, &req)
	if err != nil {
		fmt.Println("Error fetching user activity:", err)
		InternalServerError(ctx, fmt.Errorf("Failed to fetch user activity"))
		return
	}

	// Return activity statistics as JSON response
	ctx.JSON(http.StatusOK, stats)
}

// UpdateWorkingHoursHandler updates the working hours of a kitchen.
// @Summary Update kitchen working hours
// @Description Updates the working hours of a kitchen based on the provided request payload.
// @ID updateWorkingHours
// @Tags qo'shimcha Api
// @Accept json
// @Produce json
// @Param request body genproto.UpdateWorkingHoursRequest true "Request payload"
// @Success 200 {object} genproto.UpdateWorkingHoursResponse "Updated working hours information"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/kitchen/update-working-hours [put]
func (h *Handler) UpdateWorkingHoursHandler(ctx *gin.Context) {
	var req pb.UpdateWorkingHoursRequest
	fmt.Println("+++++++")
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid request payload"))
		return
	}
	parse, err := uuid.Parse(req.KitchenId)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("kitchen id not uuid"))
		return
	}
	req.KitchenId = parse.String()

	response, err := h.UserService.UpdateWorkingHours(ctx, &req)
	if err != nil {
		fmt.Println("++++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to update working hours"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateNutritionInfoHandler updates the nutrition information of a dish.
// @Summary Update nutrition information of a dish
// @Description Updates the nutrition information of a dish based on the provided request payload.
// @ID updateNutritionInfo
// @Tags qo'shimcha Api
// @Accept json
// @Produce json
// @Param meal_id path string true "Meal ID"
// @Param request body genproto.UpdateNutritionInfoRequest true "Request payload"
// @Success 200 {object} genproto.Dish "Updated dish information"
// @Router /api/order_service/meal/{meal_id}/update-nutrition-info [put]
func (h *Handler) UpdateNutritionInfoHandler(ctx *gin.Context) {
	fmt.Println("++++++++")

	var req pb.UpdateNutritionInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, fmt.Errorf("Invalid request payload"))
		return
	}

	mealId := ctx.Param("meal_id")
	fmt.Println("Meal ID from URL:", mealId) // Debug statement

	parse, err := uuid.Parse(mealId)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("meal_id is not a valid UUID"))
		return
	}

	req.DishId = parse.String()
	fmt.Println("Parsed UUID:", req.DishId) // Debug statement

	dish, err := h.OrderService.UpdateNutritionInfo(ctx, &req)
	if err != nil {
		fmt.Println("++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to update nutrition info"))
		return
	}

	ctx.JSON(http.StatusOK, dish)
}
