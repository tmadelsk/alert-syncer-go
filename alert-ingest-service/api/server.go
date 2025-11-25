package api

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/tmadelsk/alert-ingest-service/ingestion"
    "github.com/tmadelsk/alert-ingest-service/health"
    "github.com/tmadelsk/alert-ingest-service/services/db"
    "github.com/tmadelsk/alert-ingest-service/rate"
)

type Server struct {
    svc              *ingestion.Service
    monitor          *health.Monitor
    listAlertsClient *services_db.DBListAlertsClient
    limiter          rate.Limiter
}

func NewServer(svc *ingestion.Service, monitor *health.Monitor, listAlertsClient *services_db.DBListAlertsClient, limiter rate.Limiter) *Server {
    return &Server{ svc: svc, monitor: monitor, listAlertsClient: listAlertsClient, limiter: limiter }
}

func (s *Server) Run(addr string) error {
    r := chi.NewRouter()
    r.Handle("/alerts", NewWrapper(s.limiter, "ListAlerts", s.handleListAlerts))
    r.Handle("/sync", NewWrapper(s.limiter, "Sync", s.handleSync))
    r.Handle("/health", NewWrapper(s.limiter, "Health", s.handleHealth))

    return http.ListenAndServe(addr, r)
}