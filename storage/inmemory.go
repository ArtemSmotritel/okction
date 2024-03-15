package storage

import (
	"fmt"
	"github.com/artemsmotritel/oktion/types"
	"time"
)

type InMemoryStore struct {
	users      []types.User
	auctions   []types.Auction
	categories []types.Category
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

	s.auctions = []types.Auction{{
		ID:          1,
		Name:        "auction1",
		Description: "lorem",
		IsActive:    true,
		IsPrivate:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, {
		ID:          2,
		Name:        "auction2",
		Description: "lorem ipsum",
		IsActive:    true,
		IsPrivate:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}}

	s.categories = []types.Category{{
		ID:   1,
		Name: "Sport",
	},
		{
			ID:   11,
			Name: "Clothes",
		},
	}

	return nil
}

func (s *InMemoryStore) GetAuctionByID(id int64) (*types.Auction, error) {
	for i := 0; i < len(s.auctions); i++ {
		if s.auctions[i].ID == id {
			auction := types.CopyAuction(&s.auctions[i])
			return &auction, nil
		}
	}

	return nil, fmt.Errorf("no auction with id=%d", id)
}

func (s *InMemoryStore) GetAuctions() ([]types.Auction, error) {
	res := make([]types.Auction, len(s.auctions))

	for i := 0; i < len(s.auctions); i++ {
		res[i] = types.CopyAuction(&s.auctions[i])
	}

	return res, nil
}

var id int64 = 0

func (s *InMemoryStore) SaveAuction(auction *types.Auction) (*types.Auction, error) {
	auction.ID = id
	s.auctions = append(s.auctions, types.CopyAuction(auction))
	id++

	return auction, nil
}

func (s *InMemoryStore) DeleteAuction(id int64) error {
	idx := -1

	for i := 0; i < len(s.auctions); i++ {
		if id == s.auctions[i].ID {
			idx = i
		}
	}

	if idx != -1 {
		s.auctions[idx] = s.auctions[len(s.auctions)-1]
		s.auctions = s.auctions[:len(s.auctions)-1]
	}

	return nil
}

func (s *InMemoryStore) GetCategories() ([]types.Category, error) {
	res := make([]types.Category, 0)

	for _, c := range s.categories {
		cc := types.Category{
			ID:   c.ID,
			Name: c.Name,
		}

		res = append(res, cc)
	}

	return res, nil
}
