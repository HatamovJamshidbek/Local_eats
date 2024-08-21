package handlers

import (
	pb "auth_serice/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param input body genproto.RegisterRequest true "Registration details"
// @Success 200 {object} genproto.RegisterResponse
// @Failure 400 {object} string
// @Router /api/auth_service/register [post]
func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var request pb.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		BadRequest(ctx, fmt.Errorf("request error: %v", err))
		return
	}

	response, err := h.userStorage.Users().Register(&request)
	if err != nil {
		InternalServerError(ctx, fmt.Errorf("Failed to register: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// LoginHandler godoc
// @Summary User login
// @Description Logs in a user with the provided credentials
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param input body genproto.LoginRequest true "Login credentials"
// @Success 200 {object} genproto.LoginResponse
// @Failure 400 {object} string
// @Router /api/auth_service/login [post]
func (h *Handler) LoginHandler(ctx *gin.Context) {
	request := pb.LoginRequest{}
	response, err := h.userStorage.Users().Login(&request)
	if err != nil {
		BadRequest(ctx, fmt.Errorf("register is that malumot not found"))
		return
	}
	ctx.JSON(200, response)
}
