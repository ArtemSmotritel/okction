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

	mux.HandleFunc("GET /set", setCookie)

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
	mux.HandleFunc("GET /my-auctions", s.handleGetMyAuctions)
	mux.HandleFunc("GET /my-auctions/{id}/edit", s.handleEditAuction)
	mux.HandleFunc("POST /my-auctions/{id}/lots", s.handleCreateAuctionLot)
	mux.HandleFunc("GET /users", s.handleGetUsers)
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", s.handleGetUserByID)
	mux.HandleFunc("PUT /users/{id}", s.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", s.handleDeleteUser)

	mux.HandleFunc("GET /auctions", s.handleGetAuctions)
	mux.HandleFunc("GET /auctions/new", s.handleNewAuction)
	mux.HandleFunc("GET /auctions/{id}", s.handleGetAuctionByID)
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
