package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/LFroesch/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	// add *database.Queries to the api config struct
	db       *database.Queries
	platform string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// load the .env with godotenv
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	// Open db connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	// get queries from database
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		// add queries / platform to config
		db:       dbQueries,
		platform: platform,
	}

	// Router
	mux := http.NewServeMux()

	// API Router
	mux.Handle( // HomePage + RootDir for Logo / Index file etc.
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// mux.HandleFunc("POST /api/validate_chirp", chirpHandler) Disabled & Moved to Create Chirp Logic

	// Admin Router
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Basic User Router
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	// Create Server assign Handler / Address
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// CLI Directory
	log.Printf("| Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("| Printing Directory:")
	log.Println("| Base API ---")
	log.Println("| Home Screen   | http://localhost:8080/app/")
	log.Println("| Server Logo   | http://localhost:8080/app/assets/logo.png") // this is a bonus link
	log.Println("| Server Health | http://localhost:8080/api/healthz")
	log.Println("| Admin API ---")
	log.Println("| List Users    | N/A") // Build this
	log.Println("| List Chirps   | N/A") // Build this
	log.Println("| Check Metrics | http://localhost:8080/admin/metrics")
	log.Println("| Reset All     | http://localhost:8080/admin/reset")
	log.Println("| Chirp / User API ---")
	log.Println("| Create User   | http://localhost:8080/api/users")
	log.Println("| Create Chirp  | http://localhost:8080/api/chirps")
	log.Println("----")
	log.Println("| Reminder to check browser caching and restart the server after any changes")
	log.Fatal(srv.ListenAndServe())
}
