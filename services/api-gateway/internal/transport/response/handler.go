package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
)

type HTTPResponseHandler struct {
	log *slog.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(log *slog.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(resp any, statusCode int) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(resp); err != nil {
		h.log.Error("failed to encode response", logger.Err(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string, statusCode int) {
	response := map[string]string{
		"error":   err.Error(),
		"message": msg,
	}
	h.JSONResponse(response, statusCode)
}
