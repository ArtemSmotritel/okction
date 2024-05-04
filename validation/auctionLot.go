package validation

import (
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
)

type AuctionLotUpdateValidator struct {
	Errors  map[string]string
	Request *types.AuctionLotUpdateRequest
}

func NewAuctionLotUpdateValidator(request *types.AuctionLotUpdateRequest) *AuctionLotUpdateValidator {
	return &AuctionLotUpdateValidator{
		Errors:  make(map[string]string),
		Request: request,
	}
}

func (v *AuctionLotUpdateValidator) Validate() (bool, error) {
	if v.Request.Name == "" {
		v.Errors["name"] = "Auction Lot name is required"
	}

	if v.Request.Description == "" {
		v.Errors["description"] = "Auction Lot description is required"
	}

	if minimalBid, err := utils.StringToDecimal(v.Request.MinimalBidStr); err != nil {
		v.Errors["minimalBid"] = "Minimal Bid price must be a number"
	} else {
		v.Request.MinimalBid = minimalBid
	}

	if reservePrice, err := utils.StringToDecimal(v.Request.ReservePriceStr); err != nil {
		v.Errors["reservePrice"] = "Reserve Price price must be a number"
	} else {
		v.Request.ReservePrice = reservePrice
	}

	if binPrice, err := utils.StringToDecimal(v.Request.BinPriceStr); err != nil {
		v.Errors["binPrice"] = "Bin Price must be a number"
	} else {
		v.Request.BinPrice = binPrice
	}

	return len(v.Errors) == 0, nil
}
