package transport

import (
	"encoding/json"
	"first-api/internal/todo"
	"fmt"
	"net/http"
)

type TodoItem struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

type TodoItems struct {
	Items []todo.Item `json:"docs"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(todoSvc *todo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		todos, err := todoSvc.GetAll(q)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&TodoItems{
			Items: todos,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})

	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var t TodoItem

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		err = todoSvc.Add(t.Title)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(&TodoItem{
			Title:  t.Title,
			Status: "NOT_STARTED",
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})

	return &Server{
		mux: mux,
	}
}

func (server *Server) Serve() error {
	fmt.Println("server running on 8000")

	err := http.ListenAndServe(":8000", server.mux)
	if err != nil {
		return err
	}

	return nil
}
