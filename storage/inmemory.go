package storage

import (
	"fmt"
	"github.com/artemsmotritel/oktion/types"
)

type InMemoryStore struct {
	users []types.User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

func (s *InMemoryStore) GetUserByID(id int64) (*types.User, error) {
	for i := 0; i < len(s.users); i++ {
		if id == s.users[i].ID {
			u := types.CopyUser(&s.users[i])
			return &u, nil
		}
	}

	return nil, nil
}

func (s *InMemoryStore) GetUsers() ([]types.User, error) {
	res := make([]types.User, len(s.users))

	for i := 0; i < len(s.users); i++ {
		res[i] = types.CopyUser(&s.users[i])
	}

	return res, nil
}

func (s *InMemoryStore) SaveUser(user *types.User) error {
	s.users = append(s.users, types.CopyUser(user))
	return nil
}

func (s *InMemoryStore) UpdateUser(id int64, request *types.UserUpdateRequest) error {
	for i := 0; i < len(s.users); i++ {
		if id == s.users[i].ID {
			s.users[i].FirstName = request.FirstName
			s.users[i].LastName = request.LastName
			return nil
		}
	}

	return fmt.Errorf("no user with id=%d", id)
}

func (s *InMemoryStore) DeleteUser(id int64) error {
	idx := -1

	for i := 0; i < len(s.users); i++ {
		if id == s.users[i].ID {
			idx = i
		}
	}

	if idx != -1 {
		s.users[idx] = s.users[len(s.users)-1]
		s.users = s.users[:len(s.users)-1]
	}
	return nil
}

func (s *InMemoryStore) SeedData() error {
	s.users = []types.User{{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
	}, {
		ID:        2,
		FirstName: "Jane",
		LastName:  "Doe",
	}, {
		ID:        3,
		FirstName: "Abobus",
		LastName:  "Sus",
	},
	}
	return nil
}
