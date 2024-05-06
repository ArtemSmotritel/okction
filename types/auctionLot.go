package types

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"net/url"
	"time"
)

type AuctionLot struct {
	ID           int64
	AuctionID    int64
	Name         string
	Description  string
	CategoryId   int64
	IsActive     bool
	IsClosed     bool
	MinimalBid   decimal.Decimal
	ReservePrice decimal.Decimal
	BinPrice     decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

func CopyAuctionLot(auctionLot *AuctionLot) *AuctionLot {
	return &AuctionLot{
		ID:        auctionLot.ID,
		AuctionID: auctionLot.AuctionID,
		Name:      auctionLot.Name,
	}
}

type AuctionLotUpdateRequest struct {
	ID              int64
	AuctionID       int64
	Name            string
	Description     string
	CategoryId      int64
	CategoryIdStr   string
	MinimalBid      decimal.Decimal
	ReservePrice    decimal.Decimal
	BinPrice        decimal.Decimal
	MinimalBidStr   string
	ReservePriceStr string
	BinPriceStr     string
}

func NewAuctionLotUpdateRequest(values url.Values, lotId, auctionId int64) (*AuctionLotUpdateRequest, error) {
	return &AuctionLotUpdateRequest{
		ID:              lotId,
		AuctionID:       auctionId,
		Name:            values.Get("name"),
		Description:     values.Get("description"),
		CategoryIdStr:   values.Get("category"),
		MinimalBidStr:   values.Get("minimalBid"),
		ReservePriceStr: values.Get("reservePrice"),
		BinPriceStr:     values.Get("binPrice"),
	}, nil
}
