package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func OK(gn *gin.Context) {
	gn.JSON(200, gin.H{
		"status":  http.StatusOK,
		"time":    time.Now(),
		"success": true,
	})
	gn.Header("Content-Type", "application/json")

}
func Created(gn *gin.Context) {
	gn.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"time":    time.Now(),
		"success": true,
	})
	gn.Header("Content-Type", "application/json")

}
func InternalServerError(gn *gin.Context, err error) {
	fmt.Println("salom")
	gn.JSON(http.StatusInternalServerError, gin.H{
		"status":  http.StatusInternalServerError,
		"time":    time.Now(),
		"message": err.Error(),
		"success": false,
	})
	gn.Header("Content-Type", "application/json")

}
func BadRequest(gn *gin.Context, err error) {
	gn.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"time":    time.Now(),
		"message": err.Error(),
		"success": false,
	})
	gn.Header("Content-Type", "application/json")

}
