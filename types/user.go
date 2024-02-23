package types

type UserCreateRequest struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type UserUpdateRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type User struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func CreateUser(id int64, firstName, lastName string) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func CopyUser(user *User) User {
	return *CreateUser(user.ID, user.FirstName, user.LastName)
}

func MapUserCreateRequest(request UserCreateRequest) *User {
	return CreateUser(request.ID, request.FirstName, request.LastName)
}
