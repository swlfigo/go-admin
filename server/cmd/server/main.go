package main

import (
	"fmt"
	"log"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/database"
	"go-admin/internal/router"
	"go-admin/pkg/jwt"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	if err := database.MigrateAndSeedWithAdmin(db, cfg.Admin.Username, cfg.Admin.Password); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	jm := jwt.New(cfg.JWT.Secret,
		time.Duration(cfg.JWT.AccessMinutes)*time.Minute,
		time.Duration(cfg.JWT.RefreshDays)*24*time.Hour,
	)
	r := router.Setup(cfg, db, jm)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
