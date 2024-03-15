package api

import (
	"context"
	"fmt"
	"github.com/artemsmotritel/oktion/storage"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
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

	return setUserInfoToContextMiddleware(loggingMiddleware(mux, s.logger))
}

func (s *Server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	hxBoosted, err := utils.ExtractValueFromContext[bool](r.Context(), "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	handler := templates.NewProfilePageHandler(!hxBoosted)
	handler.ServeHTTP(w, r)
}

func getCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("userId")
	return cookie, err
}
func setUserInfoToContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := extractUserIDFromCookie(r)

		if err != nil {
			r = r.WithContext(context.WithValue(r.Context(), "isAuthorized", false))
		} else {
			r = r.WithContext(context.WithValue(r.Context(), "userId", id))
			r = r.WithContext(context.WithValue(r.Context(), "isAuthorized", true))
		}

		r = r.WithContext(context.WithValue(r.Context(), "hxBoosted", r.Header.Get("HX-Boosted") == "true"))

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hxBoosted, _ := utils.ExtractValueFromContext[bool](r.Context(), "hxBoosted")

		l := fmt.Sprintf("New Request: method - %s, url - %s, hx-boosted - %t", r.Method, r.URL.Path, hxBoosted)
		logger.Println(l)

		isAuth := r.Context().Value("isAuthorized").(bool)
		l = fmt.Sprintf("User: isAuthorized - %t", isAuth)
		if isAuth {
			id := r.Context().Value("userId").(int64)
			l += fmt.Sprintf(", id - %d", id)
		}
		logger.Println(l)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewNotFoundPageHandler()
	handler.ServeHTTP(w, r)
}

func (s *Server) badRequestError(w http.ResponseWriter, _ *http.Request, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func (s *Server) internalError(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Something went very wrong at our part...", http.StatusInternalServerError)
}

func (s *Server) handleNewAuction(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewCreateAuctionPageHandler()
	handler.ServeHTTP(w, r)
}

func (s *Server) handleEditAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	auction, err := s.store.GetAuctionByID(id)

	if err != nil {
		s.internalError(w, r)
		return
	}

	if auction == nil {
		s.handleNotFound(w, r)
		return
	}

	handler := templates.NewEditAuctionPageHandler(auction, make([]types.AuctionLot, 0))
	handler.ServeHTTP(w, r)
}

func (s *Server) handleGetMyAuctions(w http.ResponseWriter, r *http.Request) {
	_, err := extractUserIDFromCookie(r)

	if err != nil {
		s.badRequestError(w, r, "not authorized")
		return
	}

	auctions, err := s.store.GetAuctions()

	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleCreateAuctionLot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	auction, err := s.store.GetAuctionByID(id)

	if err != nil {
		s.internalError(w, r)
		return
	}

	if auction == nil {
		s.handleNotFound(w, r)
		return
	}

	handler := templates.NewAuctionLotListItemHandler(id)
	handler.ServeHTTP(w, r)
}

func extractUserIDFromCookie(r *http.Request) (int64, error) {
	cookie, err := getCookie(r)

	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(cookie.Value, 10, 64)
}
