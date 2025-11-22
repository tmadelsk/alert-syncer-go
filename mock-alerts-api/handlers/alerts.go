package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/tmadelsk/mock-alerts-api/model"
	"github.com/tmadelsk/mock-alerts-api/util"
	"math/rand"
)

var sources = []string{
	"siem-1","sime-2","siem-3","sime-4","sime-5",
	"siem-6","sime-7","siem-8","sime-9","sime-10",
}

var severities = []string{"low","medium","high","critical"}

func AlertsHandler(w http.ResponseWriter, r *http.Request) {
	sinceParam := r.URL.Query().Get("since")
	sinceTime, err := time.Parse(time.RFC3339, sinceParam)
	if err != nil {
		http.Error(w, "invalid since parameter", http.StatusBadRequest)
		return
	}

	// simulate failure
	if util.ShouldFail() {
		http.Error(w, "simulated upstream error", http.StatusInternalServerError)
		return
	}

	// simulate delay
	if util.ShouldDelay() {
		time.Sleep(util.RandomDelay())
	}

	// generate random number of alerts
	count := rand.Intn(5)
	alerts := make([]model.Alert, 0, count)

	for i:=0; i<count; i++ {
		source := sources[rand.Intn(len(sources))]
		severity := severities[rand.Intn(len(severities))]
		created := sinceTime.Add(time.Duration(rand.Intn(60)) * time.Minute)
		alerts = append(
			alerts, model.Alert{
				Source: source,
				Severity: severity,
				Description: "Simulated alert for " + source,
				CreatedAt: created,
			},
		)
	}

	resp := model.AlertsResponse{Alerts: alerts}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
