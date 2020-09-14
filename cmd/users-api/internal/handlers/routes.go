package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	newrelic "github.com/newrelic/go-agent"

	"github.com/igomonov88/users/internal/mid"
	"github.com/igomonov88/users/internal/platform/auth"
	"github.com/igomonov88/users/internal/platform/web"
)

// User  represents the user API method handler set.
type User struct {
	db            *sqlx.DB
	authenticator *auth.Authenticator
	relict        newrelic.Application
}

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB,
	relic newrelic.Application, authenticator *auth.Authenticator) http.Handler {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		build: build,
		db:    db,
	}

	// Register user check endpoint.
	u := User{
		db:            db,
		authenticator: authenticator,
		relict:        relic,
	}
	app.Handle("GET", "/v1/health", check.Health)

	// This route is not authenticated
	app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	app.Handle(http.MethodPost, "/v1/users", u.Create)
	app.Handle(http.MethodGet, "/v1/users/email/:email", u.EmailExist)
	app.Handle(http.MethodGet, "/v1/users/user_name/:user_name", u.UserNameExists)

	app.Handle(http.MethodPost, "/v1/users/update", u.Update, mid.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/users/delete", u.Delete, mid.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/users/:user_id", u.Retrieve, mid.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/users/by_email/:email", u.RetrieveByEmail, mid.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/users/by_user_name/:user_name", u.RetrieveByUserName, mid.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/users/update_avatar", u.UpdateAvatar, mid.Authenticate(authenticator))

	return app
}
