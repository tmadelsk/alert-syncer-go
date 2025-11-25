package health

import (
    "sync"
    "time"
)

type Monitor struct {
    mu               sync.RWMutex
    lastSuccessful   time.Time
    lastError        error
}

func NewMonitor() *Monitor {
    return &Monitor{}
}

func (m *Monitor) ReportSyncSuccess(at time.Time) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.lastSuccessful = at
    m.lastError = nil
}

func (m *Monitor) ReportSyncError(err error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.lastError = err
}

type HealthStatus struct {
    Status             string    `json:"status"`
    LastSuccessfulSync time.Time `json:"last_successful_sync"`
    LastError          string    `json:"last_error,omitempty"`
}

func (m *Monitor) GetStatus() HealthStatus {
    m.mu.RLock()
    defer m.mu.RUnlock()

    status := "ok"
    if m.lastError != nil {
        status = "degraded"
    }
    return HealthStatus{
        Status:            status,
        LastSuccessfulSync: m.lastSuccessful,
        LastError:          func() string { if m.lastError != nil { return m.lastError.Error() } ; return "" }(),
    }
}