package gateway

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/example/grpc-user-service/gen/user/v1"
	"github.com/example/grpc-user-service/internal/service"
)

// Handler exposes JSON HTTP endpoints backed by the same gRPC UserServer.
type Handler struct {
	users *service.UserServer
}

func NewHandler(users *service.UserServer) *Handler {
	return &Handler{users: users}
}

func (h *Handler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", h.health)
	mux.HandleFunc("/api/users", h.usersRoot)
	mux.HandleFunc("/api/users/", h.userByID)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) usersRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) userByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id = strings.Trim(id, "/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.users.CreateUser(r.Context(), &userv1.CreateUserRequest{
		Name:  body.Name,
		Email: body.Email,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"rpc":  "CreateUser",
		"user": userJSON(resp.User),
	})
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request, id string) {
	resp, err := h.users.GetUser(r.Context(), &userv1.GetUserRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"rpc":  "GetUser",
		"user": userJSON(resp.User),
	})
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	pageSize := int32(10)
	if v := r.URL.Query().Get("page_size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid page_size")
			return
		}
		pageSize = int32(n)
	}

	resp, err := h.users.ListUsers(r.Context(), &userv1.ListUsersRequest{
		PageSize:  pageSize,
		PageToken: r.URL.Query().Get("page_token"),
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	users := make([]map[string]any, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, userJSON(u))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"rpc":              "ListUsers",
		"users":            users,
		"next_page_token":  resp.NextPageToken,
	})
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request, id string) {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.users.UpdateUser(r.Context(), &userv1.UpdateUserRequest{
		Id:    id,
		Name:  body.Name,
		Email: body.Email,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"rpc":  "UpdateUser",
		"user": userJSON(resp.User),
	})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request, id string) {
	_, err := h.users.DeleteUser(r.Context(), &userv1.DeleteUserRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"rpc":     "DeleteUser",
		"deleted": true,
		"id":      id,
	})
}

func userJSON(u *userv1.User) map[string]any {
	return map[string]any{
		"id":              u.Id,
		"name":            u.Name,
		"email":           u.Email,
		"created_at_unix": u.CreatedAtUnix,
		"updated_at_unix": u.UpdatedAtUnix,
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{
		"error": msg,
		"code":  http.StatusText(code),
	})
}

func writeGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpCode := http.StatusInternalServerError
	switch st.Code() {
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
	case codes.NotFound:
		httpCode = http.StatusNotFound
	case codes.AlreadyExists:
		httpCode = http.StatusConflict
	case codes.Canceled:
		httpCode = 499
	}

	writeJSON(w, httpCode, map[string]any{
		"error":      st.Message(),
		"grpc_code":  st.Code().String(),
		"http_status": httpCode,
	})
}
