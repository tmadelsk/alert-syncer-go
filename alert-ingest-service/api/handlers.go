package api

import (
    "context"
    "encoding/json"
    "net/http"
    "log"
    "net/url"
    "strconv"
    "time"
    "errors"

    "github.com/tmadelsk/alert-ingest-service/services/db"
)

func (s *Server) handleListAlerts(w http.ResponseWriter, r *http.Request) {

    params, err := parseParams(r.URL.Query())
    if err != nil {
        log.Println("recorded error for list alerts API %s", err)
        jsonError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    listResult, err := s.listAlertsClient.Do(r.Context(), params)
    if err != nil {
        log.Println("recorded error for list alerts API %s", err)
        jsonError(w, http.StatusInternalServerError, "internal error")
        return
    }
    alerts := listResult.(services_db.ListAlertsResult).Alerts
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alerts)
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
    go s.svc.RunOnce(context.Background())
    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(`{"status":"triggered"}`))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    status := s.monitor.GetStatus()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func parseParams(query url.Values) (services_db.ListAlertsParams, error) {

    var params services_db.ListAlertsParams

    if v := query.Get("source"); v != "" {
        params.Source = &v
    }
    if v := query.Get("severity"); v != "" {
        params.Severity = &v
    }
    if v := query.Get("since"); v != "" {
        t, err := time.Parse(time.RFC3339, v)
        if err != nil {
            return params, errors.New("invalid since param")
        }
        params.Since = &t
    }
    if v := query.Get("limit"); v != "" {
        l, err := strconv.Atoi(v)
        if err != nil {
            return params, errors.New("invalid limit param")
        }
        params.Limit = l
    }
    if v := query.Get("offset"); v != "" {
        o, err := strconv.Atoi(v)
        if err != nil {
            return params, errors.New("invalid offset param")
        }
        params.Offset = o
    }

    return params, nil
}

func jsonError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{
        "error":   message,
    })
}
