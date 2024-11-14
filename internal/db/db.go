package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

type Item struct {
	Title     string
	Status    string
	Id        int
	CreatedAt time.Time
}

func New(
	user, password, host, dbName string, port int,
) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Insert(ctx context.Context, item Item) error {
	query := `INSERT INTO todos (title, status) VALUES ($1, $2);`

	_, err := db.pool.Exec(ctx, query, item.Title, "NOT_STARTED")
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAll(ctx context.Context, qStr string) ([]Item, error) {
	var rows pgx.Rows
	var err error

	query := `SELECT * FROM todos`
	if qStr != "" {
		query += ` WHERE title ILIKE @qStr`
		rows, err = db.pool.Query(ctx, query, pgx.NamedArgs{"qStr": fmt.Sprintf("%%%s%%", qStr)})
	} else {
		rows, err = db.pool.Query(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item

		err := rows.Scan(&item.Id, &item.Status, &item.Title, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return items, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
