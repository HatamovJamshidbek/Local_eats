package handlers

import (
	pb "api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"strings"
)

// CreateKitchen Kitchen
// @Summary Create a new kitchen
// @Description Create a new kitchen based on the provided request
// @Tags Kitchen
// @Accept  json
// @Produce  json
// @Param input body genproto.CreateKitchenRequest true "Kitchen details to create"
// @Success 200 {object} genproto.CreateKitchenResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/kitchen/create [post]
func (h *Handler) CreateKitchen(ctx *gin.Context) {
	var request pb.CreateKitchenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, fmt.Errorf("fild not true input"))
		return
	}

	if len(request.PhoneNumber) != 0 {
		if len(request.PhoneNumber) != 12 {
			BadRequest(ctx, fmt.Errorf("telefon raqam uzunligi yetarli emas"))
			return
		}

		phoneParts := strings.Split(request.PhoneNumber, "-")
		if len(phoneParts) != 4 {
			BadRequest(ctx, fmt.Errorf("telefon raqam formati to'g'ri emas, to'g'ri format -> 90-229-62-47"))
			return
		}

		for _, part := range phoneParts {
			if _, err := strconv.Atoi(part); err != nil {
				BadRequest(ctx, fmt.Errorf("telefon raqam faqat raqamlardan iborat bo'lishi kerak, to'g'ri format -> 998-90-229-62-47"))
				return
			}
		}
		request.PhoneNumber = "+998-" + request.PhoneNumber
	}

	response, err := h.UserService.CreateKitchen(ctx, &request)
	if err != nil {
		fmt.Println("+++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to create kitchen"))
		return
	}

	ctx.JSON(200, response)
}

// UpdateKitchen godoc
// @Summary Update an existing kitchen
// @Description Update an existing kitchen based on the provided request
// @Tags Kitchen
// @Accept  json
// @Produce  json
// @Param kitchen_id path string true "Kitchen ID to update"
// @Param input body genproto.UpdateKitchenRequest true "Updated kitchen details"
// @Success 200 {object} genproto.UpdateKitchenResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/kitchen/{kitchen_id} [put]
func (h *Handler) UpdateKitchen(ctx *gin.Context) {
	var request pb.UpdateKitchenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, err)
		return
	}

	id := ctx.Param("kitchen_id")
	parse, err := uuid.Parse(id)
	if err != nil {
		return
	}
	request.Id = cast.ToString(parse)
	response, err := h.UserService.UpdateKitchen(ctx, &request)
	if err != nil {
		fmt.Println("+++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to update kitchen"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetKitchenById godoc
// @Summary Get a kitchen by ID
// @Description Retrieve a kitchen by its ID
// @Tags Kitchen
// @Accept  json
// @Produce  json
// @Param id path string true "Kitchen ID to fetch"
// @Success 200 {object} genproto.KitchenResponse
// @Failure 404 {object} string
// @Router /api/user_service/kitchen/{id} [get]
func (h *Handler) GetKitchenById(ctx *gin.Context) {
	id := ctx.Param("id")
	parse, err := uuid.Parse(id)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("error. this id not uuid"))
		return
	}
	fmt.Println(id)
	request := &pb.IdRequest{Id: cast.ToString(parse)}
	response, err := h.UserService.GetByIdKitchen(ctx, request)
	if err != nil {
		fmt.Println("+++++++", err)
		InternalServerError(ctx, fmt.Errorf("kitchen not found"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetAllKitchens godoc
// @Summary Get all kitchens
// @Description Retrieve a list of all kitchens with optional pagination
// @Tags Kitchen
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} genproto.KitchensResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/kitchen/get_all [get]

func (h *Handler) GetAllKitchens(ctx *gin.Context) {
	var request pb.LimitOffset

	limitStr := ctx.Query("limit")
	if len(limitStr) != 0 {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid limit parameter"))
			return
		}
		request.Limit = int64(limit)
	}

	offsetStr := ctx.Query("offset")

	if len(offsetStr) != 0 {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid offset parameter"))
			return
		}
		request.Offset = int64(offset)
	}
	response, err := h.UserService.GetAll(ctx, &request)
	if err != nil {
		fmt.Println("Error fetching kitchens:", err)
		InternalServerError(ctx, fmt.Errorf("failed to fetch kitchens"))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SearchKitchens godoc
// @Summary Search kitchens
// @Description Search for kitchens based on various criteria
// @Tags Kitchen
// @Accept  json
// @Produce  json
// @Param name query string false "Name of the kitchen to search for"
// @Param rating query float64 false "Rating of the kitchen"
// @Param address query string false "Address of the kitchen"
// @Param total_orders query int false "Total orders for the kitchen"
// @Param phone_number query string false "Phone number of the kitchen"
// @Param cuisine_type query string false "Cuisine type of the kitchen"
// @Param description query string false "Description of the kitchen"
// @Param owner_id query string false "Owner ID of the kitchen"
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} genproto.KitchensResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/kitchen/search [get]
func (h *Handler) SearchKitchens(ctx *gin.Context) {
	request := pb.SearchKitchenRequest{
		Name:        ctx.Query("name"),
		CuisineType: ctx.Query("cuisine_type"),
		Address:     ctx.Query("address"),
		PhoneNumber: ctx.Query("phone_number"),
		OwnerId:     ctx.Query("owner_id"),
	}

	if len(request.PhoneNumber) != 0 {
		if len(request.PhoneNumber) != 12 {
			BadRequest(ctx, fmt.Errorf("telefon raqam uzunligi yetarli emas"))
			return
		}

		phoneParts := strings.Split(request.PhoneNumber, "-")
		if len(phoneParts) != 4 {
			BadRequest(ctx, fmt.Errorf("telefon raqam formati to'g'ri emas, to'g'ri format -> 90-229-62-47"))
			return
		}

		for _, part := range phoneParts {
			if _, err := strconv.Atoi(part); err != nil {
				BadRequest(ctx, fmt.Errorf("telefon raqam faqat raqamlardan iborat bo'lishi kerak, to'g'ri format -> 998-90-229-62-47"))
				return
			}
		}
		request.PhoneNumber = "+998-" + request.PhoneNumber
	}
	request.OwnerId = ctx.Query("owner_id")
	if len(request.OwnerId) != 0 {
		if !Parse(request.OwnerId) {
			BadRequest(ctx, fmt.Errorf("bunday turdagi id yoq"))
			return
		}
	}
	total_rating := ctx.Query("total_rating")
	if len(total_rating) != 0 {
		number, err := IsNumber(total_rating)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("is not  number"))
			return
		} else {
			request.TotalOrder = int64(number)
		}
	}
	rating := ctx.Query("rating")
	if len(rating) != 0 {
		rating1, err := strconv.ParseFloat(rating, 32)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("is not float"))
			return
		} else {
			request.Rating = float32(rating1)
		}
	}
	limitStr := ctx.Query("limit")
	if len(limitStr) != 0 {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid limit parameter"))
			return
		}
		request.LimitOffset.Limit = int64(limit)
	}

	offsetStr := ctx.Query("offset")

	if len(offsetStr) != 0 {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			BadRequest(ctx, fmt.Errorf("invalid offset parameter"))
			return
		}
		request.LimitOffset.Offset = int64(offset)
	}

	response, err := h.UserService.SearchKitchen(ctx, &request)
	if err != nil {
		fmt.Println("++++++++++", err)
		InternalServerError(ctx, fmt.Errorf("Failed to search kitchens"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
