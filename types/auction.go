package types

import (
	"database/sql"
	"errors"
	"net/url"
	"time"
)

type AuctionCreateRequest struct {
	url.Values
}

func (request *AuctionCreateRequest) name() (string, error) {
	name := request.Get("name")
	if name == "" {
		return "", errors.New("no name was provided for the auction")
	}
	return name, nil
}

func (request *AuctionCreateRequest) description() (string, error) {
	description := request.Get("description")
	if description == "" {
		return "", errors.New("no description was provided for the auction")
	}
	return description, nil
}

func (request *AuctionCreateRequest) private() (bool, error) {
	isPrivate := false
	private := request.Get("private")

	if private == "on" {
		isPrivate = true
	}

	return isPrivate, nil
}

type Auction struct {
	ID          int64        `json:"id,omitempty"`
	OwnerId     int64        `json:"ownerId,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	IsActive    bool         `json:"isActive,omitempty"`
	IsPrivate   bool         `json:"isPrivate,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	DeletedAt   sql.NullTime `json:"-"`
}

func CreateAuction(id int64, ownerId int64, name string, description string, isPrivate bool) *Auction {
	return &Auction{
		ID:          id,
		OwnerId:     ownerId,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func CopyAuction(auction *Auction) Auction {
	newAuction := CreateAuction(auction.ID, auction.OwnerId, auction.Name, auction.Description, auction.IsPrivate)
	newAuction.IsPrivate = auction.IsPrivate
	newAuction.CreatedAt = auction.CreatedAt
	newAuction.UpdatedAt = auction.UpdatedAt
	newAuction.DeletedAt = auction.DeletedAt

	return *newAuction
}

func MapAuctionCreateRequest(values url.Values, ownerId int64) (*Auction, error) {
	request := AuctionCreateRequest{values}

	name, err := request.name()
	if err != nil {
		return nil, err
	}

	description, err := request.description()
	if err != nil {
		return nil, err
	}

	isPrivate, err := request.private()
	if err != nil {
		return nil, err
	}

	auction := &Auction{
		IsActive:    false,
		OwnerId:     ownerId,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
	}

	return auction, nil
}

type AuctionUpdateRequest struct {
	ID          int64
	Name        string
	Description string
	IsPrivate   bool
}

func NewAuctionUpdateRequest(values url.Values, id int64) AuctionUpdateRequest {
	return AuctionUpdateRequest{
		Name:        values.Get("name"),
		Description: values.Get("description"),
		IsPrivate:   values.Get("private") == "on",
		ID:          id,
	}
}

type AuctionFilter struct {
	CategoryId int64
	Name       string
	PerPage    int
	Page       int
	SortBy     string
}

type AuctionFilterBuilder struct {
	filter AuctionFilter
}

func NewAuctionFilterBuilder() *AuctionFilterBuilder {
	return &AuctionFilterBuilder{}
}

func (b *AuctionFilterBuilder) WithCategoryId(id int64) *AuctionFilterBuilder {
	b.filter.CategoryId = id
	return b
}

func (b *AuctionFilterBuilder) WithName(name string) *AuctionFilterBuilder {
	b.filter.Name = name
	return b
}

func (b *AuctionFilterBuilder) WithPerPage(perPage int) *AuctionFilterBuilder {
	b.filter.PerPage = perPage
	return b
}

func (b *AuctionFilterBuilder) WithPage(page int) *AuctionFilterBuilder {
	b.filter.Page = page
	return b
}

func (b *AuctionFilterBuilder) WithSortBy(field string) *AuctionFilterBuilder {
	b.filter.SortBy = field
	return b
}

func (b *AuctionFilterBuilder) Build() AuctionFilter {
	if b.filter.SortBy == "" {
		b.filter.SortBy = "id"
	}
	return b.filter
}
