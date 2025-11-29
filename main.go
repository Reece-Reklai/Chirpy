package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Reece-Reklai/go_serve/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	databaseQuery  *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metrics() string {
	return fmt.Sprint(cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetMetric() {
	cfg.fileserverHits.Swap(int32(0))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func main() {
	var apiCfg apiConfig
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("failed to open database connection")
	}
	dbQueries := database.New(db)
	apiCfg.databaseQuery = dbQueries
	staticDir := "./public/"
	headerMethod := map[string]string{
		"GET":    "GET",
		"POST":   "POST",
		"PUT":    "PUT",
		"DELETE": "DELETE",
	}
	endPoints := map[string]string{
		"api":   "/api",
		"admin": "/admin",
	}
	router := Router{Mux: http.NewServeMux()}
	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router.Mux,
	}
	router.handlers(&apiCfg, staticDir, headerMethod, endPoints)
	log.Fatal(server.ListenAndServe())
}
