package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.User.Update")
	defer span.End()

	txn := u.relict.StartTransaction("update user", w, r)
	defer txn.End()

	req := UpdateUserRequest{}
	if err := web.Decode(r, &req); err != nil {
		return web.Respond(ctx, w, nil, http.StatusBadRequest)
	}

	if err := storage.Update(ctx, u.db, req.UserID, req.Name, req.Email); err != nil {
		return web.Respond(ctx, w, nil, http.StatusInternalServerError)
	}

	return web.Respond(ctx, w, UpdateUserResponse{}, http.StatusOK)
}
