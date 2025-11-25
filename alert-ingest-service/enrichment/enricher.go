package enrichment

import (
    "math/rand"
    "net"
    "time"
)

type Enricher interface {
    Enrich(raw RawAlert) EnrichedAlert
}

type RawAlert struct {
    Source      string
    Severity    string
    Description string
    CreatedAt   time.Time
}

type EnrichedAlert struct {
    Source         string
    Severity       string
    Description    string
    CreatedAt      time.Time
    EnrichmentType string
    IPAddress      string
    FetchedAt      time.Time
}

type SimpleEnricher struct{}

func NewSimpleEnricher() *SimpleEnricher {
    rand.Seed(time.Now().UnixNano())
    return &SimpleEnricher{}
}

func (e *SimpleEnricher) Enrich(raw RawAlert) EnrichedAlert {
    return EnrichedAlert{
        Source:         raw.Source,
        Severity:       raw.Severity,
        Description:    raw.Description,
        CreatedAt:      raw.CreatedAt,
        EnrichmentType: "random_ip_lookup",
        IPAddress:      generateRandomIP(),
        FetchedAt:      time.Now().UTC(),
    }
}

func generateRandomIP() string {
    // generate random IPv4 address
    return net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))).String()
}
