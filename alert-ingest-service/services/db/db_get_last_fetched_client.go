package services_db

import (
    "context"
    "fmt"
    "time"

    "gorm.io/gorm"

	"github.com/tmadelsk/alert-ingest-service/services"
)

type GetLastFetchedParams struct {
    // none needed ?
}

type GetLastFetchedResult struct {
    LastFetched time.Time
}

type DBGetLastFetchedClient struct {
    services.BaseClient
    db *gorm.DB
}

func NewDBGetLastFetchedClient(db *gorm.DB) *DBGetLastFetchedClient {
    return &DBGetLastFetchedClient{
        BaseClient: services.BaseClient{ClientName: "GetLastFetchedClient"},
        db:         db,
    }
}

func (c *DBGetLastFetchedClient) Do(ctx context.Context, params interface{}) (interface{}, error) {
    _, ok := params.(GetLastFetchedParams)
    if !ok {
        return nil, fmt.Errorf("invalid params type for DBGetLastFetchedClient")
    }
    result, err := c.DoRequest(ctx, params, func(ctx context.Context, _params interface{}) (interface{}, error) {
        var meta Metadata
        err := c.db.WithContext(ctx).First(&meta, "key = ?", "last_fetched_at").Error
        if err != nil {
            if err == gorm.ErrRecordNotFound {
                return GetLastFetchedResult{LastFetched: time.Time{}}, nil
            }
            return nil, err
        }
        t, err := time.Parse(time.RFC3339, meta.Value)
        if err != nil {
            return nil, err
        }
        return GetLastFetchedResult{LastFetched: t}, nil
    })
    if err != nil {
        return nil, err
    }
    return result, nil
}
