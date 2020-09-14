package handlers

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Avatar   string `json:"avatar"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type DeleteUserRequest struct {
	UserID string `json:"user_id"`
}

type DeleteUserResponse struct {
}

type EmailExistRequest struct{}

type EmailExistResponse struct {
	Exist bool `json:"exist"`
}

type RetrieveUserRequest struct {
}

type RetrieveUserResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UpdateAvatarRequest struct {
	UserID string `json:"user_id"`
	Avatar string `json:"avatar"`
}

type UpdateAvatarResponse struct{}

type UpdateUserRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type UpdateUserResponse struct{}

type UserNameExistRequest struct{}

type UserNameExistResponse struct {
	Exist bool `json:"exist"`
}
