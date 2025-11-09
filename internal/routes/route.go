package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"meetingroomreservation/config"
	"meetingroomreservation/internal/controller"
	"meetingroomreservation/internal/middleware"
	"meetingroomreservation/internal/repository"
	"meetingroomreservation/internal/services"
)

type Router struct {
	Engine *gin.Engine
}

func Setup(db *gorm.DB, cfg *config.Config) *Router {
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, cfg.JWTSecret)
	userCtrl := controller.NewUserController(userSvc)

	roomRepo := repository.NewRoomRepository(db)
	roomService := services.NewRoomService(roomRepo)
	roomController := controller.NewRoomController(roomService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", userCtrl.Register)
		api.POST("/login", userCtrl.Login)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			users := protected.Group("/users")
			{
				users.GET("/", userCtrl.GetAll)
			}

			rooms := protected.Group("/rooms")
				rooms.POST("/", roomController.Create)
				rooms.PUT("/:id", roomController.Update)
				rooms.DELETE("/:id", roomController.Delete)
		}
		
		roomsPublic := api.Group("/rooms")
		roomsPublic.GET("/", roomController.GetAll)
		roomsPublic.GET("/:id", roomController.GetByID)
		
	}

	return &Router{Engine: r}
}

func (r *Router) Run(addr string) error {
	return r.Engine.Run(addr)
}
