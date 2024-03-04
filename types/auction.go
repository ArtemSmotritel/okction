package types

import "time"

type AuctionCreateRequest struct {
	ID int64 `json:"id,omitempty"`
}

type AuctionUpdateRequest struct {
}

type Auction struct {
	ID          int64     `json:"id,omitempty"`
	OwnerId     int64     `json:"ownerId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"isActive,omitempty"`
	IsPrivate   bool      `json:"isPrivate,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
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

func MapAuctionCreateRequest(request AuctionCreateRequest) *Auction {
	return nil
}
