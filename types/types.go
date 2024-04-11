package types

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	PhotoURL  string    `json:"photo_url"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}
type Post struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Description string    `json:"description"`
	Photos      []*Photo  `json:"photos"`
	CreatedAt   time.Time `json:"created_at"`
}
type PostStore interface {
	GetPosts() ([]*Post, error)
	CreatePost(CreatePostPayload) error
}

type CreatePostPayload struct {
	UserID      uuid.UUID      `json:"user_id"`
	Description string         `json:"description" validate:"required,min=3,max=100"`
	Photo       multipart.File `json:"photo" validate:"required"`
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
