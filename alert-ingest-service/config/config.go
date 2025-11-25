package config

import (
    "time"
    "os"
    "strconv"
)

type Config struct {
    UpstreamURL   string
    DB            DBConfig
    SyncInterval  time.Duration
    HTTPPort      string
}

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

func Load() *Config {
    // load from env vars with defaults
    intervalEnvVar := getEnv("SYNC_INTERVAL", "5")
    intervalInt, err := strconv.Atoi(intervalEnvVar)

    var interval time.Duration
    if err != nil {
        interval = 5 * time.Minute
    } else {
        interval = time.Duration(intervalInt) * time.Minute
    }
    // optionally parse env var like SYNC_INTERVAL
    return &Config{
        UpstreamURL:  getEnv("UPSTREAM_URL", "http://mock-alerts-api:8081/alerts"),
        DB: DBConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "password"),
            Name:     getEnv("DB_NAME", "alertsdb"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        SyncInterval: interval,
        HTTPPort:     getEnv("HTTP_PORT", ":8080"),
    }
}

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}