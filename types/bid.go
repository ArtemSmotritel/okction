package types

import (
	"github.com/shopspring/decimal"
	"net/url"
	"time"
)

type Bid struct {
	ID           int64
	AuctionLotId int64
	UserId       int64
	CreatedAt    time.Time
	Value        decimal.Decimal
}

type UserBid struct {
	Bid
	IsWonByUser bool
	IsLotActive bool
}

type BidMakeRequest struct {
	AuctionLotId int64
	UserId       int64
	Value        decimal.Decimal
	ValueStr     string
}

func NewBidMakeRequest(values url.Values, auctionLotId int64, userId int64) *BidMakeRequest {
	return &BidMakeRequest{
		ValueStr:     values.Get("value"),
		AuctionLotId: auctionLotId,
		UserId:       userId,
	}
}
