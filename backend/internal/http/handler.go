package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"subscribe_tracker/backend/internal/domain"
	"subscribe_tracker/backend/internal/usecase"
)

type Handler struct {
	Auth          usecase.AuthUsecase
	Subscriptions usecase.SubscriptionUsecase
	Tokens        usecase.TokenManager
}

func NewHandler(auth usecase.AuthUsecase, subscriptions usecase.SubscriptionUsecase, tokens usecase.TokenManager) Handler {
	return Handler{
		Auth:          auth,
		Subscriptions: subscriptions,
		Tokens:        tokens,
	}
}

type contextKey string

const userIDKey contextKey = "user_id"

func (h Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.handleRegister)
			r.Post("/login", h.handleLogin)
		})

		r.Group(func(r chi.Router) {
			r.Use(h.authMiddleware)
			r.Get("/subscriptions", h.handleListSubscriptions)
			r.Post("/subscriptions", h.handleCreateSubscription)
			r.Put("/subscriptions/{id}", h.handleUpdateSubscription)
			r.Delete("/subscriptions/{id}", h.handleDeleteSubscription)
		})
	})

	return r
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenValue := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
		if tokenValue == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}

		userID, err := h.Tokens.Parse(tokenValue)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(userIDKey).(string)
	return value, ok
}

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string     `json:"token"`
	User  userResult `json:"user"`
}

type userResult struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	result, err := h.Auth.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toAuthResponse(result))
}

func (h Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	result, err := h.Auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAuthResponse(result))
}

func toAuthResponse(result usecase.AuthResult) authResponse {
	return authResponse{
		Token: result.Token,
		User:  userResult{ID: result.User.ID, Name: result.User.Name, Email: result.User.Email},
	}
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
	case errors.Is(err, usecase.ErrEmailExists):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email already in use"})
	case errors.Is(err, usecase.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "request failed"})
	}
}

type subscriptionPayload struct {
	ServiceName string `json:"service_name"`
	BankName    string `json:"bank_name"`
	CardLast4   string `json:"card_last4"`
	Billing     string `json:"billing_cycle"`
	ChargeDate  string `json:"charge_date"`
}

type subscriptionResult struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	BankName    string `json:"bank_name"`
	CardLast4   string `json:"card_last4"`
	Billing     string `json:"billing_cycle"`
	ChargeDate  string `json:"charge_date"`
}

func (h Handler) handleListSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	items, err := h.Subscriptions.List(r.Context(), userID)
	if err != nil {
		writeSubscriptionError(w, err)
		return
	}

	results := make([]subscriptionResult, 0, len(items))
	for _, item := range items {
		results = append(results, toSubscriptionResult(item))
	}

	writeJSON(w, http.StatusOK, results)
}

func (h Handler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	payload, err := parseSubscriptionPayload(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	item, err := h.Subscriptions.Create(r.Context(), userID, payload)
	if err != nil {
		writeSubscriptionError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toSubscriptionResult(item))
}

func (h Handler) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	payload, err := parseSubscriptionPayload(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	item, err := h.Subscriptions.Update(r.Context(), userID, id, payload)
	if err != nil {
		writeSubscriptionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toSubscriptionResult(item))
}

func (h Handler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.Subscriptions.Delete(r.Context(), userID, id); err != nil {
		writeSubscriptionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func parseSubscriptionPayload(r *http.Request) (usecase.SubscriptionInput, error) {
	var payload subscriptionPayload
	if err := decodeJSON(r, &payload); err != nil {
		return usecase.SubscriptionInput{}, errors.New("invalid payload")
	}

	return usecase.SubscriptionInput{
		ServiceName: payload.ServiceName,
		BankName:    payload.BankName,
		CardLast4:   payload.CardLast4,
		Billing:     payload.Billing,
		ChargeDate:  payload.ChargeDate,
	}, nil
}

func toSubscriptionResult(item domain.Subscription) subscriptionResult {
	return subscriptionResult{
		ID:          item.ID,
		ServiceName: item.ServiceName,
		BankName:    item.BankName,
		CardLast4:   item.CardLast4,
		Billing:     item.Billing,
		ChargeDate:  item.ChargeDate.Format("2006-01-02"),
	}
}

func writeSubscriptionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
	case errors.Is(err, usecase.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	case errors.Is(err, usecase.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "subscription not found"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "request failed"})
	}
}

func decodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
