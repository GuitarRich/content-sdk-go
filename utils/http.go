package utils

import (
	"net/http"
	"strconv"
	"strings"
)

// headerSupportedTypes constrains the types supported by GetHeaderOrDefault.
type headerSupportedTypes interface {
	~string | ~bool | ~int | ~int64 | ~float64
}

// GetHeaderOrDefault returns the header value for the given key parsed into T,
// or the provided default if absent or unparsable.
func GetHeaderOrDefault[T headerSupportedTypes](req *http.Request, key string, defaultValue T) T {
	if req == nil {
		return defaultValue
	}
	raw := strings.TrimSpace(req.Header.Get(key))
	if raw == "" {
		return defaultValue
	}

	var anyDefault any = defaultValue
	switch anyDefault.(type) {
	case string:
		return any(raw).(T)
	case bool:
		if v, err := strconv.ParseBool(raw); err == nil {
			return any(v).(T)
		}
	case int:
		if v, err := strconv.ParseInt(raw, 10, 0); err == nil {
			return any(int(v)).(T)
		}
	case int64:
		if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return any(v).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(raw, 64); err == nil {
			return any(v).(T)
		}
	default:
		// unsupported type; fall through to return default
	}
	return defaultValue
}
