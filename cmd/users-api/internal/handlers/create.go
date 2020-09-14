package handlers

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.User.Create")
	defer span.End()

	txn := u.relict.StartTransaction("create user", w, r)
	defer txn.End()

	cur := CreateUserRequest{}
	if err := web.Decode(r, &cur); err != nil {
		return errors.Wrap(err, "")
	}

	usr, err := storage.Create(ctx, u.db, cur.Email, cur.Name, cur.Avatar, cur.Password)
	if err != nil {
		switch err {
		case storage.ErrEmailAlreadyExist, storage.ErrUserNameAlreadyExist:
			return web.ResponseError(ctx, w, err)
		default:
			return web.ResponseError(ctx, w, err)
		}
	}

	resp := CreateUserResponse{UserID:usr.ID}

	return web.Respond(ctx, w, &resp, http.StatusOK)
}
