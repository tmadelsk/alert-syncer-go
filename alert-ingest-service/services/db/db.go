package services_db

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "time"

    "github.com/tmadelsk/alert-ingest-service/config"
)

func NewDB(cfg config.DBConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
    )
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to DB: %w", err)
    }
    // Get *sql.DB to set pool options if needed
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get sql.DB from gorm DB: %w", err)
    }
    // Connection pool tuning TODO: move it to config or env variables
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    // Auto-migrate schema
    if err := db.AutoMigrate(&Alert{}, &Metadata{}); err != nil {
        return nil, fmt.Errorf("auto-migrate failed: %w", err)
    }

    return db, nil
}
