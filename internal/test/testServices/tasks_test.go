package testServices

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"togo/internal/services"
)

func TestTasks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", services.ListTasks)
	rr := httptest.NewRecorder()
	// body post
	// s := []byte("{\"created_date\": 2021-03-02}")
	// reader := bytes.NewBuffer(s)

	req, _ := http.NewRequest("GET", "/tasks?created_date=2021-03-02", nil)
	ctx := req.Context()
	//var UserAuthKey int8
	ctx = context.Background()
	//ctx = context.WithValue(ctx, UserAuthKey, "firstUser")
	req = req.WithContext(ctx)
	mux.ServeHTTP(rr, req)
	//http.DefaultServeMux.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())
	if rr.Code != 200 {
		t.Errorf("Response code is %v", rr.Code)
	}
}
