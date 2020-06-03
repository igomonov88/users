package storage_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/igomonov88/users/internal/storage"
	"github.com/igomonov88/users/internal/tests"
)


func TestUsers(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	ctx := tests.Context()

	t.Log("Given the need to work with users records.")

	u := struct {
		email string
		name string
		avatar string
		password string
	}{
		email:"gopher@gmail.com",
		name:"gopher",
		avatar: "",
		password: "qwerty",
	}

	// Create User Tests
	{
		nu, err := storage.Create(ctx, db,u.email, u.name, u.avatar, u.password)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to add new user to storage: %s", tests.Failed, err)
		}

		if cmp.Diff(nu.Name, u.name) != "" {
			t.Fatalf("\t%s\tCreated user should have same name as was provided.", tests.Failed)
		}

		if cmp.Diff(nu.Email, u.email) != "" {
			t.Fatalf("\t%s\tCreated user should have same email as was provided", tests.Failed)
		}

		if cmp.Diff(nu.Avatar, u.avatar) != "" {
			t.Fatalf("\t%s\tCreated user should have same avatar as was provided", tests.Failed)
		}
		t.Logf("\t%s\tShould be able to add new user to storage.", tests.Success)
	}


}
