package services

import (
    "database/sql"
    "errors"
    "fmt"
    "io"
    "net"
    "strings"
)

type FailureType string

const (
    RequestError        FailureType = "RequestError"
    InternalServerError FailureType = "InternalServerError"
    NonRetryableError   FailureType = "NonRetryableError"
    UnknownError        FailureType = "UnknownError"
)

type UpstreamError struct {
    Type    FailureType
    Message string
    Err     error
}

func (e *UpstreamError) Error() string {
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// translateError inspects err and returns the appropriate FailureType
// only an example how it may look like, not a real implementation
func translateError(err error) FailureType {
    if err == nil {
        return ""
    }
    // Unwrap errors if needed
    // Example: network timeout or connection refused => RequestError
    var netErr net.Error
    if errors.As(err, &netErr) {
        if netErr.Timeout() {
            return RequestError
        }
        return RequestError
    }
    // SQL/DB errors of certain type
    if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrNoRows) {
        // Consider connection errors as InternalServerError
        return InternalServerError
    }
    // HTTP status code error inside our wrapper
    if ue, ok := err.(*UpstreamError); ok {
        // If bursting through nested UpstreamError, propagate
        return ue.Type
    }
    // Check message patterns (not ideal but placeholder)
    msg := err.Error()
    if strings.Contains(msg, "status 5") || strings.Contains(msg, "server error") {
        return InternalServerError
    }
    // dummy implementation, this case is much more complex in reality
    if strings.Contains(msg, "non-retryable") {
        return NonRetryableError
    }
    // For other known conditions:
    if errors.Is(err, io.EOF) {
        return RequestError
    }
    // Fallback to UnknownError
    return UnknownError
}

// wrapError wraps an existing err with context and the appropriate type
func wrapError(err error) error {
    if err == nil {
        return nil
    }
    ft := translateError(err)
    return &UpstreamError{
        Type:    ft,
        Message: err.Error(),
        Err:     err,
    }
}
