package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscribe_tracker/backend/internal/config"
	"subscribe_tracker/backend/internal/db"
	httpapi "subscribe_tracker/backend/internal/http"
	"subscribe_tracker/backend/internal/repository/postgres"
	"subscribe_tracker/backend/internal/security"
	"subscribe_tracker/backend/internal/usecase"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" {
		log.Fatal("DATABASE_URL and JWT_SECRET must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.ApplyMigrations(context.Background(), pool, cfg.MigrationsDir); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	tokenManager := security.NewJWTManager([]byte(cfg.JWTSecret), 7*24*time.Hour)

	userRepo := postgres.NewUserRepository(pool)
	subRepo := postgres.NewSubscriptionRepository(pool)

	authUC := usecase.NewAuthUsecase(userRepo, tokenManager)
	subUC := usecase.NewSubscriptionUsecase(subRepo)

	handler := httpapi.NewHandler(authUC, subUC, tokenManager)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           withCORS(cfg.CorsOrigins, handler.Routes()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("API listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func withCORS(allowed []string, next http.Handler) http.Handler {
	allowedSet := make(map[string]struct{})
	for _, origin := range allowed {
		if origin != "" {
			allowedSet[origin] = struct{}{}
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if len(allowedSet) == 0 {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if _, ok := allowedSet[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
