package api

import (
	"fmt"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"github.com/artemsmotritel/oktion/validation"
	"net/http"
	"strconv"
	"strings"
)

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

	auctionLots, err := s.store.GetAuctionLotsByAuctionID(auction.ID)
	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewEditAuctionPageHandler(auction, auctionLots)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleGetMyAuctions(w http.ResponseWriter, r *http.Request) {
	_, err := extractUserIDFromCookie(r)

	if err != nil {
		s.badRequestError(w, r, "not authorized")
		return
	}

	ownerId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		s.handleUnauthorized(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(ownerId)

	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleGetAuctions(w http.ResponseWriter, r *http.Request) {
	categoryIdOrNameStr := r.URL.Query().Get("category")
	name := r.URL.Query().Get("name")
	pageStr := r.URL.Query().Get("page")

	filterBuilder := types.NewAuctionFilterBuilder()
	filterBuilder = filterBuilder.SetPerPage(10)
	if name != "" {
		filterBuilder = filterBuilder.SetName("%" + strings.ToLower(name) + "%")
	}
	if val, err := strconv.ParseInt(categoryIdOrNameStr, 10, 64); err == nil {
		filterBuilder = filterBuilder.SetCategoryId(val)
	}
	if val, err := strconv.Atoi(pageStr); err == nil {
		filterBuilder = filterBuilder.SetPage(val)
	}

	filter := filterBuilder.Build()

	auctions, err := s.store.GetAuctions(filter)
	if err != nil {
		s.internalError(w, r)
		return
	}

	totalAuctionCount, categoryName, err := s.store.CountAuctionsAndGetCategoryName(filter)
	if err != nil {
		s.internalError(w, r)
		return
	}

	pageParamBuilder := types.NewAuctionsListPageParameterBuilder()
	pageParamBuilder = pageParamBuilder.SetAuctions(auctions)
	pageParamBuilder = pageParamBuilder.SetAuctionsFound(totalAuctionCount)
	pageParamBuilder = pageParamBuilder.SetFilter(filter)
	if categoryName.Valid {
		pageParamBuilder = pageParamBuilder.SetCategoryName(categoryName.String)
	}

	pageParam := pageParamBuilder.Build()

	urlToPushInBrowserUrl := "/auctions?"
	if filter.Name != "" {
		urlToPushInBrowserUrl += "name=" + name
	}
	if filter.CategoryId != 0 {
		urlToPushInBrowserUrl += "category=" + strconv.FormatInt(filter.CategoryId, 10)
	}

	w.Header().Set("HX-Push-Url", urlToPushInBrowserUrl)
	handler := templates.NewAuctionsListPageHandler(pageParam, r.Context())
	handler.ServeHTTP(w, r)
}

func (s *Server) handleGetAuctionView(w http.ResponseWriter, r *http.Request) {
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

	lots, err := s.store.GetAuctionLotsByAuctionID(auction.ID)
	if err != nil {
		s.internalError(w, r)
		return
	}

	user, err := s.store.GetUserByID(auction.OwnerId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	isAuth, err := utils.ExtractValueFromContext[bool](r.Context(), "isAuthorized")
	if err != nil {
		isAuth = false
	}

	pageParam := types.AuctionViewPageParam{
		Auction:      auction,
		Lots:         lots,
		Owner:        user,
		IsAuthorized: isAuth,
	}

	handler := templates.NewAuctionViewHandler(&pageParam, r.Context())
	handler.ServeHTTP(w, r)
}

func (s *Server) handleCreateAuction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.badRequestError(w, r, "Couldn't parse form")
		return
	}

	ownerID, err := extractUserIDFromCookie(r)

	if err != nil {
		s.badRequestError(w, r, "Invalid cookie")
	}

	auction, err := types.MapAuctionCreateRequest(r.Form, ownerID)

	if err != nil {
		s.badRequestError(w, r, "Bad form request: "+err.Error())
		return
	}

	savedAuction, err := s.store.SaveAuction(auction)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Add("HX-Push-Url", fmt.Sprintf("/my-auctions/%d/edit", savedAuction.ID))
	handler := templates.NewEditAuctionPageHandler(savedAuction, []types.AuctionLot{})
	w.WriteHeader(http.StatusCreated)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleDeleteAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = s.store.DeleteAuction(id); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpdateAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = r.ParseForm(); err != nil {
		s.badRequestError(w, r, err.Error())
		return
	}

	updateRequest := types.NewAuctionUpdateRequest(r.Form, id)
	validator := validation.NewAuctionUpdateValidator(updateRequest)
	ok, err := validator.Validate()
	if err != nil {
		s.internalError(w, r)
		return
	}

	if !ok {
		auctionWithBadData := types.Auction{
			ID:          id,
			Name:        updateRequest.Name,
			Description: updateRequest.Description,
			IsActive:    true,
			IsPrivate:   updateRequest.IsPrivate,
		}
		w.Header().Set("HX-Retarget", "#create-auction-form-1")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.Header().Set("HX-Replace-Url", fmt.Sprintf("/my-auctions/%s/edit", utils.IdToString(id)))
		handler := templates.NewAuctionEditFormErrorBadRequestHandler(&auctionWithBadData, validator.Errors)
		handler.ServeHTTP(w, r)
		return
	}

	updatedAuction, err := s.store.UpdateAuction(updateRequest)
	if err != nil {
		s.internalError(w, r)
		return
	}

	auctionLots, err := s.store.GetAuctionLotsByAuctionID(updatedAuction.ID)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/my-auctions/%s/edit", utils.IdToString(id)))
	w.WriteHeader(http.StatusCreated)
	handler := templates.NewEditAuctionPageHandler(updatedAuction, auctionLots)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleArchiveAuction(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = s.store.SetAuctionActiveStatus(id, false); err != nil {
		s.internalError(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(userId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleReinstateAuction(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = s.store.SetAuctionActiveStatus(id, true); err != nil {
		s.internalError(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(userId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleCloseAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
		return
	}

	if err = s.store.SetWinnersToAllAuctionLots(id); err != nil {
		s.internalError(w, r)
		return
	}

	if err = s.store.CloseAuction(id); err != nil {
		s.internalError(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(userId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}
