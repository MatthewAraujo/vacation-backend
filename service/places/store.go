package places

import (
	"database/sql"
	"strconv"

	"github.com/MatthewAraujo/vacation-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetTopTenFamous() ([]*types.Post, error) {
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
