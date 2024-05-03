package storage

import (
	"context"
	"database/sql"
	"github.com/artemsmotritel/oktion/types"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type PostgresqlStore struct {
	connection *pgx.Conn
	logger     *log.Logger
}

func NewPostgresqlStore(conn *pgx.Conn, logger *log.Logger) *PostgresqlStore {
	return &PostgresqlStore{
		connection: conn,
		logger:     logger,
	}
}

func (p *PostgresqlStore) GetUserByID(id int64) (*types.User, error) {
	var user types.User

	query := "SELECT id, email, phone, fullname, password FROM users where id = $1"
	err := p.connection.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Email, &user.Phone, &user.FullName, &user.Password)
	if err != nil {
		// TODO: think of a normal way to log an error
		p.logger.Printf("ERROR: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (p *PostgresqlStore) GetUsers() ([]types.User, error) {
	rows, err := p.connection.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Email, &user.FullName, &user.Phone)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (p *PostgresqlStore) SaveUser(user *types.User) error {
	query := "INSERT INTO users (email, fullname, phone, password) VALUES ($1, $2, $3, $4)"
	args := []any{user.Email, user.FullName, user.Phone, user.Password}
	_, err := p.connection.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser DOES NOT update the user password or email
func (p *PostgresqlStore) UpdateUser(id int64, request types.UserUpdateRequest) (*types.User, error) {
	// intentionally skip email update for now
	query := "UPDATE users SET fullname = $1, phone = $2 WHERE id = $3 RETURNING fullname, email, phone, password"
	args := []any{request.FullName, request.Phone, id}

	var user types.User
	user.ID = id

	err := p.connection.QueryRow(context.Background(), query, args...).Scan(&user.FullName, &user.Email, &user.Phone, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresqlStore) DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := p.connection.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresqlStore) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	query := "SELECT id, email, phone, fullname, password FROM users where email = $1"
	err := p.connection.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Phone, &user.FullName, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresqlStore) GetAuctionsByOwnerId(ownerId int64) ([]types.Auction, error) {
	query := "SELECT id, name, description, is_active, is_private, created_at, updated_at, deleted_at, owner_id FROM auction WHERE owner_id = $1"

	rows, err := p.connection.Query(context.Background(), query, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	auctions := make([]types.Auction, 0)

	for rows.Next() {
		var auction types.Auction
		err := rows.Scan(&auction.ID, &auction.Name, &auction.Description, &auction.IsActive, &auction.IsPrivate, &auction.CreatedAt, &auction.UpdatedAt, &auction.DeletedAt, &auction.OwnerId)
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, auction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return auctions, nil
}

func (p *PostgresqlStore) GetOwnerIDByAuctionID(auctionId int64) (int64, error) {
	query := "SELECT owner_id FROM auction WHERE id = $1"
	var ownerId int64
	err := p.connection.QueryRow(context.Background(), query, auctionId).Scan(&ownerId)
	if err != nil {
		return 0, err
	}

	return ownerId, nil
}

func (p *PostgresqlStore) GetAuctionByID(id int64) (*types.Auction, error) {
	query := "SELECT id, name, description, is_active, is_private, created_at, updated_at, deleted_at, owner_id FROM auction WHERE id = $1"
	var auction types.Auction

	err := p.connection.QueryRow(context.Background(), query, id).Scan(&auction.ID, &auction.Name, &auction.Description, &auction.IsActive, &auction.IsPrivate, &auction.CreatedAt, &auction.UpdatedAt, &auction.DeletedAt, &auction.OwnerId)
	if err != nil {
		return nil, err
	}

	return &auction, nil
}

func (p *PostgresqlStore) GetAuctions() ([]types.Auction, error) {
	return make([]types.Auction, 0), nil
}

func (p *PostgresqlStore) SaveAuction(auction *types.Auction) (*types.Auction, error) {
	query := "INSERT INTO auction (name, description, is_active, is_private, owner_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at"
	args := []any{auction.Name, auction.Description, auction.IsActive, auction.IsPrivate, auction.OwnerId}
	var (
		id        int64
		createdAt time.Time
	)

	err := p.connection.QueryRow(context.Background(), query, args...).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	return &types.Auction{
		ID:          id,
		OwnerId:     auction.OwnerId,
		Name:        auction.Name,
		Description: auction.Description,
		IsActive:    auction.IsActive,
		IsPrivate:   auction.IsPrivate,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		DeletedAt:   sql.NullTime{},
	}, nil
}

func (p *PostgresqlStore) DeleteAuction(id int64) error {
	query := "DELETE FROM auction WHERE id = $1"
	_, err := p.connection.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresqlStore) GetAuctionLotsByAuctionID(auctionId int64) ([]types.AuctionLot, error) {
	query := "SELECT id, name, description, is_active, minimal_bid, reserve_price, bin_price, created_at, updated_at, deleted_at, auction_id FROM auction_lot WHERE auction_id = $1"
	rows, err := p.connection.Query(context.Background(), query, auctionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lots := make([]types.AuctionLot, 0)
	for rows.Next() {
		var lot types.AuctionLot

		if err := rows.Scan(&lot.ID, &lot.Name, &lot.Description, &lot.IsActive, &lot.MinimalBid, &lot.ReservePrice, &lot.BinPrice, &lot.CreatedAt, &lot.UpdatedAt, &lot.DeletedAt, &lot.AuctionID); err != nil {
			return nil, err
		}

		lots = append(lots, lot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lots, err
}

func (p *PostgresqlStore) SaveAuctionLot(auctionLot *types.AuctionLot) (*types.AuctionLot, error) {
	query := "INSERT INTO auction_lot (NAME, DESCRIPTION, IS_ACTIVE, MINIMAL_BID, RESERVE_PRICE, BIN_PRICE, AUCTION_ID) VALUES (@name, @description, @is_active, @minimal_bid, @reserve_price, @bin_price, @auction_id) RETURNING id, created_at"
	args := pgx.NamedArgs{
		"name":          auctionLot.Name,
		"description":   auctionLot.Description,
		"is_active":     auctionLot.IsActive,
		"minimal_bid":   auctionLot.MinimalBid,
		"reserve_price": auctionLot.ReservePrice,
		"bin_price":     auctionLot.BinPrice,
		"auction_id":    auctionLot.AuctionID,
	}

	var (
		id        int64
		createdAt time.Time
	)

	err := p.connection.QueryRow(context.Background(), query, args).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	return &types.AuctionLot{
		ID:           id,
		AuctionID:    auctionLot.AuctionID,
		Name:         auctionLot.Name,
		Description:  auctionLot.Description,
		IsActive:     auctionLot.IsActive,
		MinimalBid:   auctionLot.MinimalBid,
		ReservePrice: auctionLot.ReservePrice,
		BinPrice:     auctionLot.BinPrice,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}, nil
}

func (p *PostgresqlStore) GetAuctionLotCount(auctionId int64) (int, error) {
	query := "SELECT COUNT(id) FROM auction_lot WHERE auction_id = $1"
	var count int
	err := p.connection.QueryRow(context.Background(), query, auctionId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *PostgresqlStore) GetAuctionLotByID(auctionLotId int64) (*types.AuctionLot, error) {
	query := "SELECT name, description, is_active, minimal_bid, reserve_price, bin_price, created_at, updated_at, deleted_at, auction_id FROM auction_lot WHERE id = $1"

	var lot types.AuctionLot

	err := p.connection.QueryRow(context.Background(), query, auctionLotId).Scan(&lot.Name, &lot.Description, &lot.IsActive, &lot.MinimalBid, &lot.ReservePrice, &lot.BinPrice, &lot.CreatedAt, &lot.UpdatedAt, &lot.DeletedAt, &lot.AuctionID)
	if err != nil {
		return nil, err
	}

	return &lot, nil
}

func (p *PostgresqlStore) GetCategories() ([]types.Category, error) {
	query := "SELECT id, name FROM category"

	rows, err := p.connection.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]types.Category, 0)

	for rows.Next() {
		var category types.Category

		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (p *PostgresqlStore) SeedData() error {
	return nil
}
