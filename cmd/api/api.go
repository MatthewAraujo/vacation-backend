package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MatthewAraujo/vacation-backend/service/places"
	"github.com/MatthewAraujo/vacation-backend/service/post"
	"github.com/MatthewAraujo/vacation-backend/service/user"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type APIServer struct {
	addr  string
	db    *sql.DB
	redis *redis.Client
}

func NewAPIServer(addr string, db *sql.DB, redis *redis.Client) *APIServer {
	return &APIServer{
		addr:  addr,
		db:    db,
		redis: redis,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	// if the api changes in the future we can just change the version here, and the old version will still be available
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	postStore := post.NewStore(s.db)
	postHandler := post.NewHandler(postStore, userStore)
	postHandler.RegisterRoutes(subrouter)

	placesStore := places.NewStore(s.db, s.redis)
	placesHandler := places.NewHandler(placesStore)
	placesHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
