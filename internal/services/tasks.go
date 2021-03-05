package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"togo/internal/driver"
	"togo/internal/storages"

	postgres "togo/internal/storages/pgres"

	"github.com/google/uuid"
)

var dd = driver.Connect()

// Listtask
func ListTasks(resp http.ResponseWriter, req *http.Request) {
	s := &postgres.PostgresDB{}
	s.DB = dd.SQL
	id, _ := userIDFromCtx(req.Context())
	tasks, err := s.RetrieveTasks(
		req.Context(),
		sql.NullString{
			String: id,
			Valid:  true,
		},
		value(req, "created_date"),
	)

	resp.Header().Set("Content-Type", "application/json")

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string][]*storages.Task{
		"data": tasks,
	})

	//fmt.Fprintln(resp, `{"results": [ {"topic": "Clojure", "id": 1000} ]}`)

}

func AddTask(resp http.ResponseWriter, req *http.Request) {
	s := &postgres.PostgresDB{}
	s.DB = dd.SQL
	var str string

	t := &storages.Task{}
	err := json.NewDecoder(req.Body).Decode(t)
	defer req.Body.Close()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	userID, _ := userIDFromCtx(req.Context())
	t.ID = uuid.New().String()
	t.UserID = userID
	t.CreatedDate = now.Format("2006-01-02")

	resp.Header().Set("Content-Type", "application/json")

	str, err = s.AddTask(req.Context(), t)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if str != "" {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": str,
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string]*storages.Task{
		"data": t,
	})
}

func value(req *http.Request, p string) sql.NullString {
	return sql.NullString{
		String: req.FormValue(p),
		Valid:  true,
	}
}

type userAuthKey int8

func userIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(userAuthKey(0))
	fmt.Println("ctx", ctx)
	fmt.Println("v", v)
	fmt.Println("authkey", userAuthKey(0))
	id, ok := v.(string)
	fmt.Println("id", id)
	return id, ok
}
