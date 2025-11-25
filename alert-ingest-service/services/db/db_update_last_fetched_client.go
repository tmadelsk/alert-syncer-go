package services_db

import (
    "context"
    "fmt"
    "time"

    "gorm.io/gorm"

	"github.com/tmadelsk/alert-ingest-service/services"
)

type UpdateLastFetchedParams struct {
    At time.Time
}

type UpdateLastFetchedResult struct {
    // nothing needed ?
}

type DBUpdateLastFetchedClient struct {
    services.BaseClient
    db *gorm.DB
}

func NewDBUpdateLastFetchedClient(db *gorm.DB) *DBUpdateLastFetchedClient {
    return &DBUpdateLastFetchedClient{
        BaseClient: services.BaseClient{ClientName: "UpdateLastFetchedClient"},
        db:         db,
    }
}

func (c *DBUpdateLastFetchedClient) Do(ctx context.Context, params interface{}) (interface{}, error) {
    p, ok := params.(UpdateLastFetchedParams)
    if !ok {
        return nil, fmt.Errorf("invalid params type for DBUpdateLastFetchedClient")
    }
    result, err := c.DoRequest(ctx, p, func(ctx context.Context, _params interface{}) (interface{}, error) {
        meta := Metadata{
            Key:       "last_fetched_at",
            Value:     p.At.UTC().Format(time.RFC3339),
            UpdatedAt: time.Now().UTC(),
        }
        // we use Save so upsert
        if err := c.db.WithContext(ctx).Save(&meta).Error; err != nil {
            return nil, err
        }
        return UpdateLastFetchedResult{}, nil
    })
    if err != nil {
        return nil, err
    }
    return result, nil
}
