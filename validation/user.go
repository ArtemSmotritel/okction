package validation

import (
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/artemsmotritel/oktion/types"
	"net/mail"
	"net/url"
	"strings"
)

type SignUpValidator struct {
	Email           string
	Password        string
	ConfirmPassword string
	FullName        string
	Errors          map[string]string
}

type LoginValidator struct {
	Email    string
	Password string
	Errors   map[string]string
}

func NewLoginValidator() *LoginValidator {
	return &LoginValidator{}
}

func NewSignUpValidator() *SignUpValidator {
	return &SignUpValidator{}
}

type UserIdentityProvider interface {
	GetUserByEmail(string) (*types.User, error)
}

func (u *LoginValidator) Validate(values url.Values, identityProvider UserIdentityProvider) (bool, error) {
	if identityProvider == nil {
		return false, errors.New("user identity provider is nil")
	}

	u.Errors = make(map[string]string)
	u.Email = values.Get("email")
	u.Password = values.Get("password")

	var user *types.User
	var err error

	if valid, message := validateEmail(u.Email); !valid {
		u.Errors["email"] = message
	} else {
		user, err = identityProvider.GetUserByEmail(u.Email)
		if err != nil {
			return false, err
		}
	}

	if user == nil {
		u.Errors["email"] = "Invalid email or password"
		return false, nil
	}

	if u.Password == "" {
		u.Errors["password"] = "Enter a password"
	} else {
		isSame, err := argon2id.ComparePasswordAndHash(u.Password, user.Password)
		if err != nil {
			return false, err
		}
		if !isSame {
			u.Errors["password"] = "Invalid email or password"
		}
	}

	return len(u.Errors) == 0, nil
}

func (u *LoginValidator) Values() map[string]string {
	return map[string]string{
		"email":    u.Email,
		"password": u.Password,
	}
}

func (u *SignUpValidator) Values() map[string]string {
	return map[string]string{
		"email":            u.Email,
		"password":         u.Password,
		"confirm-password": u.ConfirmPassword,
	}
}

func (u *SignUpValidator) Validate(values url.Values, identityProvider UserIdentityProvider) (bool, error) {
	if identityProvider == nil {
		return false, errors.New("user identity provider is nil")
	}
	u.Errors = make(map[string]string)
	u.Email = values.Get("email")
	u.Password = values.Get("password")

	var user *types.User
	var err error

	if valid, message := validateEmail(u.Email); !valid {
		u.Errors["email"] = message
	} else {
		user, err = identityProvider.GetUserByEmail(u.Email)
		if err != nil {
			return false, err
		}
	}

	if user != nil {
		u.Errors["email"] = "This email is already taken"
	}

	if u.Password == "" {
		u.Errors["password"] = "Enter a password"
	}

	if u.Password != "" && u.Password != values.Get("confirm-password") {
		u.Errors["confirm-password"] = "Your passwords don't match"
	}

	return len(u.Errors) == 0, nil
}

var userId int64

func MapUserCreateRequestToUser(request *SignUpValidator) *types.User {
	userId++
	return &types.User{
		ID:       userId,
		FullName: request.FullName,
		Email:    request.Email,
		Password: request.Password,
	}
}

func validateEmail(email string) (isValid bool, validationMessage string) {
	e := strings.TrimSpace(email)
	isValid = true

	if e == "" {
		isValid = false
		validationMessage = "Enter your email"
	} else if !IsEmailValid(e) {
		isValid = false
		validationMessage = "Enter a valid email"
	}

	return
}

func IsEmailValid(email string) bool {
	// TODO implement a correct validation so that "abobus correct@email.com" would not be correct
	_, err := mail.ParseAddress(email)
	return err == nil
}
