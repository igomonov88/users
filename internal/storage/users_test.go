package storage_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/igomonov88/users/internal/platform/auth"
	"github.com/igomonov88/users/internal/storage"
	"github.com/igomonov88/users/internal/tests"
)

func TestUser(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	ctx := tests.Context()

	t.Log("Given the need to work with users records.")

	u := struct {
		email    string
		name     string
		avatar   string
		password string
	}{
		email:    "gopher@gmail.com",
		name:     "gopher",
		avatar:   "",
		password: "qwerty",
	}

	// Create User Tests
	{
		nu, err := storage.Create(ctx, db, u.email, u.name, u.avatar, u.password)
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

		// Retrieve User Test
		{
			ru, err := storage.Retrieve(ctx, db, nu.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user from storage: %s", tests.Failed, err)
			}

			if cmp.Diff(nu.Name, ru.Name) != "" {
				t.Fatalf("\t%s\tRetrieved user should have same name as was on created state.", tests.Failed)
			}

			if cmp.Diff(nu.Email, ru.Email) != "" {
				t.Fatalf("\t%s\tRetrieved user should have same email as was on created state.", tests.Failed)
			}

			if cmp.Diff(nu.Avatar, ru.Avatar) != "" {
				t.Fatalf("\t%s\tRetrieved user should have same avatar as was on created state.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to retrieve user from storage.", tests.Success)
		}

		// Update User Test
		{
			err := storage.Update(ctx, db, nu.ID, "igor", "myEmail@gmail.com")
			if err != nil {
				t.Fatalf("\t%s\tShould be able to update user : %s", tests.Failed, err)
			}

			ru, err := storage.Retrieve(ctx, db, nu.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user: %s", tests.Failed, err)
			}

			if cmp.Diff(ru.ID, nu.ID) != "" {
				t.Fatalf("\t%s\tRetrieved user should have same user_id after update operation.", tests.Failed)
			}

			if cmp.Diff(ru.Name, "igor") != "" {
				t.Fatalf("\t%s\tShould be able to update user's name.", tests.Failed)
			}

			if cmp.Diff(ru.Email, "myEmail@gmail.com") != "" {
				t.Fatalf("\t%s\tShould be able to update user's email.", tests.Failed)
			}

			t.Logf("\t%s\tShould be able to update user info.", tests.Success)
		}

		// Does Email Exist Test
		{
			exist, err := storage.DoesEmailExist(ctx, db, "myEmail@gmail.com")
			if err != nil {
				t.Fatalf("\t%s\tShould be able to check email exist: %s", tests.Failed, err)
			}

			if !exist {
				t.Fatalf("\t%s\tExisting email is not showing as exist: %s", tests.Failed, err)
			}

			t.Logf("\t%s\tShould be able to check email exist.", tests.Success)
		}

		// Does User Name Exist Test
		{
			exist, err := storage.DoesUserNameExist(ctx, db, "igor")
			if err != nil {
				t.Fatalf("\t%s\tShould be able to check user_name exist: %s", tests.Failed, err)
			}

			if !exist {
				t.Fatalf("\t%s\tExisting user name is not showing as exist: %s", tests.Failed, err)
			}

			t.Logf("\t%s\tShould be able to check user name exist.", tests.Success)
		}

		// Update Avatar Test
		{
			err := storage.UpdateAvatar(ctx, db, nu.ID, "myAvatarURL")
			if err != nil {
				t.Fatalf("\t%s\tShould be able to update avatar: %s ", tests.Failed, err)
			}

			ru, err := storage.Retrieve(ctx, db, nu.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user: %s", tests.Failed, err)
			}

			if cmp.Diff(ru.Avatar, "myAvatarURL") != "" {
				t.Fatalf("t%s\tShould be able to update user's avatar.", tests.Failed)
			}

			t.Logf("\t%s\tShould be able to update user's avatar.", tests.Success)
		}

		// DeleteAvatar User Test
		{
			err := storage.DeleteAvatar(ctx, db, nu.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to delete avatar: %s ", tests.Failed, err)
			}

			ru, err := storage.Retrieve(ctx, db, nu.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user: %s", tests.Failed, err)
			}

			if cmp.Diff(ru.Avatar, "") != "" {
				t.Fatalf("t%s\tShould be able to delete user's avatar.", tests.Failed)
			}

			t.Logf("\t%s\tShould be able to delete user's avatar.", tests.Success)
		}

		// Delete User Test
		{
			err := storage.Delete(ctx, db, nu.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to delete user: %s ", tests.Failed, err)
			}

			_ , err = storage.Retrieve(ctx, db, nu.ID)
			if err != storage.ErrNotFound {
				t.Fatalf("\t%s\tShould be able to delete user: %s ", tests.Failed, err)
			}

			t.Logf("\t%s\tShould be able to delete user.", tests.Success)
		}
	}
}

// TestAuthenticate validates the behavior around authenticating users.
func TestAuthenticate (t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to authenticate users")
	{
		t.Log("\tWhen handling a single User.")
		{
			ctx := tests.Context()

			u := struct {
				email    string
				name     string
				avatar   string
				password string
			}{
				email:    "gopher@gmail.com",
				name:     "gopher",
				avatar:   "",
				password: "qwerty",
			}

			cu, err := storage.Create(ctx, db, u.email,u.name,u.avatar,u.password)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create user: %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create user.", tests.Success)

			claims, err := storage.Authenticate(ctx, db, time.Now(), u.email, u.password)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to generate claims: %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to generate claims.", tests.Success)

			want := auth.Claims{}
			want.Subject = cu.ID
			want.ExpiresAt = time.Now().Add(time.Hour).Unix()
			want.IssuedAt = time.Now().Unix()

			if diff := cmp.Diff(want, claims); diff != "" {
				t.Fatalf("\t%s\tShould get back the expected claims. Diff:\n%s", tests.Failed, diff)
			}
			t.Logf("\t%s\tShould get back the expected claims.", tests.Success)
		}
	}
}