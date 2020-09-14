package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) EmailExist(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.User.EmailExist")
	defer span.End()

	txn := u.relict.StartTransaction("email exist", w, r)
	defer txn.End()

	email, ok := params["email"]
	if !ok {
		return web.Respond(ctx, w, EmailExistResponse{}, http.StatusBadRequest)
	}

	exist, err := storage.DoesEmailExist(ctx, u.db, email)
	if err != nil {
		return web.Respond(ctx, w, web.ErrorResponse{
			Error: "Internal Server Error",}, http.StatusInternalServerError)
	}

	resp := EmailExistResponse{Exist: exist}

	return web.Respond(ctx, w, &resp, http.StatusOK)
}
