package storage

import (
	"database/sql"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
)

type Storage interface {
	GetUserByID(id int64) (*types.User, error)
	GetUsers() ([]types.User, error)
	SaveUser(user *types.User) (*types.User, error)
	UpdateUser(id int64, request types.UserUpdateRequest) (*types.User, error)
	DeleteUser(id int64) error
	GetUserByEmail(email string) (*types.User, error)

	GetAuctionsByOwnerId(ownerId int64) ([]types.Auction, error)
	GetOwnerIDByAuctionID(auctionId int64) (int64, error)
	GetAuctionByID(id int64) (*types.Auction, error)
	GetAuctions(filter types.AuctionFilter) ([]types.Auction, error)
	CountAuctionsAndGetCategoryName(filter types.AuctionFilter) (int, sql.NullString, error)
	SaveAuction(auction *types.Auction) (*types.Auction, error)
	DeleteAuction(id int64) error
	UpdateAuction(auction types.AuctionUpdateRequest) (*types.Auction, error)
	SetAuctionActiveStatus(auctionId int64, isActive bool) error
	CloseAuction(auctionId int64) error
	CheckAuctionStatus(auctionId int64) (utils.Status, error)

	GetAuctionLotsByAuctionID(auctionId int64) ([]types.AuctionLot, error)
	SaveAuctionLot(auctionLot *types.AuctionLot) (*types.AuctionLot, error)
	GetAuctionLotCount(auctionId int64) (int, error)
	GetAuctionLotByID(auctionLotId int64) (*types.AuctionLot, error)
	UpdateAuctionLot(auctionLotId int64, lot *types.AuctionLotUpdateRequest) (*types.AuctionLot, error)
	SetAuctionLotActiveStatus(auctionLotId int64, isActive bool) error
	SetUserFavoriteAuctionLot(userId, auctionLotId int64, isFavorite bool) error
	CheckAuctionLotStatus(lotId int64) (utils.Status, error)
	GetSavedAuctionLots(userId int64) ([]types.AuctionLot, error)
	DoesUserSavedAuctionLot(userId, auctionLotId int64) (bool, error)
	CanBidOnAuctionLot(auctionLotId int64) (bool, error)

	SaveAuctionLotBid(request *types.BidMakeRequest) (*types.Bid, error)
	GetUserBids(userId int64) ([]types.UserBid, error)

	GetCategories() ([]types.Category, error)

	SeedData() error
}
