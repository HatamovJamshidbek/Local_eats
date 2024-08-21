package api

import (
	_ "api_get_way/api/docs"
	"api_get_way/api/handlers"
	"api_get_way/genproto"
	"fmt"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger" // Import gin-swagger middleware
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

// RouterApi @title LocalEats
// @version 1.0
// @description API service
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func RouterApi(con1 *grpc.ClientConn, con2 *grpc.ClientConn) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	paymentCon := genproto.NewUserServiceClient(con1)
	reservationCon := genproto.NewOrderServiceClient(con2)
	h := handlers.NewHandler(paymentCon, reservationCon)
	fmt.Println(h)

	authRoutes := router.Group("/")
	//authRoutes.Use(middleware.MiddleWare())

	{
		meal := router.Group("/api/order_service/meal/")
		{
			meal.POST("/create", h.CreateMealHandler)
			meal.PUT("/:meal_id", h.UpdateMealHandler)
			meal.DELETE("/:meal_id", h.DeleteMealHandler)
			meal.GET("/:kitchen_id/meals", h.GetMealHandler)
		}

		order := authRoutes.Group("/api/order_service/order")
		{
			order.POST("/create", h.CreateOrderHandler)
			order.PUT("/:order_id/status", h.UpdateOrderHandler)
			order.GET("/chef/:kitchen_id ", h.GetOrdersForChefHandler)
			order.GET("/customer/:user_id ", h.GetOrdersForCustomerHandler)
			order.GET("/:id", h.GetOrderByIdHandler)
		}

		comment := authRoutes.Group("/api/order_service/comment")
		{
			comment.POST("/create", h.CreateCommentHandler)
			comment.GET("/:kitchen_id", h.GetReviewsForKitchenHandler)
		}
		payment := authRoutes.Group("/api/order_service")
		{
			payment.POST("/payments", h.CreatePaymentHandler)
		}

		kitchen := authRoutes.Group("/api/user_service/kitchen")
		{
			kitchen.POST("/create", h.CreateKitchen)
			kitchen.PUT("/:kitchen_id", h.UpdateKitchen)
			kitchen.GET("/:id", h.GetKitchenById)
			kitchen.GET("/get_all", h.GetAllKitchens)
			kitchen.GET("/search", h.SearchKitchens)
		}

		userProfile := authRoutes.Group("/api/user_service/users")
		{
			userProfile.GET("/:id/profile", h.GetUserProfile)
			userProfile.PUT("/:id/profile", h.UpdateUserProfile)
			userProfile.GET("/logout", h.LogoutHandler)
			userProfile.PUT("/update_password", h.UpdatePasswordHandler)
			userProfile.PUT("/update_token", h.UpdateTokenHandler)

		}
		extraApi := authRoutes.Group("/api")
		{
			//extraApi.PUT("/order_service/meal/update-nutrition-info/:meal_id", h.UpdateNutritionInfoHandler)
			//extraApi.GET("/user_service/kitchen/statistics/:kitchen_id", h.GetKitchenStatisticsHandler)
			//extraApi.GET(" /user_service/users/activity/:user_id", h.GetUserActivityHandler)
			extraApi.PUT("/user_service/kitchen/update-working-hours", h.UpdateWorkingHoursHandler)
		}
		router.GET("/api/user_service/kitchen/statistics/:kitchen_id", h.GetKitchenStatisticsHandler)
		router.PUT("/api/order_service/meal/:meal_id/update-nutrition-info", h.UpdateNutritionInfoHandler)
		router.GET("/api/user_service/users/activity/:user_id", h.GetUserActivityHandler)

		return router
	}
}
