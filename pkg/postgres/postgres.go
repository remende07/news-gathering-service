package postgres

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

// Конструктор объекта БД
func New(connstr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	s := Store{
		db: db,
	}

	return &s, nil
}

// Загружает в БД посты
func (s *Store) CreatePosts(posts []models.Post) error {
	for _, p := range posts {
		_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (title, content, published_at, link)
		VALUES ($1, $2, $3, $4)`,
			p.Title, p.Content, p.PubTime, p.Link)

		if err != nil {
			return err
		}
	}

	return nil
}

// Возвращает последние n постов
func (s *Store) Posts(n int) ([]models.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
		id,
		title,
		content,
		published_at,
		link
		FROM posts
		ORDER BY published_at DESC
		LIMIT $1`,
		n)

	if err != nil {
		return nil, err
	}

	posts := []models.Post{}

	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}
