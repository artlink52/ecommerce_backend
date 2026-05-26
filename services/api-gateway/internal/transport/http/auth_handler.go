package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"github.com/artlink52/ecommerce_backend/services/api-gateway/internal/transport/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserClient interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	userClient UserClient
}

func NewAuthHandler(userClient UserClient) *AuthHandler {
	return &AuthHandler{userClient: userClient}

}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request body", http.StatusBadRequest)
		return
	}

	log.Info("register request", slog.String("email", req.Email))

	userID, err := h.userClient.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Warn("register failed", slog.String("email", req.Email), logger.Err(err))
		writeGRPCError(responseHandler, err)
		return
	}

	log.Info("user registered", slog.String("email", req.Email), slog.Int64("user_id", userID))
	responseHandler.JSONResponse(map[string]any{"user_id": userID}, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request body", http.StatusBadRequest)
		return
	}

	log.Info("login request", slog.String("email", req.Email))

	token, err := h.userClient.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Warn("login failed", slog.String("email", req.Email), logger.Err(err))
		writeGRPCError(responseHandler, err)
		return
	}

	log.Info("login successful", slog.String("email", req.Email))
	responseHandler.JSONResponse(map[string]any{"token": token}, http.StatusOK)
}

func writeGRPCError(rh *response.HTTPResponseHandler, err error) {
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.InvalidArgument:
		rh.ErrorResponse(err, st.Message(), http.StatusBadRequest)
	case codes.AlreadyExists:
		rh.ErrorResponse(err, st.Message(), http.StatusConflict)
	case codes.Unauthenticated:
		rh.ErrorResponse(err, st.Message(), http.StatusUnauthorized)
	case codes.NotFound:
		rh.ErrorResponse(err, st.Message(), http.StatusNotFound)
	default:
		rh.ErrorResponse(errors.New("internal server error"), "internal server error", http.StatusInternalServerError)
	}
}
