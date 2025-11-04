package middleware

import (
	"net/http"
	"strings"

	"github.com/guitarrich/content-sdk-go/debug"
)

// LocaleConfig contains configuration for locale middleware
type LocaleConfig struct {
	// DefaultLanguage is the fallback language
	DefaultLanguage string

	// SupportedLanguages is the list of supported languages
	SupportedLanguages []string

	// CookieName is the name of the locale cookie
	CookieName string

	// UseAcceptLanguage enables Accept-Language header parsing
	UseAcceptLanguage bool

	// CookieSecure sets the Secure attribute
	CookieSecure bool

	// CookieHTTPOnly sets the HttpOnly attribute
	CookieHTTPOnly bool

	// CookieSameSite sets the SameSite attribute
	CookieSameSite http.SameSite
}

// LocaleMiddleware handles language/locale detection
type LocaleMiddleware struct {
	config LocaleConfig
}

// NewLocaleMiddleware creates a new locale middleware
func NewLocaleMiddleware(config LocaleConfig) *LocaleMiddleware {
	// Set defaults
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = "en"
	}
	if config.CookieName == "" {
		config.CookieName = "sc_locale"
	}
	if !config.CookieSecure {
		config.CookieSecure = true
	}
	if !config.CookieHTTPOnly {
		config.CookieHTTPOnly = false // Allow JavaScript access for locale
	}
	if config.CookieSameSite == 0 {
		config.CookieSameSite = http.SameSiteLaxMode
	}

	return &LocaleMiddleware{
		config: config,
	}
}

// Handle processes the locale middleware
func (m *LocaleMiddleware) Handle(ctx Context, next HandlerFunc) error {
	path := ctx.Path()
	var locale string

	debug.Locale("processing locale for path=%s", path)

	// 1. Try to extract locale from URL path (e.g., /fr/page)
	locale = m.extractLocaleFromPath(path)
	if locale != "" && m.isSupported(locale) {
		debug.Locale("locale from path: %s", locale)
		ctx.Set(LocaleKey, locale)
		m.setCookie(ctx, locale)
		return next(ctx)
	}

	// 2. Try to get locale from query parameter
	localeParam := ctx.Request().URL.Query().Get("sc_lang")
	if localeParam == "" {
		localeParam = ctx.Request().URL.Query().Get("locale")
	}
	if localeParam != "" && m.isSupported(localeParam) {
		debug.Locale("locale from query param: %s", localeParam)
		locale = localeParam
		ctx.Set(LocaleKey, locale)
		m.setCookie(ctx, locale)
		return next(ctx)
	}

	// 3. Try to get locale from cookie
	if cookie, err := ctx.Cookie(m.config.CookieName); err == nil && cookie != nil {
		if m.isSupported(cookie.Value) {
			debug.Locale("locale from cookie: %s", cookie.Value)
			locale = cookie.Value
			ctx.Set(LocaleKey, locale)
			return next(ctx)
		}
	}

	// 4. Try Accept-Language header
	if m.config.UseAcceptLanguage {
		acceptLang := ctx.Header("Accept-Language")
		if acceptLang != "" {
			locale = m.parseAcceptLanguage(acceptLang)
			if locale != "" && m.isSupported(locale) {
				debug.Locale("locale from Accept-Language: %s", locale)
				ctx.Set(LocaleKey, locale)
				m.setCookie(ctx, locale)
				return next(ctx)
			}
		}
	}

	// 5. Fall back to default language
	locale = m.config.DefaultLanguage
	debug.Locale("using default locale: %s", locale)
	ctx.Set(LocaleKey, locale)
	m.setCookie(ctx, locale)

	return next(ctx)
}

// extractLocaleFromPath extracts locale from URL path
func (m *LocaleMiddleware) extractLocaleFromPath(path string) string {
	// Path format: /fr/page or /fr-CA/page
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		firstPart := parts[0]
		// Check if first part looks like a locale (2 chars or 2-2 chars like fr-CA)
		if len(firstPart) == 2 || (len(firstPart) == 5 && firstPart[2] == '-') {
			return firstPart
		}
	}
	return ""
}

// parseAcceptLanguage parses Accept-Language header
func (m *LocaleMiddleware) parseAcceptLanguage(header string) string {
	// Parse Accept-Language header (e.g., "en-US,en;q=0.9,fr;q=0.8")
	languages := strings.SplitSeq(header, ",")
	for lang := range languages {
		// Remove quality value if present
		parts := strings.Split(lang, ";")
		language := strings.TrimSpace(parts[0])

		// Try exact match first
		if m.isSupported(language) {
			return language
		}

		// Try language prefix (en-US -> en)
		if idx := strings.Index(language, "-"); idx > 0 {
			prefix := language[:idx]
			if m.isSupported(prefix) {
				return prefix
			}
		}
	}
	return ""
}

// isSupported checks if a locale is supported
func (m *LocaleMiddleware) isSupported(locale string) bool {
	if len(m.config.SupportedLanguages) == 0 {
		return true // All languages supported if none specified
	}

	locale = strings.ToLower(locale)
	for _, supported := range m.config.SupportedLanguages {
		if strings.ToLower(supported) == locale {
			return true
		}
	}
	return false
}

// setCookie sets the locale cookie
func (m *LocaleMiddleware) setCookie(ctx Context, locale string) {
	ctx.SetCookie(&http.Cookie{
		Name:     m.config.CookieName,
		Value:    locale,
		Path:     "/",
		Secure:   m.config.CookieSecure,
		HttpOnly: m.config.CookieHTTPOnly,
		SameSite: m.config.CookieSameSite,
		MaxAge:   365 * 24 * 60 * 60, // 1 year
	})
}
