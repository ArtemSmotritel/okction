package types

import (
	"database/sql"
	decimal "github.com/jackc/pgx-shopspring-decimal"
	"time"
)

type AuctionLot struct {
	ID           int64
	AuctionID    int64
	Name         string
	Description  string
	IsActive     bool
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
