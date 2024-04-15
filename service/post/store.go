package post

import (
	"database/sql"
	"strconv"

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
	query := "SELECT p.*, ph.photoID, ph.postID, ph.url_photo, l.locationID, l.latitude, l.longitude FROM posts p JOIN photos ph ON p.postID = ph.postID JOIN locations l ON ph.photoID = l.photoID"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	posts := []*types.Post{}

	for rows.Next() {
		p, err := scanRowIntoPostWithPhotosAndLocation(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *Store) GetPostByID(postID string) (*types.Post, error) {
	query := "SELECT p.*, ph.photoID, ph.postID, ph.url_photo, l.locationID, l.latitude, l.longitude FROM posts p JOIN photos ph ON p.postID = ph.postID JOIN locations l ON ph.photoID = l.photoID WHERE p.postID = ?"
	rows := s.db.QueryRow(query, postID)
	p, err := scanOneRowIntoPostWithPhotosAndLocation(rows)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Store) EditPost(postID string, post types.EditPostPayload) (*types.Post, error) {
	_, err := s.db.Exec("UPDATE posts SET description = ? WHERE postID = ?", post.Description, postID)
	if err != nil {
		return nil, err
	}
	return s.GetPostByID(postID)
}

func (s *Store) DeletePost(postID string) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE postID = ?", postID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreatePost(post types.CreatePostPayload) error {
	postID := uuid.New()
	_, err := s.db.Exec("INSERT INTO posts(postID,userID,description) VALUES (?, ?, ?)", postID, post.UserID, post.Description)
	if err != nil {
		return err
	}
	pi, err := GetPhotoInfos()
	if err != nil {
		return err
	}

	photoURL := pi.PhotoURL
	photoLocation := pi.Location

	photoID := uuid.New()
	_, err = s.db.Exec("INSERT INTO photos(photoID,postID,url_photo) VALUES (?,?, ?)", photoID, postID, photoURL)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO locations(locationID, photoID, latitude,longitude) VALUES (?,?,?,?)", uuid.New(), photoID, photoLocation.Latitude, photoLocation.Longitude)
	if err != nil {
		return err
	}

	return nil
}

func scanOneRowIntoPostWithPhotosAndLocation(rows *sql.Row) (*types.Post, error) {
	p := new(types.Post)
	ph := new(types.Photo)
	l := new(types.LocationBD)

	err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.Description, &ph.ID, &ph.PostID, &ph.PhotoURL, &l.LocationID, &l.Latitude, &l.Longitude)
	if err != nil {
		return nil, err
	}
	p.Photos = append(p.Photos, ph)

	latitude, err := strconv.ParseFloat(l.Latitude, 64)
	if err != nil {
		return nil, err
	}

	longitude, err := strconv.ParseFloat(l.Longitude, 64)
	if err != nil {
		return nil, err
	}

	p.Photos[0].Location = types.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return p, nil
}

func scanRowIntoPostWithPhotosAndLocation(rows *sql.Rows) (*types.Post, error) {
	p := new(types.Post)
	ph := new(types.Photo)
	l := new(types.LocationBD)

	err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.Description, &ph.ID, &ph.PostID, &ph.PhotoURL, &l.LocationID, &l.Latitude, &l.Longitude)
	if err != nil {
		return nil, err
	}
	p.Photos = append(p.Photos, ph)

	latitude, err := strconv.ParseFloat(l.Latitude, 64)
	if err != nil {
		return nil, err
	}

	longitude, err := strconv.ParseFloat(l.Longitude, 64)
	if err != nil {
		return nil, err
	}

	p.Photos[0].Location = types.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return p, nil
}
