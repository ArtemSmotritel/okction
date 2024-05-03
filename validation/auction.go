package validation

import (
	"errors"
	"github.com/artemsmotritel/oktion/types"
)

type AuctionUpdateValidator struct {
	Errors  map[string]error
	Request types.AuctionUpdateRequest
}

func NewAuctionUpdateValidator(request types.AuctionUpdateRequest) AuctionUpdateValidator {
	return AuctionUpdateValidator{
		Errors:  make(map[string]error),
		Request: request,
	}
}

func (v *AuctionUpdateValidator) Validate() (bool, error) {
	if v.Request.Name == "" {
		v.Errors["Name"] = errors.New("auction name is required")
	}

	if v.Request.Description == "" {
		v.Errors["Description"] = errors.New("auction description is required")
	}

	return len(v.Errors) == 0, nil
}
