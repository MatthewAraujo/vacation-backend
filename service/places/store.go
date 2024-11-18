package places

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MatthewAraujo/vacation-backend/types"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	db    *sql.DB
	redis *redis.Client
}

func NewStore(db *sql.DB, redis *redis.Client) *Store {
	return &Store{
		db:    db,
		redis: redis,
	}
}

func (db *Store) getTopTenFromRedis() ([]*types.Post, error) {
	ctx := context.Background()
	places := "places:top-ten"

	postJSON, err := db.redis.Get(ctx, places).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache expired or key not found")
		}
		return nil, fmt.Errorf("failed to get posts from Redis: %w", err)
	}

	var posts []*types.Post
	if err := json.Unmarshal([]byte(postJSON), &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts JSON: %w", err)
	}

	return posts, nil
}

func (db *Store) setTopTenToRedis(posts []*types.Post) error {
	ctx := context.Background()
	places := "places:top-ten"

	postJSON, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("failed to serialize posts: %w", err)
	}

	err = db.redis.Set(ctx, places, postJSON, 24*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set posts in Redis: %w", err)
	}

	return nil
}

func (s *Store) GetTopTenFamous() ([]*types.Post, error) {
	places, err := s.getTopTenFromRedis()
	if err != nil && len(places) < 10 {
		return s.fetchTopTenFromDB()
	}
	return places, nil
}

func (s *Store) fetchTopTenFromDB() ([]*types.Post, error) {
	query := `
		SELECT 
			p.*, 
			ph.photoID, 
			ph.postID, 
			ph.url_photo, 
			l.locationID, 
			l.latitude, 
			l.longitude 
		FROM 
			posts p
		JOIN 
			photos ph ON p.postID = ph.postID
		JOIN 
			locations l ON ph.photoID = l.photoID
		GROUP BY 
			p.postID, ph.photoID, l.locationID
		ORDER BY 
			p.favorite DESC
		LIMIT 10;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*types.Post
	for rows.Next() {
		p, err := scanRowIntoPostWithPhotosAndLocation(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	s.setTopTenToRedis(posts)

	return posts, nil
}

func scanRowIntoPostWithPhotosAndLocation(rows *sql.Rows) (*types.Post, error) {
	p := new(types.Post)
	ph := new(types.Photo)
	l := new(types.LocationBD)

	err := rows.Scan(
		&p.ID,
		&p.UserID,
		&p.CreatedAt,
		&p.Description,
		&p.Favorite,
		&ph.ID,
		&ph.PostID,
		&ph.PhotoURL,
		&l.LocationID,
		&l.Latitude,
		&l.Longitude,
	)
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
