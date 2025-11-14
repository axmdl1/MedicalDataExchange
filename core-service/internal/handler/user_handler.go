package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	userv1 "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/user"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type UserHandler struct {
	Client userv1.UserServiceClient
}

func NewUserHandler(cli userv1.UserServiceClient) *UserHandler {
	return &UserHandler{Client: cli}
}

func mdFromRequest(r *http.Request) metadata.MD {
	md := metadata.MD{}
	if v := r.Header.Get("Authorization"); v != "" {
		md.Append("authorization", v)
	}
	return md
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Temporary struct for JSON decoding
	var req struct {
		User struct {
			ID          int64  `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			Type        string `json:"type"`
			Password    string `json:"password"`
			ClinicID    *int64 `json:"clinic_id"` // nullable int64
		} `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to protobuf User
	protoUser := &userv1.User{
		Id:          req.User.ID,
		FirstName:   req.User.FirstName,
		LastName:    req.User.LastName,
		Email:       req.User.Email,
		PhoneNumber: req.User.PhoneNumber,
		Type:        req.User.Type,
		Password:    req.User.Password,
	}
	if req.User.ClinicID != nil {
		protoUser.ClinicId = wrapperspb.Int64(*req.User.ClinicID)
	}

	ctx := metadata.NewOutgoingContext(r.Context(), mdFromRequest(r))
	resp, err := h.Client.CreateUser(ctx, &userv1.CreateUserRequest{User: protoUser})
	if err != nil {
		http.Error(w, err.Error(), httpStatusFromErr(err))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// POST /users/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userv1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Login публичный — без metadata
	resp, err := h.Client.Login(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), httpStatusFromErr(err))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	ctx := metadata.NewOutgoingContext(r.Context(), mdFromRequest(r))
	resp, err := h.Client.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), httpStatusFromErr(err))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GET /users?clinic_id=&type=
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var clinic *wrapperspb.Int64Value
	if s := q.Get("clinic_id"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			clinic = wrapperspb.Int64(v)
		}
	}
	ctx := metadata.NewOutgoingContext(r.Context(), mdFromRequest(r))
	resp, err := h.Client.ListUsers(ctx, &userv1.ListUsersRequest{
		ClinicId: clinic,
		Type:     q.Get("type"),
	})
	if err != nil {
		http.Error(w, err.Error(), httpStatusFromErr(err))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var u userv1.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.Id = id
	ctx := metadata.NewOutgoingContext(r.Context(), mdFromRequest(r))
	resp, err := h.Client.UpdateUser(ctx, &userv1.UpdateUserRequest{User: &u})
	if err != nil {
		http.Error(w, err.Error(), httpStatusFromErr(err))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	ctx := metadata.NewOutgoingContext(r.Context(), mdFromRequest(r))
	resp, err := h.Client.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), httpStatusFromErr(err))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// ---- map gRPC error -> HTTP status ----
func httpStatusFromErr(err error) int {
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.InvalidArgument:
			return http.StatusBadRequest
		case codes.Unauthenticated:
			return http.StatusUnauthorized
		case codes.PermissionDenied:
			return http.StatusForbidden
		case codes.NotFound:
			return http.StatusNotFound
		default:
			return http.StatusBadGateway
		}
	}
	return http.StatusBadGateway
}
