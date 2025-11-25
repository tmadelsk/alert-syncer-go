package services_db

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "github.com/tmadelsk/alert-ingest-service/services"
)

// InsertAlertParams holds the alert data that will be inserted.
type InsertAlertParams struct {
    Alert *Alert
}

// InsertAlertResult holds any metadata from the operation.
type InsertAlertResult struct {
    RowsAffected int64
}

// DBInsertAlertClient executes the insert via GORM and uses BaseClient for metrics/latency/etc.
type DBInsertAlertClient struct {
    services.BaseClient
    db *gorm.DB
}

// NewDBInsertAlertClient constructs a new client given a GORM DB connection.
func NewDBInsertAlertClient(db *gorm.DB) *DBInsertAlertClient {
    return &DBInsertAlertClient{
        BaseClient: services.BaseClient{ClientName: "InsertAlertClient"},
        db:         db,
    }
}

// Do executes the insert.
func (c *DBInsertAlertClient) Do(ctx context.Context, params interface{}) (interface{}, error) {
    p, ok := params.(InsertAlertParams)
    if !ok {
        return nil, fmt.Errorf("invalid params type for DBInsertAlertClient")
    }

    result, err := c.DoRequest(ctx, p, func(ctx context.Context, _params interface{}) (interface{}, error) {
        // TODO: add alert deduplication
        // there may be some cases where not all alerts of the previous batch were uploaded correctly
        // in this case the timestamp of the last fetched batch won't be updated and we may end with duplicated alerts
        // also, upstream may be out for a while and then we'll have a to upload a large amount of alerts at once
        res := c.db.WithContext(ctx).Create(p.Alert)
        if err := res.Error; err != nil {
            return nil, err
        }
        return InsertAlertResult{RowsAffected: res.RowsAffected}, nil
    })

    if err != nil {
        return nil, err
    }
    return result, nil
}