package types

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	PhotoURL  string    `json:"photo_url"`
	Location  Location  `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type LocationBD struct {
	LocationID string
	PhotoID    string
	Latitude   string
	Longitude  string
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type PhotoInfo struct {
	PhotoURL string
	Location Location
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"description"`
	Photos      []*Photo  `json:"photos"`
	Favorite    int       `json:"favorite"`
	CreatedAt   time.Time `json:"created_at"`
}
type PostStore interface {
	GetPosts() ([]*Post, error)
	CreatePost(CreatePostPayload) error
	GetPostByID(postID string) (*Post, error)
	EditPost(postID string, post EditPostPayload) (*Post, error)
	DeletePost(postID string) error
}

type CreatePostPayload struct {
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"description" validate:"required,min=3,max=100"`
	// Photo       multipart.FileHeader
	// `json:"photo" validate:"required"`
}

type EditPostPayload struct {
	Description string `json:"description" validate:"required,min=3,max=100"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	CreateUser(u *User) error
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type PlacesStore interface {
	GetTopTenFamous() ([]*Post, error)
}
