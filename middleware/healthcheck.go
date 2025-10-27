package middleware

import (
	"encoding/json"
	"net/http"
)

// HealthcheckConfig contains configuration for healthcheck middleware
type HealthcheckConfig struct {
	// Path is the healthcheck endpoint path (default: /healthz)
	Path string

	// Response is the custom response to return
	Response interface{}
}

// HealthcheckMiddleware provides a health check endpoint
type HealthcheckMiddleware struct {
	config HealthcheckConfig
}

// NewHealthcheckMiddleware creates a new healthcheck middleware
func NewHealthcheckMiddleware(config HealthcheckConfig) *HealthcheckMiddleware {
	if config.Path == "" {
		config.Path = "/healthz"
	}

	if config.Response == nil {
		config.Response = map[string]string{
			"status": "ok",
		}
	}

	return &HealthcheckMiddleware{
		config: config,
	}
}

// Handle processes the healthcheck middleware
func (m *HealthcheckMiddleware) Handle(ctx Context, next HandlerFunc) error {
	// Check if this is the healthcheck path
	if ctx.Path() == m.config.Path {
		return ctx.JSON(http.StatusOK, m.config.Response)
	}

	// Not the healthcheck path, continue to next middleware
	return next(ctx)
}

// HealthcheckResponse is the default healthcheck response
type HealthcheckResponse struct {
	Status  string            `json:"status"`
	Version string            `json:"version,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// HealthcheckHandler is a standalone handler for health checks
func HealthcheckHandler(version string) HandlerFunc {
	return func(ctx Context) error {
		response := HealthcheckResponse{
			Status:  "ok",
			Version: version,
			Details: map[string]string{
				"message": "Sitecore Content SDK Go is running",
			},
		}
		return ctx.JSON(http.StatusOK, response)
	}
}

// ReadinessHandler checks if the application is ready to receive traffic
func ReadinessHandler(checks map[string]func() bool) HandlerFunc {
	return func(ctx Context) error {
		allReady := true
		details := make(map[string]string)

		for name, checkFunc := range checks {
			if checkFunc() {
				details[name] = "ready"
			} else {
				details[name] = "not ready"
				allReady = false
			}
		}

		status := "ready"
		statusCode := http.StatusOK

		if !allReady {
			status = "not ready"
			statusCode = http.StatusServiceUnavailable
		}

		response := map[string]interface{}{
			"status":  status,
			"details": details,
		}

		ctx.Response().Header().Set("Content-Type", "application/json")
		ctx.Response().WriteHeader(statusCode)
		return json.NewEncoder(ctx.Response()).Encode(response)
	}
}
