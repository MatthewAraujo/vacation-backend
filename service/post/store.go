package post

import (
	"database/sql"
	"fmt"
	"os"

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
	query := "SELECT p.*, ph.photoID, ph.postID,ph.url_photo, ph.location FROM posts p JOIN photos ph ON p.postID = ph.postID"
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

func (s *Store) GetPostByID(postID string) (*types.Post, error) {
	query := "SELECT p.*, ph.photoID, ph.postID, ph.url_photo, ph.location FROM posts p JOIN photos ph ON p.postID = ph.postID WHERE p.postID = ?"
	rows := s.db.QueryRow(query, postID)
	p, err := scanOneRowIntoPostWithPhotos(rows)
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
	photoLocation := getFileLocation()
	photoURL := "https://www.google.com"

	_, err = s.db.Exec("INSERT INTO photos(photoID,postID,url_photo,location) VALUES (?,?, ?, ?)", uuid.New(), postID, photoURL, photoLocation)
	if err != nil {
		return err
	}
	return nil
}

func scanOneRowIntoPostWithPhotos(rows *sql.Row) (*types.Post, error) {
	p := new(types.Post)
	ph := new(types.Photo)
	err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.Description, &ph.ID, &ph.PostID, &ph.PhotoURL, &ph.Location)
	if err != nil {
		return nil, err
	}
	p.Photos = append(p.Photos, ph)
	return p, nil

}
func scanRowIntoPostWithPhotos(rows *sql.Rows) (*types.Post, error) {
	p := new(types.Post)
	ph := new(types.Photo)
	err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.Description, &ph.ID, &ph.PostID, &ph.PhotoURL, &ph.Location)
	if err != nil {
		return nil, err
	}
	p.Photos = append(p.Photos, ph)
	return p, nil
}

func getFileLocation() string {
	folderPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho atual:", err)
		return ""
	}
	folderPath += "/tmp"

	// Abrir a pasta
	folder, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Erro ao abrir a pasta:", err)
		return ""
	}
	defer folder.Close()

	// Ler o conteúdo da pasta
	files, err := folder.Readdir(1) // Obter apenas um arquivo
	if err != nil {
		fmt.Println("Erro ao ler conteúdo da pasta:", err)
		return ""
	}

	// Verificar se há arquivos na pasta
	if len(files) == 0 {
		fmt.Println("Pasta vazia.")
		return ""
	}

	// Obter o nome do primeiro arquivo
	firstName := files[0].Name()

	return firstName // Add this line to fix the "missing return" error
}
