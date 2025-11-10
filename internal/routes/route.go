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

	reservationRepo := repository.NewReservationRepository(db)
	reservationService := services.NewReservationService(reservationRepo)
	reservationController := controller.NewReservationController(reservationService)

	checkinRepo := repository.NewCheckinRepository(db)
	checkinService := services.NewCheckinService(checkinRepo, reservationRepo)
	checkinController := controller.NewCheckinController(checkinService)

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

			reservations := protected.Group("/reservations")
			{
				reservations.POST("/", reservationController.Create)
				reservations.PUT("/:id", reservationController.Update)
				reservations.DELETE("/:id", reservationController.Delete)
				reservations.GET("/:id/qrcode", reservationController.GetQRCode)
			}

			checkins := protected.Group("/checkin")
			{
				checkins.POST("/:reservation_id", checkinController.Checkin)
				checkins.POST("/:reservation_id/checkout", checkinController.Checkout)
				checkins.POST("/qrcode", checkinController.CheckinByQRCode)
			}

			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnlyMiddleware())
			{
				admin.PUT("/reservations/:id/approve", reservationController.Approve)
				admin.PUT("/reservations/:id/reject", reservationController.Reject)
			}
		}

		roomsPublic := api.Group("/rooms")
		roomsPublic.GET("/", roomController.GetAll)
		roomsPublic.GET("/:id", roomController.GetByID)

		reservationsPublic := api.Group("/reservations")
		reservationsPublic.GET("/", reservationController.GetAll)
		reservationsPublic.GET("/:id", reservationController.GetByID)
		
		
	}

	return &Router{Engine: r}
}

func (r *Router) Run(addr string) error {
	return r.Engine.Run(addr)
}
