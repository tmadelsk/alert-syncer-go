package services

import (
	"context"
	"time"
	"log"
)

// Client is the interface all upstream clients implement
type Client interface {
	// Invoke the upstream and return some generic result or error
	Do(ctx context.Context, params interface{}) (interface{}, error)
}

// BaseClient provides shared behaviour; you embed or compose it
type BaseClient struct {
	// placeholder fields: metricsCollector, httpClient, etc.
	ClientName string
}

// DoRequest is the template method: it wraps the concrete makeRequest
func (b *BaseClient) DoRequest(ctx context.Context, params interface{}, 
	makeRequest func(ctx context.Context, params interface{}) (interface{}, error)) (interface{}, error) {

		// TODO: move it to configs
		const maxRetries = 3
		baseDelay := 500 * time.Millisecond

		var lastErr error

		for attempt := 0; attempt < maxRetries; attempt++ {
			start := time.Now()
			result, err := makeRequest(ctx, params)
			duration := time.Since(start)

			// TODO: emit latency metric with client name dimension instead of logging it
			log.Printf("last call made by %s took %s", b.ClientName, duration)

			if err == nil {
            	// TODO: emit success metric with client name dimension
            	return result, nil
        	}

			// translate error to FailureType
			failureType := translateError(err)
			// TODO: emit failure meric with failure type and client name dimensions. 
			// Failure type will allow us to decide if we should trigger internal alert 
			// (e.g. for internal server or unknown errors) or not (e.g. RequestError)
			log.Printf("error of type %s happened for client type %s", failureType, b.ClientName)

			// if the error type is non-retryable, break early
        	switch failureType {
        	case RequestError:
            	// maybe retryable
        	case InternalServerError:
            	// maybe retryable
        	case UnknownError:
            	// maybe retryable or maybe not
			case NonRetryableError:
				return nil, &UpstreamError{Type: failureType, Err: err}
        	default:
            	// unknown type — treat as non-retryable for safety
            	return nil, &UpstreamError{Type: failureType, Err: err}
        	}

			lastErr = err

        	// wait for backoff
        	delay := baseDelay * (1 << attempt)
			// we can also add a jitter here
        	time.Sleep(delay)
		}

		failureType := translateError(lastErr)
		return nil, &UpstreamError{Type: failureType, Err: lastErr}
	}
