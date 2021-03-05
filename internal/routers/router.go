package routers

import (
	"log"
	"net/http"

	"togo/internal/services"
)

// ToDoService implement HTTP server
type ToDoService struct {
	// JWTKey string
	// Store  *postgres.PostgresDB
}

func (s *ToDoService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	switch req.URL.Path {
	case "/login":
		services.GetAuthToken(resp, req)
		return
	case "/tasks":
		var ok bool
		req, ok = services.ValidToken(req)
		if !ok {
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch req.Method {
		case http.MethodGet:
			services.ListTasks(resp, req)
		case http.MethodPost:
			services.AddTask(resp, req)
		}
		return
	}
}
