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
	rows, err := s.db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}

	posts := make([]*types.Post, 0)
	for rows.Next() {
		p, err := scanRowIntoPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (s *Store) CreatePost(post types.CreatePostPayload) error {
	_, err := s.db.Exec("INSERT INTO posts(postID,userID,description) VALUES (?, ?, ?)", uuid.New(), post.UserID, post.Description)
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
