package types

import "net/url"

type User struct {
	ID       int64  `json:"id,omitempty"`
	FullName string `json:"firstName,omitempty"`
	Password string `json:"-"`
	Phone    string `json:"phone,omitempty"`
	Email    string
}

type UserUpdateRequest struct {
	ID       int64
	FullName string
	Email    string
	Phone    string
}

func NewUserUpdateRequest(values url.Values, id int64) UserUpdateRequest {
	return UserUpdateRequest{
		ID:       id,
		FullName: values.Get("fullName"),
		Email:    values.Get("email"),
		Phone:    values.Get("phone"),
	}
}

func CreateUser(id int64, fullName string, email string, password string) *User {
	return &User{
		ID:       id,
		FullName: fullName,
		Email:    email,
		Password: password,
	}
}

func CopyUser(user *User) *User {
	return CreateUser(user.ID, user.FullName, user.Email, user.Password)
}
