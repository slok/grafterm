package v1

import (
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// API is a wrapper of v1.API.
type API interface {
	promv1.API
}
