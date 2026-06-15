package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"todolist/internal/apperrors"
	"todolist/internal/model"
	"todolist/internal/service"
)

const maxBody = 1 << 20

type TodoHTTPHandler struct {
	svc    *service.TodoService
	logger *slog.Logger
}

func NewTodoHTTPHandler(svc *service.TodoService, logger *slog.Logger) *TodoHTTPHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &TodoHTTPHandler{svc: svc, logger: logger}
}

func (h *TodoHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /todos", http.HandlerFunc(h.create))
	mux.Handle("GET /todos", http.HandlerFunc(h.list))
	mux.Handle("GET /todos/{id}", http.HandlerFunc(h.get))
	mux.Handle("PUT /todos/{id}", http.HandlerFunc(h.update))
	mux.Handle("DELETE /todos/{id}", http.HandlerFunc(h.delete))
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

type errBody struct {
	Error string `json:"error"`
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, errBody{Error: msg})
}

func statusFromErr(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrTaskTitleEmpty),
		errors.Is(err, apperrors.ErrDueDateInvalid),
		errors.Is(err, apperrors.ErrTodoIDNotPositive):
		return http.StatusBadRequest
	case errors.Is(err, apperrors.ErrTodoNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func decodeJSON[T any](r *http.Request, dst *T) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return err
	}
	return nil
}

func (h *TodoHTTPHandler) create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	todo, err := h.svc.Create(r.Context(), req)
	if err != nil {
		code := statusFromErr(err)
		if code == http.StatusInternalServerError {
			h.logger.Error("create", slog.String("error", err.Error()))
		}
		writeErr(w, code, service.PublicErrorMessage(err))
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHTTPHandler) get(w http.ResponseWriter, r *http.Request) {
	todo, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		code := statusFromErr(err)
		if code == http.StatusInternalServerError {
			h.logger.Error("get", slog.String("error", err.Error()))
		}
		writeErr(w, code, service.PublicErrorMessage(err))
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHTTPHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("include_completed")
	include := false
	if q != "" {
		b, err := strconv.ParseBool(q)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "invalid include_completed")
			return
		}
		include = b
	}
	todos, err := h.svc.List(r.Context(), include)
	if err != nil {
		h.logger.Error("list", slog.String("error", err.Error()))
		writeErr(w, http.StatusInternalServerError, service.PublicErrorMessage(err))
		return
	}
	if todos == nil {
		todos = []model.Todo{}
	}
	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHTTPHandler) update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	todo, err := h.svc.Update(r.Context(), r.PathValue("id"), req)
	if err != nil {
		code := statusFromErr(err)
		if code == http.StatusInternalServerError {
			h.logger.Error("update", slog.String("error", err.Error()))
		}
		writeErr(w, code, service.PublicErrorMessage(err))
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHTTPHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
		code := statusFromErr(err)
		if code == http.StatusInternalServerError {
			h.logger.Error("delete", slog.String("error", err.Error()))
		}
		writeErr(w, code, service.PublicErrorMessage(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
