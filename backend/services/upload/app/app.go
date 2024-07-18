package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
}

func RegisterAPIRoutes(r *mux.Router) {
	// r.HandleFunc("/deploy", deployHandler).Methods("POST")
	// r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (a *App) Start() error {
	a.Router = mux.NewRouter()
	RegisterAPIRoutes(a.Router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // All headers
		AllowCredentials: true,
	})

	handler := c.Handler(a.Router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	fmt.Println("Starting server on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, handler))

	return nil
}
