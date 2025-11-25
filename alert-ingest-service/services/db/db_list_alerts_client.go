package services_db

import (
    "context"
    "fmt"
    "time"

    "gorm.io/gorm"

	"github.com/tmadelsk/alert-ingest-service/services"
)

// ListAlertsParams holds filtering/pagination info
type ListAlertsParams struct {
    Source   *string
    Severity *string
    Since    *time.Time
    Limit    int
    Offset   int
}

// ListAlertsResult returns alerts
type ListAlertsResult struct {
    Alerts []Alert
}

// DBListAlertsClient executes a list query via GORM
type DBListAlertsClient struct {
    services.BaseClient
    db *gorm.DB
}

func NewDBListAlertsClient(db *gorm.DB) *DBListAlertsClient {
    return &DBListAlertsClient{
        BaseClient: services.BaseClient{ClientName: "ListAlertsClient"},
        db:         db,
    }
}

func (c *DBListAlertsClient) Do(ctx context.Context, params interface{}) (interface{}, error) {
    p, ok := params.(ListAlertsParams)
    if !ok {
        return nil, fmt.Errorf("invalid params type for DBListAlertsClient")
    }
    result, err := c.DoRequest(ctx, p, func(ctx context.Context, _params interface{}) (interface{}, error) {
        var alerts []Alert
        q := c.db.WithContext(ctx).Model(&Alert{})
        if p.Source != nil {
            q = q.Where("source = ?", *p.Source)
        }
        if p.Severity != nil {
            q = q.Where("severity = ?", *p.Severity)
        }
        if p.Since != nil {
            q = q.Where("created_at >= ?", *p.Since)
        }
        if p.Limit > 0 {
            q = q.Limit(p.Limit)
        }
        if p.Offset > 0 {
            q = q.Offset(p.Offset)
        }
        if err := q.Find(&alerts).Error; err != nil {
            return nil, err
        }
        return ListAlertsResult{Alerts: alerts}, nil
    })
    if err != nil {
        return nil, err
    }
    return result, nil
}