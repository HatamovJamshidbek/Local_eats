package handlers

import (
	"api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateMealHandler handles the creation of a new meal item.
// @Summary Create Meal
// @Description Create a new meal
// @Tags Meal
// @Accept json
// @Produce json
// @Param Create body genproto.CreateMealRequest true "Create Menu"
// @Success 200 {object} genproto.MealResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/meal/create [post]
func (h *Handler) CreateMealHandler(ctx *gin.Context) {
	var request genproto.CreateMealRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, err)
		fmt.Println("_________", err)
		return
	}
	_, err := h.UserService.GetByIdKitchen(ctx, &genproto.IdRequest{Id: request.KitchenId})
	if err != nil {
		fmt.Println("_________", err)
		BadRequest(ctx, fmt.Errorf("kitchen_id mavjud emas"))
		return
	}

	response, err := h.OrderService.CreateMeal(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)

		InternalServerError(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

// UpdateMealHandler handles the update of a meal item.
// @Summary Update Meal
// @Description Update an existing meal
// @Tags Meal
// @Accept json
// @Produce json
// @Param meal_id path string true "Meal ID"
// @Param Update body genproto.UpdateMealRequest true "Update Menu"
// @Success 200 {object} genproto.MealResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/meal/{meal_id} [put]
func (h *Handler) UpdateMealHandler(ctx *gin.Context) {
	var request genproto.UpdateMealRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, err)
		fmt.Println("_________", err)
		return
	}
	request.Id = ctx.Param("meal_id")
	response, err := h.OrderService.UpdateMeal(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)

		InternalServerError(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

// DeleteMealHandler handles the deletion of a meal item.
// @Summary Delete Meal
// @Description Delete a meal
// @Tags Meal
// @Produce json
// @Param meal_id path string true "Meal ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/meal/{meal_id} [delete]
func (h *Handler) DeleteMealHandler(ctx *gin.Context) {
	var request genproto.IdRequest

	request.Id = ctx.Param("meal_id")
	_, err := h.OrderService.Delete(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)

		InternalServerError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Dish successfully deleted",
	})
}

// GetMealHandler retrieves meals based on query parameters.
// @Summary Get Meals
// @Description Get meals based on query parameters
// @Tags Meal
// @Produce json
// @Param name query string false "Name"
// @Param category query string false "Category"
// @Param available query string false "Available"
// @Param price query string false "Price"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param kitchen_id path string true "Kitchen ID"
// @Success 200 {object} genproto.MealsResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/order_service/meal/{kitchen_id}/meals [get]
func (h *Handler) GetMealHandler(ctx *gin.Context) {

	request := genproto.GetAllMealRequest{
		Name:     ctx.Query("name"),
		Category: ctx.Query("category"),
	}
	price := ctx.Query("price")
	if price != "" {
		number, err := IsNumber(price)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("pric error xatto"))
			return
		}
		request.Price = float32(number)
	}

	limit := ctx.Query("limit")
	if len(limit) != 0 {
		limit1, err := IsNumber(limit)
		if err != nil {
			BadRequest(ctx, err)
			return
		}
		request.LimitOffset.Limit = int64(limit1)
	}
	fmt.Println("+++++++++++")
	offset := ctx.Query("offset")
	if len(offset) != 0 {
		offset1, err := IsNumber(offset)
		if err != nil {
			BadRequest(ctx, err)
			return
		}
		request.LimitOffset.Offset = int64(offset1)
	}

	kitchenId := ctx.Param("kitchen_id")
	if len(kitchenId) != 0 {
		parse, err := uuid.Parse(kitchenId)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("kitchen id bunday bo'lishi mumkin emas"))
			return
		}
		request.KitchenId = parse.String()
	}

	fmt.Println("+++++++++++")
	response, err := h.OrderService.GetAllMeal(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++++", err)
		InternalServerError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"meals": response,
	})
}
