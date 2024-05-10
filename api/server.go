package api

import (
	"github.com/artemsmotritel/oktion/storage"
	"github.com/artemsmotritel/oktion/templates"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type Server struct {
	listenAddress string
	store         storage.Storage
	logger        *log.Logger
}

func NewServer(listenAddress string, store storage.Storage, logger *log.Logger) *Server {
	return &Server{
		listenAddress: listenAddress,
		store:         store,
		logger:        logger,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.listenAddress, s.newConfiguredRouter())
}

func (s *Server) newConfiguredRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	homePaths := []string{"/", "/home"}
	categories, _ := s.store.GetCategories()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)

		if slices.Contains[[]string](homePaths, path) {
			handler := templates.NewIndexPageHandler(categories)
			handler.ServeHTTP(w, r)
			return
		}

		s.handleNotFound(w, r)
	})
	mux.HandleFunc("GET /profile", s.handleGetProfile)
	mux.Handle("GET /my-auctions", s.onlyAuthorizedMiddleware(http.HandlerFunc(s.handleGetMyAuctions)))
	mux.Handle("GET /my-auctions/{id}/edit", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleEditAuction), "id"))
	mux.Handle("POST /my-auctions/{id}/lots", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleCreateAuctionLot), "id"))
	mux.Handle("GET /my-auctions/{auctionId}/lots/{lotId}/edit", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleEditAuctionLot), "auctionId"))

	mux.Handle("GET /login", templates.NewLoginPageHandler())
	mux.HandleFunc("POST /login", s.handleLogin)
	mux.Handle("GET /sign-up", templates.NewSignUpPageHandler())
	mux.HandleFunc("POST /sign-up", s.handleSignUp)
	mux.HandleFunc("POST /logout", s.handleLogout)

	mux.HandleFunc("GET /users", s.handleGetUsers)
	mux.HandleFunc("GET /users/{id}", s.handleGetUserByID)
	mux.HandleFunc("PUT /users/{id}", s.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", s.handleDeleteUser)

	mux.HandleFunc("GET /auctions", s.handleGetAuctions)
	mux.HandleFunc("GET /auctions/new", s.handleNewAuction)
	mux.HandleFunc("GET /auctions/{id}/view", s.handleGetAuctionView)

	mux.Handle("PUT /auctions/{id}", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleUpdateAuction), "id"))
	mux.Handle("POST /auctions/{id}/archive", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleArchiveAuction), "id"))
	mux.Handle("POST /auctions/{id}/reinstate", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleReinstateAuction), "id"))
	mux.Handle("PUT /auctions/{auctionId}/lots/{lotId}", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleUpdateAuctionLot), "auctionId"))
	mux.Handle("POST /auctions/{auctionId}/lots/{lotId}/archive", s.protectAuctionsMiddleware(s.handleSetAuctionLotActiveStatus(false), "auctionId"))
	mux.Handle("POST /auctions/{auctionId}/lots/{lotId}/reinstate", s.protectAuctionsMiddleware(s.handleSetAuctionLotActiveStatus(true), "auctionId"))
	mux.HandleFunc("GET /auctions/{auctionId}/lots/{lotId}/view", s.handleViewAuctionLot)

	// TODO: finish these and make UI for them
	mux.HandleFunc("GET /categories/{id}/auctions", s.handleGetCategoryAuctions)
	mux.Handle("POST /auctions/{auctionId}/lots/{lotId}/make-favorite", s.onlyAuthorizedMiddleware(s.handleSetUserFavoriteAuctionLot(true)))
	mux.Handle("POST /auctions/{auctionId}/lots/{lotId}/unmake-favorite", s.onlyAuthorizedMiddleware(s.handleSetUserFavoriteAuctionLot(false)))
	mux.Handle("POST /auctions/{auctionId}/lots/{lotId}/bid", s.onlyNotAuctionOwnerMiddleware(s.onlyOpenAuction(http.HandlerFunc(s.handleMakeBid), "auctionId"), "auctionId"))
	mux.Handle("GET /my-bids", s.onlyAuthorizedMiddleware(http.HandlerFunc(s.handleGetMyBids)))
	mux.Handle("POST /auctions/{id}/close", s.protectAuctionsMiddleware(http.HandlerFunc(s.handleCloseAuction), "id"))

	mux.HandleFunc("POST /auctions", s.handleCreateAuction)
	mux.HandleFunc("DELETE /auctions/{id}", s.handleDeleteAuction)

	return setUserInfoToContextMiddleware(loggingMiddleware(redirectUserMiddleware(mux), s.logger))
}

func extractUserIDFromCookie(r *http.Request) (int64, error) {
	cookie, err := getCookie(r)

	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(cookie.Value, 10, 64)
}
