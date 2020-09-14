package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.User.Delete")
	defer span.End()

	txn := u.relict.StartTransaction("delete user", w, r)
	defer txn.End()

	req := DeleteUserRequest{}
	if err := web.Decode(r, &req); err != nil {
		return web.Respond(ctx, w, DeleteUserResponse{}, http.StatusBadRequest)
	}

	if err := storage.Delete(ctx, u.db, req.UserID); err!= nil {
		return web.Respond(ctx, w, DeleteUserResponse{}, http.StatusInternalServerError)
	}

	return web.Respond(ctx, w, DeleteUserResponse{}, http.StatusOK)
}
