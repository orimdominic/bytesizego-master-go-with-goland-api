package todo_test

import (
	"context"
	"first-api/internal/db"
	"first-api/internal/todo"
	"reflect"
	"strings"
	"testing"
)

type MockDB struct {
	items []db.Item
}

func (m *MockDB) Insert(ctx context.Context, item db.Item) (_ error) {
	m.items = append(m.items, item)
	return nil
}

func (m *MockDB) GetAll(ctx context.Context, qStr string) (_ []db.Item, _ error) {
	var items []db.Item
	for _, item := range m.items {
		if strings.Contains(strings.ToLower(item.Title), strings.ToLower(qStr)) {
			items = append(items, item)
		}
	}

	return items, nil
}

func TestService_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		q          string
		todosToAdd []string
		want       []todo.Item
	}{
		{
			name:       "given an empty query string, return all items",
			q:          "",
			todosToAdd: []string{"get groceries", "get tools", "study golang"},
			want: []todo.Item{
				{Title: "get groceries", Status: "NOT_STARTED"},
				{Title: "get tools", Status: "NOT_STARTED"},
				{Title: "study golang", Status: "NOT_STARTED"},
			},
		},
		{
			name:       "given a query string, return items containing that string",
			q:          "golang",
			todosToAdd: []string{"get groceries", "get tools", "study golang"},
			want: []todo.Item{
				{Title: "study golang", Status: "NOT_STARTED"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{items: []db.Item{}}
			svc := todo.NewService(m)

			for _, toAdd := range tt.todosToAdd {
				svc.Add(toAdd)
			}

			if got, _ := svc.GetAll(tt.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
