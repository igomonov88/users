package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) UpdateAvatar(ctx context.Context, w http.ResponseWriter, r *http.Request,
	params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.UpdateAvatar")
	defer span.End()

	txn := u.relict.StartTransaction("update avatar", w, r)
	defer txn.End()

	req := UpdateAvatarRequest{}
	if err := web.Decode(r, &req); err != nil {
		return web.Respond(ctx, w, web.ErrorResponse{
			Error: "Bad Request",}, http.StatusBadRequest)
	}

	if err := storage.UpdateAvatar(ctx, u.db, req.UserID, req.Avatar); err != nil {
		return web.Respond(ctx, w, web.ErrorResponse{
			Error: "Internal Server Error",}, http.StatusInternalServerError)
	}

	return web.Respond(ctx, w, UpdateAvatarResponse{}, http.StatusOK)
}
