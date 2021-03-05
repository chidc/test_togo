package main

import (
	"net/http"

	"togo/internal/routers"

	_ "github.com/lib/pq"
)

func main() {
	http.ListenAndServe(":5050", &routers.ToDoService{})
}
