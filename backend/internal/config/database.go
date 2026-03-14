package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	cfg := Global.Database
	var gormLogLevel logger.LogLevel
	if Global.Server.Mode == "debug" {
		gormLogLevel = logger.Info
	} else {
		gormLogLevel = logger.Silent
	}

	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("[Database] Failed to connect: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("[Database] Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	log.Printf("[Database] Connected to PostgreSQL %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
}
