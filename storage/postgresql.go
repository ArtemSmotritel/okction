package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"github.com/jackc/pgx/v5"
	"log"
	"strconv"
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

func (p *PostgresqlStore) logError(err error, tag string) {
	p.logger.Printf("An error occurred when executing a query to the postgres db\nTAG: %s\nERROR: %+v\n", tag, err)
}

func (p *PostgresqlStore) GetUserByID(id int64) (*types.User, error) {
	var user types.User

	query := "SELECT id, email, phone, fullname, password FROM users where id = $1"
	err := p.connection.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Email, &user.Phone, &user.FullName, &user.Password)
	if err != nil {
		// TODO: think of a normal way to log an error
		p.logError(err, "get user by id")
		return nil, err
	}

	return &user, nil
}

func (p *PostgresqlStore) GetUsers() ([]types.User, error) {
	rows, err := p.connection.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		p.logError(err, "get users")
		return nil, err
	}
	defer rows.Close()

	var users []types.User

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Email, &user.FullName, &user.Phone)
		if err != nil {
			p.logError(err, "get users; rows")
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		p.logError(err, "get users; after rows")
		return nil, err
	}

	return users, nil
}

func (p *PostgresqlStore) SaveUser(user *types.User) (*types.User, error) {
	query := "INSERT INTO users (email, fullname, phone, password) VALUES ($1, $2, $3, $4) RETURNING id, email, fullname, phone"
	args := []any{user.Email, user.FullName, user.Phone, user.Password}

	var savedUser types.User
	returningArgs := []any{&savedUser.ID, &savedUser.Email, &savedUser.FullName, &savedUser.Phone}
	err := p.connection.QueryRow(context.Background(), query, args...).Scan(returningArgs...)
	if err != nil {
		p.logError(err, "save user")
		return nil, err
	}

	return &savedUser, nil
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
		p.logError(err, "update user")
		return nil, err
	}

	return &user, nil
}

func (p *PostgresqlStore) DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := p.connection.Exec(context.Background(), query, id)
	if err != nil {
		p.logError(err, "delete user")
		return err
	}

	return nil
}

func (p *PostgresqlStore) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	query := "SELECT id, email, phone, fullname, password FROM users where email = $1"
	err := p.connection.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Phone, &user.FullName, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		p.logError(err, "get user by email")
		return nil, err
	}

	return &user, nil
}

func (p *PostgresqlStore) GetAuctionsByOwnerId(ownerId int64) ([]types.Auction, error) {
	query := "SELECT id, name, description, is_active, is_private, created_at, updated_at, deleted_at, owner_id FROM auction WHERE owner_id = $1"

	rows, err := p.connection.Query(context.Background(), query, ownerId)
	if err != nil {
		p.logError(err, "get auctions by owner id")
		return nil, err
	}
	defer rows.Close()
	auctions := make([]types.Auction, 0)

	for rows.Next() {
		var auction types.Auction
		err := rows.Scan(&auction.ID, &auction.Name, &auction.Description, &auction.IsActive, &auction.IsPrivate, &auction.CreatedAt, &auction.UpdatedAt, &auction.DeletedAt, &auction.OwnerId)
		if err != nil {
			p.logError(err, "get auctions by owner id; rows")
			return nil, err
		}

		auctions = append(auctions, auction)
	}

	if err = rows.Err(); err != nil {
		p.logError(err, "get auctions by owner id; after rows")
		return nil, err
	}

	return auctions, nil
}

func (p *PostgresqlStore) GetOwnerIDByAuctionID(auctionId int64) (int64, error) {
	query := "SELECT owner_id FROM auction WHERE id = $1"
	var ownerId int64
	err := p.connection.QueryRow(context.Background(), query, auctionId).Scan(&ownerId)
	if err != nil {
		p.logError(err, "get auction owner id by auction id")
		return 0, err
	}

	return ownerId, nil
}

func (p *PostgresqlStore) GetAuctionByID(id int64) (*types.Auction, error) {
	query := "SELECT id, name, description, is_active, is_private, created_at, updated_at, deleted_at, owner_id FROM auction WHERE id = $1"
	var auction types.Auction

	err := p.connection.QueryRow(context.Background(), query, id).Scan(&auction.ID, &auction.Name, &auction.Description, &auction.IsActive, &auction.IsPrivate, &auction.CreatedAt, &auction.UpdatedAt, &auction.DeletedAt, &auction.OwnerId)
	if err != nil {
		p.logError(err, "get auction by id")
		return nil, err
	}

	return &auction, nil
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
		p.logError(err, "save auction")
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
		p.logError(err, "delete auction")
		return err
	}

	return nil
}

func (p *PostgresqlStore) GetAuctionLotsByAuctionID(auctionId int64) ([]types.AuctionLot, error) {
	query := "SELECT id, name, description, is_active, minimal_bid, reserve_price, bin_price, created_at, updated_at, deleted_at, auction_id, COALESCE((SELECT category_id FROM auction_lot_categories WHERE auction_lot_id = $1), 0), is_closed FROM auction_lot WHERE auction_id = $1"
	rows, err := p.connection.Query(context.Background(), query, auctionId)
	if err != nil {
		p.logError(err, "get auction lots by auction id")
		return nil, err
	}
	defer rows.Close()

	lots := make([]types.AuctionLot, 0)
	for rows.Next() {
		var lot types.AuctionLot

		if err := rows.Scan(&lot.ID, &lot.Name, &lot.Description, &lot.IsActive, &lot.MinimalBid, &lot.ReservePrice, &lot.BinPrice, &lot.CreatedAt, &lot.UpdatedAt, &lot.DeletedAt, &lot.AuctionID, &lot.CategoryId, &lot.IsClosed); err != nil {
			p.logError(err, "get auction lots by auction id; rows")
			return nil, err
		}

		lots = append(lots, lot)
	}

	if err = rows.Err(); err != nil {
		p.logError(err, "get auction lots by auction id; after rows")
		return nil, err
	}

	return lots, err
}

func (p *PostgresqlStore) SaveAuctionLot(auctionLot *types.AuctionLot) (*types.AuctionLot, error) {
	query := "INSERT INTO auction_lot (NAME, DESCRIPTION, IS_ACTIVE, MINIMAL_BID, RESERVE_PRICE, BIN_PRICE, AUCTION_ID, IS_CLOSED) VALUES (@name, @description, @is_active, @minimal_bid, @reserve_price, @bin_price, @auction_id, @is_closed) RETURNING id, created_at"
	args := pgx.NamedArgs{
		"name":          auctionLot.Name,
		"description":   auctionLot.Description,
		"is_active":     auctionLot.IsActive,
		"minimal_bid":   auctionLot.MinimalBid,
		"reserve_price": auctionLot.ReservePrice,
		"bin_price":     auctionLot.BinPrice,
		"auction_id":    auctionLot.AuctionID,
		"is_closed":     auctionLot.IsClosed,
	}

	var (
		id        int64
		createdAt time.Time
	)

	err := p.connection.QueryRow(context.Background(), query, args).Scan(&id, &createdAt)
	if err != nil {
		p.logError(err, "save auction lot")
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
		p.logError(err, "get auction lot count")
		return 0, err
	}

	return count, nil
}

func (p *PostgresqlStore) GetAuctionLotByID(auctionLotId int64) (*types.AuctionLot, error) {
	query := "SELECT name, description, is_active, minimal_bid, reserve_price, bin_price, created_at, updated_at, deleted_at, auction_id, COALESCE((SELECT category_id FROM auction_lot_categories WHERE auction_lot_id = $1), 0), is_closed FROM auction_lot WHERE id = $1"

	var lot types.AuctionLot
	lot.ID = auctionLotId

	returningArgs := []any{&lot.Name, &lot.Description, &lot.IsActive, &lot.MinimalBid, &lot.ReservePrice, &lot.BinPrice, &lot.CreatedAt, &lot.UpdatedAt, &lot.DeletedAt, &lot.AuctionID, &lot.CategoryId, &lot.IsClosed}

	err := p.connection.QueryRow(context.Background(), query, auctionLotId).Scan(returningArgs...)
	if err != nil {
		p.logError(err, "get auction lot by id")
		return nil, err
	}

	return &lot, nil
}

func (p *PostgresqlStore) GetCategories() ([]types.Category, error) {
	query := "SELECT id, name FROM category"

	rows, err := p.connection.Query(context.Background(), query)
	if err != nil {
		p.logError(err, "get categories")
		return nil, err
	}
	defer rows.Close()

	categories := make([]types.Category, 0)

	for rows.Next() {
		var category types.Category

		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			p.logError(err, "get categories; rows")
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		p.logError(err, "get categories; after rows")
		return nil, err
	}

	return categories, nil
}

func (p *PostgresqlStore) SeedData() error {
	return nil
}

func (p *PostgresqlStore) UpdateAuction(update types.AuctionUpdateRequest) (*types.Auction, error) {
	query := "UPDATE auction SET name = @name, description = @description, is_private = @is_private, updated_at = @updated_at WHERE id = @id RETURNING name, description, is_private, is_active, updated_at, created_at, deleted_at, owner_id"
	args := pgx.NamedArgs{
		"name":        update.Name,
		"description": update.Description,
		"is_private":  update.IsPrivate,
		"updated_at":  time.Now(),
		"id":          update.ID,
	}

	var auction types.Auction
	auction.ID = update.ID

	returningArgs := []any{&auction.Name, &auction.Description, &auction.IsPrivate, &auction.IsActive, &auction.UpdatedAt, &auction.CreatedAt, &auction.DeletedAt, &auction.OwnerId}

	if err := p.connection.QueryRow(context.Background(), query, args).Scan(returningArgs...); err != nil {
		p.logError(err, "update auction")
		return nil, err
	}

	return &auction, nil
}

func (p *PostgresqlStore) SetAuctionActiveStatus(id int64, isActive bool) error {
	query := "UPDATE auction SET is_active = $1 WHERE id = $2"
	if _, err := p.connection.Exec(context.Background(), query, isActive, id); err != nil {
		p.logError(err, "set auction active status")
		return err
	}

	return nil
}

func (p *PostgresqlStore) UpdateAuctionLot(auctionLotId int64, request *types.AuctionLotUpdateRequest) (*types.AuctionLot, error) {
	updateLotCategorySubQuery := "WITH category_subquery AS (INSERT INTO auction_lot_categories (auction_lot_id, category_id) VALUES (@id, @category_id) ON CONFLICT (auction_lot_id) DO UPDATE SET category_id = @category_id RETURNING category_id), "
	updateLotQuery := "lot_subquery AS (UPDATE auction_lot SET name = @name, description = @description, minimal_bid = @minimal_bid, reserve_price = @reserve_price, bin_price = @bin_price, updated_at = @updated_at WHERE id = @id RETURNING name, description, is_active, minimal_bid, reserve_price, bin_price, updated_at, created_at, is_closed) "
	selectQuery := "SELECT * FROM category_subquery, lot_subquery"

	query := updateLotCategorySubQuery + updateLotQuery + selectQuery
	args := pgx.NamedArgs{
		"id":            auctionLotId,
		"name":          request.Name,
		"description":   request.Description,
		"minimal_bid":   request.MinimalBid,
		"reserve_price": request.ReservePrice,
		"bin_price":     request.BinPrice,
		"updated_at":    time.Now(),
		"category_id":   request.CategoryId,
	}

	var lot types.AuctionLot
	lot.ID = auctionLotId

	returningArgs := []any{&lot.CategoryId, &lot.Name, &lot.Description, &lot.IsActive, &lot.MinimalBid, &lot.ReservePrice, &lot.BinPrice, &lot.UpdatedAt, &lot.CreatedAt, &lot.IsClosed}

	if err := p.connection.QueryRow(context.Background(), query, args).Scan(returningArgs...); err != nil {
		p.logError(err, "update auction lot")
		return nil, err
	}

	return &lot, nil
}

func (p *PostgresqlStore) SetAuctionLotActiveStatus(auctionLotId int64, isActive bool) error {
	query := "UPDATE auction_lot SET is_active = $1 WHERE id = $2"

	if _, err := p.connection.Exec(context.Background(), query, isActive, auctionLotId); err != nil {
		p.logError(err, "set auction lot active status")
		return err
	}

	return nil
}

func (p *PostgresqlStore) CountAuctionsAndGetCategoryName(filter types.AuctionFilter) (int, sql.NullString, error) {
	query := `SELECT DISTINCT ON (a.id) a.id, cat.name FROM auction a
			LEFT JOIN auction_lot l ON l.auction_id = a.id
			LEFT JOIN auction_lot_categories c ON c.auction_lot_id = l.id
            LEFT JOIN category cat ON cat.id = c.category_id
		WHERE a.is_active
			AND (@category_id = 0 OR c.category_id = @category_id)
			AND (@name = '' OR LOWER(a.name) LIKE @name)
			AND (@show_deleted OR a.deleted_at IS NULL)
		GROUP BY a.id, cat.name
	`

	args := pgx.NamedArgs{
		"name":         filter.Name,
		"category_id":  filter.CategoryId,
		"order_by":     filter.SortBy,
		"show_deleted": filter.ShowDeleted,
	}

	var (
		count        int
		categoryName sql.NullString
	)

	if err := p.connection.QueryRow(context.Background(), query, args).Scan(&count, &categoryName); err != nil {
		p.logError(err, "count auctionsAndGetCategoryName")
	}

	return count, categoryName, nil
}

func (p *PostgresqlStore) GetAuctions(filter types.AuctionFilter) ([]types.Auction, error) {
	query := `SELECT DISTINCT ON (a.id) a.id, a.name, a.description, a.is_active, a.is_private, a.created_at, a.updated_at, a.deleted_at, a.owner_id FROM auction a
			LEFT JOIN auction_lot l ON l.auction_id = a.id
			LEFT JOIN auction_lot_categories c ON c.auction_lot_id = l.id
		WHERE a.is_active
			AND (@category_id = 0 OR c.category_id = @category_id)
			AND (@name = '' OR LOWER(a.name) LIKE @name)
			AND (@show_deleted OR a.deleted_at IS NULL)
		OFFSET @offset LIMIT @limit
	`

	args := pgx.NamedArgs{
		"name":         filter.Name,
		"category_id":  filter.CategoryId,
		"order_by":     filter.SortBy,
		"limit":        filter.PerPage,
		"offset":       filter.Page * filter.PerPage,
		"show_deleted": filter.ShowDeleted,
	}

	rows, err := p.connection.Query(context.Background(), query, args)
	if err != nil {
		p.logError(err, "get auctions")
	}
	defer rows.Close()

	auctions := make([]types.Auction, 0)

	for rows.Next() {
		auction := types.Auction{}
		returningArgs := []any{&auction.ID, &auction.Name, &auction.Description, &auction.IsActive, &auction.IsPrivate, &auction.CreatedAt, &auction.UpdatedAt, &auction.DeletedAt, &auction.OwnerId}

		if err = rows.Scan(returningArgs...); err != nil {
			p.logError(err, "get auctions; rows")
			return nil, err
		}

		auctions = append(auctions, auction)
	}

	if err = rows.Err(); err != nil {
		p.logError(err, "get auctions; after rows")
		return nil, err
	}

	return auctions, nil
}

func (p *PostgresqlStore) SetUserFavoriteAuctionLot(userId, auctionLotId int64, isFavorite bool) error {
	var query string
	if isFavorite {
		query = "INSERT INTO saved_auction_lots (auction_lot_id, user_id) VALUES (@auction_lot_id, @user_id) ON CONFLICT (user_id, auction_lot_id) DO NOTHING"
	} else {
		query = "DElETE FROM saved_auction_lots WHERE user_id = @user_id AND auction_lot_id = @auction_lot_id"
	}

	args := pgx.NamedArgs{
		"auction_lot_id": auctionLotId,
		"user_id":        userId,
	}

	_, err := p.connection.Exec(context.Background(), query, args)
	if err != nil {
		p.logError(err, "set user favorite auction lot: "+strconv.FormatBool(isFavorite))
		return err
	}

	return nil
}

func (p *PostgresqlStore) SaveAuctionLotBid(request *types.BidMakeRequest) (*types.Bid, error) {
	query := "INSERT INTO bid (value, user_id, auction_lot_id) VALUES (@value, @user_id, @auction_lot_id) RETURNING id, value, user_id, auction_lot_id, created_at"
	args := pgx.NamedArgs{
		"value":          request.Value,
		"user_id":        request.UserId,
		"auction_lot_id": request.AuctionLotId,
	}

	var bid types.Bid
	returningArgs := []any{&bid.ID, &bid.Value, &bid.UserId, &bid.AuctionLotId, &bid.CreatedAt}

	if err := p.connection.QueryRow(context.Background(), query, args).Scan(returningArgs...); err != nil {
		p.logError(err, "save auction lot bid")
		return nil, err
	}

	return &bid, nil
}

func (p *PostgresqlStore) GetUserBids(userId int64) ([]types.UserBid, error) {
	query := `SELECT DISTINCT ON (b.id) b.id, b.value, b.user_id, b.auction_lot_id, b.created_at, l.is_active, w.user_id = $1 AS did FROM bid b
		LEFT JOIN auction_lot l on l.id = b.auction_lot_id
		LEFT JOIN auction_lot_winner w on w.auction_lot_id = l.id
		WHERE b.user_id = $1
		AND (l.is_active = true OR w.user_id = $1)`
	//SELECT DISTINCT ON (b.id) b.id, b.value, b.auction_lot_id, l.is_active, COALESCE((w.bid_id = b.id), false) AS did_win_lot FROM bid b
	//LEFT JOIN auction_lot l on l.id = b.auction_lot_id
	//LEFT JOIN auction_lot_winner w on w.bid_id = b.id
	//WHERE b.user_id = 5;

	rows, err := p.connection.Query(context.Background(), query, userId)
	if err != nil {
		p.logError(err, "get user bids")
		return nil, err
	}
	defer rows.Close()

	bids := make([]types.UserBid, 0)
	for rows.Next() {
		bid := types.UserBid{
			Bid: types.Bid{},
		}
		returningArgs := []any{&bid.ID, &bid.Value, &bid.UserId, &bid.AuctionLotId, &bid.CreatedAt, &bid.IsLotActive, &bid.IsWonByUser}

		if err = rows.Scan(returningArgs...); err != nil {
			p.logError(err, "get user bids; rows")
			return nil, err
		}

		bids = append(bids, bid)
	}

	if err = rows.Err(); err != nil {
		p.logError(err, "get user bids; after rows")
		return nil, err
	}

	return bids, err
}

func (p *PostgresqlStore) CloseAuction(auctionId int64) error {
	query := "WITH r AS (UPDATE auction_lot SET is_closed = true WHERE auction_id = $1) UPDATE auction SET is_closed = true WHERE id = $1"

	// TODO: maybe this query should be executed in a transaction
	if _, err := p.connection.Exec(context.Background(), query, auctionId); err != nil {
		p.logError(err, "close auction")
		return err
	}

	return nil
}

func (p *PostgresqlStore) CheckAuctionStatus(auctionId int64) (utils.Status, error) {
	query := "SELECT is_closed, is_active FROM auction WHERE id = $1"

	status := utils.Status{}

	if err := p.connection.QueryRow(context.Background(), query, auctionId).Scan(&status.IsClosed, &status.IsActive); err != nil {
		p.logError(err, "is auction closed")
		return status, err
	}

	return status, nil
}

func (p *PostgresqlStore) CheckAuctionLotStatus(lotId int64) (utils.Status, error) {
	query := "SELECT l.is_closed = true OR a.is_closed = true AS closed, l.is_active = true AND a.is_active = TRUE AS active FROM auction_lot l LEFT JOIN auction a on a.id = l.auction_id WHERE l.id = $1"

	status := utils.Status{}

	if err := p.connection.QueryRow(context.Background(), query, lotId).Scan(&status.IsClosed, &status.IsActive); err != nil {
		p.logError(err, "is auction closed")
		return status, err
	}

	return status, nil
}
