package main

import (
    "context"
    "log"
    "os/signal"
    "syscall"

    "github.com/tmadelsk/alert-ingest-service/config"
    "github.com/tmadelsk/alert-ingest-service/ingestion"
    "github.com/tmadelsk/alert-ingest-service/rate"
    "github.com/tmadelsk/alert-ingest-service/api"
    "github.com/tmadelsk/alert-ingest-service/health"
    "github.com/tmadelsk/alert-ingest-service/enrichment"
    "github.com/tmadelsk/alert-ingest-service/services/db"
    "github.com/tmadelsk/alert-ingest-service/services/alerts"
)

func main() {
    cfg := config.Load()

    // Setup DB connection
    db, err := services_db.NewDB(cfg.DB)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    // Setup upstream client(s)
    alertsClient := services_alerts.NewAlertsClient(cfg.UpstreamURL)

    // Setup rate limiter (noop for now)
    limiter := &rate.NoopLimiter{}

    // Setup health monitor
    monitor := health.NewMonitor()

    // Setup simple enricher
    enricher := enrichment.NewSimpleEnricher()

    // create all necessary clients to interact with the db
    insertAlertClient := services_db.NewDBInsertAlertClient(db)
    listAlertsClient := services_db.NewDBListAlertsClient(db)
    getLastFetchedClient := services_db.NewDBGetLastFetchedClient(db)
    updateLastFetchedClient := services_db.NewDBUpdateLastFetchedClient(db)

    // Setup ingestion service
    svc := ingestion.NewService(alertsClient, enricher, cfg.SyncInterval, monitor, insertAlertClient, getLastFetchedClient, updateLastFetchedClient)

    // Start HTTP API server
    apiSrv := api.NewServer(svc, monitor, listAlertsClient, limiter)
    go func() {
        log.Printf("HTTP API listening on %s", cfg.HTTPPort)
        if err := apiSrv.Run(cfg.HTTPPort); err != nil {
            log.Fatalf("HTTP server failed: %v", err)
        }
    }()

    // Start ingestion loop
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    go svc.Start(ctx)

    <-ctx.Done()
    log.Println("Shutting down…")
}