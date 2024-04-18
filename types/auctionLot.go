package types

type AuctionLot struct {
	ID          int64
	AuctionID   int64
	Name        string
	Description string
}

func CopyAuctionLot(auctionLot *AuctionLot) *AuctionLot {
	return &AuctionLot{
		ID:        auctionLot.ID,
		AuctionID: auctionLot.AuctionID,
		Name:      auctionLot.Name,
	}
}
