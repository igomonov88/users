package handlers

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

// Token handles a request to authenticate a user. It expects a request using
// Basic Auth with a user's email and password. It responds with a JWT.
func (u *User) Token(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.User.Token")
	defer span.End()

	txn := u.relict.StartTransaction("get token", w, r)
	defer txn.End()



	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := storage.Authenticate(ctx, u.db, v.Now, email, pass)
	if err != nil {
		switch err {
		case storage.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}
	tkn := TokenResponse{}

	tkn.Token, err = u.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
