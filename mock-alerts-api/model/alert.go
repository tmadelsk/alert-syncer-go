package model

import "time"

type Alert struct {
	Source string `json:"source"`
	Severity string `json:"severity"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

type AlertsResponse struct {
	Alerts []Alert `json:"alerts"`
}
