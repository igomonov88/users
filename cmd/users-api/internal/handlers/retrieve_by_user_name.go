package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) RetrieveByUserName(ctx context.Context, w http.ResponseWriter, r *http.Request,
	params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.RetrieveByUserName")
	defer span.End()

	txn := u.relict.StartTransaction("retrieve user by name", w, r)
	defer txn.End()

	userName, ok := params["user_name"]
	if !ok {
		return web.Respond(ctx, w, RetrieveUserResponse{}, http.StatusBadRequest)
	}

	usr, err := storage.RetrieveByUserName(ctx, u.db, userName)

	if err != nil {
		switch err {
		case storage.ErrNotFound:
			return web.Respond(ctx, w, web.ErrorResponse{
				Error: "User Not Found",}, http.StatusNotFound)
		default:
			return web.Respond(ctx, w, web.ErrorResponse{
				Error: "Internal Server Error",}, http.StatusInternalServerError)
		}
	}

	resp := RetrieveUserResponse{
		UserID:   usr.ID,
		UserName: usr.Name,
		Email:    usr.Email,
		Avatar:   usr.Avatar,
	}

	return web.Respond(ctx, w, &resp, http.StatusOK)
}
