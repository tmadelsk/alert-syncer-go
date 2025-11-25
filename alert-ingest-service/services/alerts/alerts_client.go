package services_alerts

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/tmadelsk/alert-ingest-service/services"
)

type AlertsClient struct {
    services.BaseClient
    baseURL string
}

type FetchAlertsParams struct {
    Since time.Time
}

type RawAlert struct {
    Source      string    `json:"source"`
    Severity    string    `json:"severity"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}

func NewAlertsClient(baseURL string) *AlertsClient {
    return &AlertsClient{ 
        BaseClient: services.BaseClient{ClientName: "AlertsClient"},
        baseURL: baseURL,
    }
}

func (c *AlertsClient) Do(ctx context.Context, params interface{}) (interface{}, error) {
    p, ok := params.(FetchAlertsParams)
    if !ok {
        return nil, fmt.Errorf("invalid params type for AlertsClient")
    }
    return c.DoRequest(ctx, p, func(ctx context.Context, _params interface{}) (interface{}, error) {
        url := fmt.Sprintf("%s?since=%s", c.baseURL, p.Since.UTC().Format(time.RFC3339))
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return nil, err
        }
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        if resp.StatusCode >= 500 {
            return nil, fmt.Errorf("server error: status %d", resp.StatusCode)
        }
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        var payload struct {
            Alerts []RawAlert `json:"alerts"`
        }
        if err := json.Unmarshal(body, &payload); err != nil {
            return nil, err
        }
        return payload.Alerts, nil
    })
}