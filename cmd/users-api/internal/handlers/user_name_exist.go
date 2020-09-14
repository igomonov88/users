package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/igomonov88/users/internal/platform/web"
	"github.com/igomonov88/users/internal/storage"
)

func (u *User) UserNameExists(ctx context.Context, w http.ResponseWriter, r *http.Request,
	params map[string]string) error {

	ctx, span := trace.StartSpan(ctx, "handlers.User.UserNameExist")
	defer span.End()

	txn := u.relict.StartTransaction("user name exist", w, r)
	defer txn.End()

	un, ok := params["user_name"]
	if !ok {
		return web.Respond(ctx, w, UserNameExistResponse{}, http.StatusBadRequest)
	}

	exist, err := storage.DoesUserNameExist(ctx, u.db, un)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			return web.Respond(ctx, w, web.ErrorResponse{
				Error:"User Not Found",}, http.StatusNotFound)
		default:
			return web.Respond(ctx, w, web.ErrorResponse{
			Error: "Internal Server Error",}, http.StatusInternalServerError)
		}
	}

	resp := UserNameExistResponse{Exist: exist}

	return web.Respond(ctx, w, &resp, http.StatusOK)
}
