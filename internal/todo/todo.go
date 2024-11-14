package todo

import (
	"context"
	"first-api/internal/db"
	"strings"
)

type Item struct {
	Title  string
	Status string
}

type Service struct {
	db DBManager
}

type DBManager interface {
	Insert(ctx context.Context, item db.Item) error
	GetAll(ctx context.Context, qStr string) ([]db.Item, error)
}

func NewService(db DBManager) *Service {
	return &Service{
		db: db,
	}
}

func (svc *Service) Add(todo string) error {
	ctx := context.Background()
	todos, err := svc.db.GetAll(ctx, "")

	if err != nil {
		return err
	}
	for _, t := range todos {
		if strings.EqualFold(t.Title, todo) {
			return nil
		}
	}

	err = svc.db.Insert(ctx, db.Item{
		Title:  todo,
		Status: "NOT_STARTED",
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetAll(q string) ([]Item, error) {
	todos, err := svc.db.GetAll(context.Background(), q)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0)
	for _, todo := range todos {
		items = append(items, Item{
			Title:  todo.Title,
			Status: todo.Status,
		})
	}

	return items, nil
}
