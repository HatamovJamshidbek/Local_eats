package handlers

import (
	pb "api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

// GetUserProfile godoc
// @Summary Get user profile by ID
// @Description Retrieves a user's profile information based on the provided ID
// @Tags UserProfile
// @ID get-user-profile
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} genproto.UserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/users/{id}/profile [get]
func (h *Handler) GetUserProfile(ctx *gin.Context) {
	userId := ctx.Param("id")
	parse, err := uuid.Parse(userId)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("user_id xatto "))
		return
	}
	request := &pb.IdRequest{Id: parse.String()}

	response, err := h.UserService.UserProfile(ctx, request)
	if err != nil {
		fmt.Println("++++++", err)
		InternalServerError(ctx, fmt.Errorf("failed to fetch user profile"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Updates a user's profile information based on the provided data
// @Tags UserProfile
// @ID update-user-profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body genproto.UpdateUserProfileRequest true "Update User Profile Request"
// @Success 200 {object} genproto.UserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/users/{id}/profile [put]
func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
	var request pb.UpdateUserProfileRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, fmt.Errorf("invalid input data"))
		return
	}
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

	userId := ctx.Param("id")
	parse, err := uuid.Parse(userId)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("user_id not is uuid"))
		return
	}
	request.Id = parse.String()

	response, err := h.UserService.UpdateUserProfile(ctx, &request)
	if err != nil {
		fmt.Println("+++++", err)
		InternalServerError(ctx, fmt.Errorf("failed to update user profile"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateTokenHandler godoc
// @Summary Refresh access token
// @Description Refreshes the access token using a refresh token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param refresh_token query string true "Refresh token"
// @Success 200 {object} genproto.LoginResponse
// @Failure 400 {object} string
// @Router /api/user_service/users/update_token [PUT]
func (h *Handler) UpdateTokenHandler(ctx *gin.Context) {
	request := pb.UpdateTokenRequest{
		RefreshToken: ctx.Query("refresh_token"),
	}
	response, err := h.UserService.UpdateToken(ctx, &request)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("register is that malumot not found"))
		return
	}
	ctx.JSON(200, response)
}
func (h *Handler) ResetPasswordHandler(ctx *gin.Context) {
	request := pb.UpdatePasswordRequest{
		Email:    ctx.Query("email"),
		Password: ctx.Query("password"),
	}
	response, err := h.UserService.UpdatePassword(ctx, &request)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("register is that malumot not found"))
		return
	}
	ctx.JSON(200, response)
}

// UpdatePasswordHandler godoc
// @Summary Update user password
// @Description Updates the password for a user based on the provided email and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body genproto.UpdatePasswordRequest true "Update Password Request"
// @Success 200 {object} genproto.Void
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/users/update_password [put]
func (h *Handler) UpdatePasswordHandler(ctx *gin.Context) {
	var request pb.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, fmt.Errorf("invalid input data"))
		return
	}

	_, err := h.UserService.UpdatePassword(ctx, &request)
	if err != nil {
		InternalServerError(ctx, fmt.Errorf("failed to update password"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// LogoutHandler godoc
// @Summary Logout
// @Description Logout user based on the provided email and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/user_service/users/logout [put]
func (h *Handler) LogoutHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"logout": "logout successfully",
	})
}
