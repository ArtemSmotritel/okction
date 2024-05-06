package validation

import (
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"github.com/shopspring/decimal"
	"strconv"
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
	} else if comp := minimalBid.Compare(decimal.Zero); comp < 0 {
		v.Errors["minimalBid"] = "Minimal Bid must be no less than zero"
	} else {
		v.Request.MinimalBid = minimalBid
	}

	if reservePrice, err := utils.StringToDecimal(v.Request.ReservePriceStr); err != nil {
		v.Errors["reservePrice"] = "Reserve Price price must be a number"
	} else if comp := reservePrice.Compare(decimal.Zero); comp < 0 {
		v.Errors["reservePrice"] = "Reserve Price must be no less than zero"
	} else {
		v.Request.ReservePrice = reservePrice
	}

	if binPrice, err := utils.StringToDecimal(v.Request.BinPriceStr); err != nil {
		v.Errors["binPrice"] = "Bin Price must be a number"
	} else if comp := binPrice.Compare(decimal.Zero); comp < 0 {
		v.Errors["binPrice"] = "Bin Price must be no less than zero"
	} else {
		v.Request.BinPrice = binPrice
	}

	if v.Request.CategoryIdStr == "" {
		v.Errors["category"] = "Category is required"
	} else if categoryId, err := strconv.ParseInt(v.Request.CategoryIdStr, 10, 64); err != nil {
		v.Errors["category"] = "Category must have a valid value"
	} else {
		v.Request.CategoryId = categoryId
	}

	return len(v.Errors) == 0, nil
}
