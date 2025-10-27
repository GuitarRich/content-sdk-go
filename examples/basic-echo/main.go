package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/content-sdk-go/client"
	"github.com/content-sdk-go/config"
	"github.com/content-sdk-go/handlers"
	"github.com/content-sdk-go/middleware"
	"github.com/content-sdk-go/seo"
	"github.com/content-sdk-go/site"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration from environment
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Create Echo application
	e := echo.New()

	// Add basic middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	// Create Sitecore client
	sitecoreClient := client.NewSitecoreClient(client.ClientConfig{
		APIEndpoint: cfg.GetGraphQLEndpoint(),
		APIKey:      cfg.GetAPIKey(),
		SiteName:    cfg.DefaultSite,
	})

	// Add healthcheck middleware
	e.Use(middleware.AdaptMiddlewareToEcho(
		middleware.NewHealthcheckMiddleware(middleware.HealthcheckConfig{
			Path: "/healthz",
			Response: map[string]string{
				"status":  "healthy",
				"service": "content-sdk-go",
			},
		}),
	))

	// Add multisite middleware if enabled
	if cfg.Multisite.Enabled {
		fmt.Println("Multisite enabled")
		e.Use(middleware.AdaptMiddlewareToEcho(
			middleware.NewMultisiteMiddleware(middleware.MultisiteConfig{
				DefaultSite:         cfg.DefaultSite,
				UseCookieResolution: cfg.Multisite.UseCookieResolution,
			}),
		))
	}

	// Add locale middleware
	e.Use(middleware.AdaptMiddlewareToEcho(
		middleware.NewLocaleMiddleware(middleware.LocaleConfig{
			DefaultLanguage:   cfg.DefaultLanguage,
			SupportedLanguages: []string{"en", "fr", "de", "es"},
		}),
	))

	// Add redirects middleware if enabled
	if cfg.Multisite.Enabled {
		redirectsService := site.NewRedirectsService(
			cfg.GetGraphQLEndpoint(),
			cfg.GetAPIKey(),
		)
		e.Use(middleware.AdaptMiddlewareToEcho(
			middleware.NewRedirectsMiddleware(middleware.RedirectsConfig{
				RedirectsService: redirectsService,
				CacheDuration:    5 * time.Minute,
			}),
		))
	}

	// SEO Routes
	setupSEORoutes(e, cfg)

	// Editing Routes (if enabled)
	if cfg.Editing.Enabled {
		setupEditingRoutes(e, cfg, sitecoreClient)
	}

	// Catch-all route for Sitecore pages
	catchAllHandler := handlers.NewCatchAllHandler(sitecoreClient)
	e.GET("/*", func(c echo.Context) error {
		ctx := middleware.NewEchoContext(c)
		return catchAllHandler.Handle(ctx)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Starting server on port %s\n", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

// setupSEORoutes configures SEO-related routes
func setupSEORoutes(e *echo.Echo, cfg *config.Config) {
	// Robots.txt
	robotsService := seo.NewRobotsService(
		cfg.GetGraphQLEndpoint(),
		cfg.GetAPIKey(),
	)
	robotsHandler := handlers.NewRobotsHandler(handlers.RobotsHandlerConfig{
		RobotsService: robotsService,
		SitemapURLs:   []string{fmt.Sprintf("https://%s/sitemap.xml", cfg.DefaultSite)},
	})
	e.GET("/robots.txt", func(c echo.Context) error {
		ctx := middleware.NewEchoContext(c)
		return robotsHandler.Handle(ctx)
	})

	// Sitemap.xml
	sitemapService := seo.NewSitemapXmlService(
		cfg.GetGraphQLEndpoint(),
		cfg.GetAPIKey(),
	)
	sitemapHandler := handlers.NewSitemapHandler(handlers.SitemapHandlerConfig{
		SitemapService: sitemapService,
		Sites:          []string{cfg.DefaultSite},
		Languages:      []string{cfg.DefaultLanguage},
	})
	e.GET("/sitemap.xml", func(c echo.Context) error {
		ctx := middleware.NewEchoContext(c)
		return sitemapHandler.Handle(ctx)
	})
}

// setupEditingRoutes configures editing-related routes
func setupEditingRoutes(e *echo.Echo, cfg *config.Config, sitecoreClient *client.SitecoreClient) {
	// Editing config
	editingConfigHandler := handlers.NewEditingConfigHandler(handlers.EditingConfigHandlerConfig{
		SitecoreEdgeURL:       cfg.API.Edge.EdgeURL,
		SitecoreEdgeContextID: cfg.API.Edge.ContextID,
		DefaultLanguage:       cfg.DefaultLanguage,
		DefaultSite:           cfg.DefaultSite,
	})
	e.GET("/api/editing/config", func(c echo.Context) error {
		ctx := middleware.NewEchoContext(c)
		return editingConfigHandler.Handle(ctx)
	})

	// Editing render
	editingRenderHandler := handlers.NewEditingRenderHandler(sitecoreClient)
	e.POST("/api/editing/render", func(c echo.Context) error {
		ctx := middleware.NewEchoContext(c)
		return editingRenderHandler.Handle(ctx)
	})

	// Preview endpoint (for testing)
	e.GET("/api/preview", func(c echo.Context) error {
		itemID := c.QueryParam("itemId")
		language := c.QueryParam("language")
		site := c.QueryParam("site")

		if itemID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "itemId is required",
			})
		}

		// Default values
		if language == "" {
			language = cfg.DefaultLanguage
		}
		if site == "" {
			site = cfg.DefaultSite
		}

		ctx := context.Background()
		layoutService := sitecoreClient.LayoutService

		// Fetch layout data in preview mode
		layoutData, err := layoutService.FetchLayoutData(ctx, "/", site, language)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Error fetching preview: %v", err),
			})
		}

		return c.JSON(http.StatusOK, layoutData)
	})
}

