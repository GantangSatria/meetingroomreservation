package bootstrap

import (
	"fmt"
	"log"
	"meetingroomreservation/config"
	"meetingroomreservation/internal/models"
	"meetingroomreservation/internal/routes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewApp(cfg *config.Config) *App {

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Room{}, &models.Reservation{}, /* &models.CheckIn{} */); err != nil {
    log.Fatal("migration failed:", err)
}


	return &App{
		cfg: cfg,
		db:  db,
	}
}

func (a *App) Run() error {
	r := routes.Setup(a.db, a.cfg)
	addr := fmt.Sprintf(":%s", a.cfg.AppPort)
	return r.Run(addr)
}
