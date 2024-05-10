package validation

import (
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"github.com/shopspring/decimal"
)

type BidMakeValidator struct {
	Request     *types.BidMakeRequest
	Errors      map[string]string
	lotProvider AuctionLotProvider
}

type AuctionLotProvider interface {
	GetAuctionLotByID(auctionLotId int64) (*types.AuctionLot, error)
}

type SavedAuctionLotProvider struct {
	Lot *types.AuctionLot
}

func (s *SavedAuctionLotProvider) GetAuctionLotByID(auctionLotId int64) (*types.AuctionLot, error) {
	return s.Lot, nil
}

func NewBidMakeValidator(request *types.BidMakeRequest, provider AuctionLotProvider) *BidMakeValidator {
	return &BidMakeValidator{
		Request:     request,
		Errors:      make(map[string]string),
		lotProvider: provider,
	}
}

func (v *BidMakeValidator) Validate() (bool, error) {
	lot, err := v.lotProvider.GetAuctionLotByID(v.Request.AuctionLotId)
	if err != nil {
		return false, err
	}

	if lot.IsClosed {
		v.Errors["value"] = "This auction is closed already"
	} else if !lot.IsActive {
		v.Errors["value"] = "This auction is archived"
	}

	if len(v.Errors) != 0 {
		return false, nil
	}

	if v.Request.ValueStr == "" {
		v.Errors["value"] = "Bid Value is required"
	} else if value, err := utils.StringToDecimal(v.Request.ValueStr); err != nil {
		v.Errors["value"] = "Bid Value must be a number"
	} else if value.Compare(decimal.Zero) < 0 {
		v.Errors["value"] = "Bid Value can't be less than zero"
	} else {
		v.Request.Value = value

		if v.Request.Value.Cmp(lot.MinimalBid) < 0 {
			v.Errors["value"] = "Bid Value can't be less than Lot Minimal Bid"
		}
	}

	return len(v.Errors) == 0, nil
}
