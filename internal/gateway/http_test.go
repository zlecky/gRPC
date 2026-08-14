package gateway

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/grpc-user-service/internal/service"
	"github.com/example/grpc-user-service/internal/store"
)

func TestHTTPGatewayCRUD(t *testing.T) {
	h := NewHandler(service.NewUserServer(store.NewMemoryStore()))
	mux := http.NewServeMux()
	h.Mount(mux)

	create := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(`{"name":"Ada","email":"ada@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(create, req)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", create.Code, create.Body.String())
	}
	if !strings.Contains(create.Body.String(), `"rpc":"CreateUser"`) {
		t.Fatalf("unexpected create body: %s", create.Body.String())
	}

	list := httptest.NewRecorder()
	mux.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/api/users?page_size=10", nil))
	if list.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", list.Code, list.Body.String())
	}
	if !strings.Contains(list.Body.String(), "ada@example.com") {
		t.Fatalf("list missing user: %s", list.Body.String())
	}

	id := extractID(create.Body.String())
	if id == "" {
		t.Fatalf("could not parse id from %s", create.Body.String())
	}

	get := httptest.NewRecorder()
	mux.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/users/"+id, nil))
	if get.Code != http.StatusOK {
		t.Fatalf("get status=%d body=%s", get.Code, get.Body.String())
	}

	del := httptest.NewRecorder()
	mux.ServeHTTP(del, httptest.NewRequest(http.MethodDelete, "/api/users/"+id, nil))
	if del.Code != http.StatusOK {
		t.Fatalf("delete status=%d body=%s", del.Code, del.Body.String())
	}

	missing := httptest.NewRecorder()
	mux.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/users/"+id, nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", missing.Code, missing.Body.String())
	}
}

func extractID(body string) string {
	const key = `"id":"`
	i := strings.Index(body, key)
	if i < 0 {
		return ""
	}
	rest := body[i+len(key):]
	j := strings.IndexByte(rest, '"')
	if j < 0 {
		return ""
	}
	return rest[:j]
}
