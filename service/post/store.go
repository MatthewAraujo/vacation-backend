package post

import (
	"database/sql"

	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetPosts() ([]*types.Post, error) {
	query := "SELECT p.*, ph.url_photo, ph.location FROM posts p JOIN photos ph ON p.postID = ph.postID"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	posts := []*types.Post{}

	for rows.Next() {
		p, err := scanRowIntoPostWithPhotos(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *Store) CreatePost(post types.CreatePostPayload) error {
	postID := uuid.New()
	_, err := s.db.Exec("INSERT INTO posts(postID,userID,description) VALUES (?, ?, ?)", postID, post.UserID, post.Description)
	if err != nil {
		return err
	}
	photoLocation := "ucrain"
	photoURL := "https://www.google.com"

	_, err = s.db.Exec("INSERT INTO photos(photoID,postID,url_photo,location) VALUES (?,?, ?, ?)", uuid.New(), postID, photoURL, photoLocation)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoPost(rows *sql.Rows) (*types.Post, error) {
	p := new(types.Post)
	err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.Description)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func scanRowIntoPostWithPhotos(rows *sql.Rows) (*types.Post, error) {
	p := new(types.Post)
	ph := new(types.Photo)
	err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.Description, &ph.PhotoURL, &ph.Location)
	if err != nil {
		return nil, err
	}
	p.Photos = append(p.Photos, ph)
	return p, nil
}
