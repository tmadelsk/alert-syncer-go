package ingestion

import (
    "context"
    "fmt"
    "time"

    "github.com/tmadelsk/alert-ingest-service/services"
    "github.com/tmadelsk/alert-ingest-service/services/db"
    "github.com/tmadelsk/alert-ingest-service/services/alerts"
    "github.com/tmadelsk/alert-ingest-service/enrichment"
    "github.com/tmadelsk/alert-ingest-service/health"
)

type Service struct {
    fetchAlertsClient       services.Client
    enricher                enrichment.Enricher
    monitor                 *health.Monitor
    interval                time.Duration
    insertAlertClient       services.Client
    getLastFetchedClient    services.Client
    updateLastFetchedClient services.Client
}

func NewService(fetchAlertsClient services.Client, enricher enrichment.Enricher, interval time.Duration, monitor *health.Monitor, 
    insertAlertClient services.Client, getLastFetchedClient services.Client, updateLastFetchedClient services.Client) *Service {
        return &Service{
            fetchAlertsClient:       fetchAlertsClient,
            enricher:                enricher,
            insertAlertClient:       insertAlertClient,
            monitor:                 monitor,
            interval:                interval,
            getLastFetchedClient:    getLastFetchedClient,
            updateLastFetchedClient: updateLastFetchedClient,
        }
}

func (s *Service) Start(ctx context.Context) {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()

    // initial run
    s.RunOnce(ctx)

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.RunOnce(ctx)
        }
    }
}

// runOnce executes one ingestion cycle
func (s *Service) RunOnce(ctx context.Context) {
    lastFetchedResult, err := s.getLastFetchedClient.Do(ctx, services_db.GetLastFetchedParams{})
    if err != nil {
        s.monitor.ReportSyncError(fmt.Errorf("failed to get last fetched time: %w", err))
        return
    }
    lastFetched := lastFetchedResult.(services_db.GetLastFetchedResult).LastFetched

    // fetch raw alerts from upstream
    result, err := s.fetchAlertsClient.Do(ctx, services_alerts.FetchAlertsParams{Since: lastFetched})
    if err != nil {
        // we expect err to possibly be *services.UpstreamError
        s.monitor.ReportSyncError(fmt.Errorf("upstream fetch failed: %w", err))
        return
    }

    rawAlerts, ok := result.([]services_alerts.RawAlert)
    if !ok {
        s.monitor.ReportSyncError(fmt.Errorf("unexpected type returned from upstream: %T", result))
        return
    }

    for _, ra := range rawAlerts {
        enriched := s.enricher.Enrich(enrichment.RawAlert{
            Source:      ra.Source,
            Severity:    ra.Severity,
            Description: ra.Description,
            CreatedAt:   ra.CreatedAt,
        })

        modelAlert := &services_db.Alert{
            Source:         enriched.Source,
            Severity:       enriched.Severity,
            Description:    enriched.Description,
            CreatedAt:      enriched.CreatedAt,
            EnrichmentType: enriched.EnrichmentType,
            IPAddress:      enriched.IPAddress,
            FetchedAt:      enriched.FetchedAt,
        }

        _, err := s.insertAlertClient.Do(ctx, services_db.InsertAlertParams{Alert: modelAlert})
        if err != nil {
            s.monitor.ReportSyncError(fmt.Errorf("failed to insert alert: %w", err))
            return
        }
        
    }

    // update last fetched timestamp only if we were able to successfully pull and upload all new alerts
    now := time.Now().UTC()
    _, err = s.updateLastFetchedClient.Do(ctx, services_db.UpdateLastFetchedParams{At: now})
    if err != nil {
        s.monitor.ReportSyncError(fmt.Errorf("failed to update last fetched time: %w", err))
        return
    }

    s.monitor.ReportSyncSuccess(now)
}