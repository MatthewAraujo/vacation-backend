package post

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/MatthewAraujo/vacation-backend/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.PostStore
	userStore types.UserStore
}

func NewHandler(store types.PostStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/posts", h.handleGetPosts).Methods(http.MethodGet)
	// router.HandleFunc("/posts", auth.WithJWTAuth(h.handleCreatePost, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/posts", h.handleCreatePost).Methods(http.MethodPost)
}

func (h *Handler) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.GetPosts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, posts)
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var post types.CreatePostPayload

	_, err := utils.ParseMultipartForm(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	post = types.CreatePostPayload{
		Description: r.FormValue("description"),
		UserID:      uuid.MustParse(r.FormValue("user_id")),
	}

	// validation of the payload
	if err := utils.Validate.Struct(&post); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	user, err := h.userStore.GetUserByID(post.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}

	_, err = FileUploadHandler(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	post.UserID = user.ID
	err = h.store.CreatePost(post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, post)
}

func FileUploadHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("photo")
	log.Println(header.Filename)

	if err != nil {
		fmt.Fprintln(w, err)
		return "", err
	}

	defer file.Close()

	folderPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho atual:", err)
		return "", err
	}

	folderPath += "/tmp/"

	path := folderPath + header.Filename
	out, err := os.Create(path)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	return path, nil
}
