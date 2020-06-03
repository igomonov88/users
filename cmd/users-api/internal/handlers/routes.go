package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/igomonov88/users/internal/mid"
	"github.com/igomonov88/users/internal/platform/auth"
	"github.com/igomonov88/users/internal/platform/cache"
	"github.com/igomonov88/users/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, fdcClient *api.Client, c *cache.Cache) http.Handler {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		build: build,
		db:    db,
	}

	app.Handle("GET", "/v1/health", check.Health)

	return app
}
