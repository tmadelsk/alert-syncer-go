package services_db

import "time"

type Alert struct {
    ID             uint64    `gorm:"primaryKey;autoIncrement"`
    Source         string
    Severity       string
    Description    string
    CreatedAt      time.Time
    EnrichmentType string
    IPAddress      string
    FetchedAt      time.Time
}

type Metadata struct {
    Key   string    `gorm:"primaryKey;size:128"`
    Value string    `gorm:"size:256"`
    UpdatedAt time.Time
}
