package validation

import (
	"github.com/artemsmotritel/oktion/types"
)

type AuctionUpdateValidator struct {
	Errors  map[string]string
	Request types.AuctionUpdateRequest
}

func NewAuctionUpdateValidator(request types.AuctionUpdateRequest) AuctionUpdateValidator {
	return AuctionUpdateValidator{
		Errors:  make(map[string]string),
		Request: request,
	}
}

func (v *AuctionUpdateValidator) Validate() (bool, error) {
	if v.Request.Name == "" {
		v.Errors["name"] = "Auction Name is required"
	}

	if v.Request.Description == "" {
		v.Errors["description"] = "Auction Description is required"
	}

	return len(v.Errors) == 0, nil
}
