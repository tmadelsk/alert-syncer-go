package api

import (
    "log"
    "net/http"
    "time"

    "github.com/tmadelsk/alert-ingest-service/rate"
)

// HandlerFunc is the signature for the actual handler functions.
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Wrapper wraps an existing HandlerFunc with common concerns:
// rate limiting, latency measurement, success/failure placeholder metrics.
type Wrapper struct {
    limiter    rate.Limiter
    handler    HandlerFunc
    name       string
}

func NewWrapper(limiter rate.Limiter, name string, handler HandlerFunc) *Wrapper {
    return &Wrapper{
        limiter: limiter,
        handler: handler,
        name:    name,
    }
}

func (w *Wrapper) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    ctx := req.Context()

    // 1. Rate limit check
	// TODO: move the rate limiting to a middleware, before the request hits the handler
	// TODO: define meaningful rate limiting interface (based on IP/API key/customer id etc.)
    if err := w.limiter.Acquire(ctx); err != nil {
        // rate limit failure: log + return 429
        log.Printf("[wrapper=%s] rate limiter acquired error: %v", w.name, err)
        http.Error(resp, "rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    // 2. Start timer
    start := time.Now()

    // 3. Call actual handler
    w.handler(resp, req)

    // TODO: the same as for upstream integration, handle handler failures gracefully and add here error translation

    // 4. Measure duration
    duration := time.Since(start)
    log.Printf("[wrapper=%s] handler duration: %s", w.name, duration)
    // TODO: emit latency metric for w.name

    // 5. Placeholder success/failure metrics
    // inspect resp status code by using a custom ResponseWriter if needed
    log.Printf("[wrapper=%s] handler completed", w.name)
}
