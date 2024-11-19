package places

import (
	"net/http"

	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/MatthewAraujo/vacation-backend/utils"
	"github.com/MatthewAraujo/vacation-backend/xcsf"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.PlacesStore
}

func NewHandler(store types.PlacesStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/places/top-ten-famous", xcsf.WithCSF(h.handleGetTopTenFamous)).Methods(http.MethodGet)
}

func (h *Handler) handleGetTopTenFamous(w http.ResponseWriter, r *http.Request) {
	places, err := h.store.GetTopTenFamous()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	utils.WriteJSON(w, http.StatusOK, places)
	utils.Logger(http.StatusOK, r)
}
