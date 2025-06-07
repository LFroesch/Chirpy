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
	db         *database.Queries
	platform   string
	jwt_secret string
	polka_key  string
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
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	polka_key := os.Getenv("POLKA_KEY")
	if jwt_secret == "" {
		log.Fatal("POLKA_KEY must be set")
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
		db:         dbQueries,
		platform:   platform,
		jwt_secret: jwt_secret,
		polka_key:  polka_key,
	}

	// Router
	mux := http.NewServeMux()
	// Home / Healthz / Logo / Index Router
	mux.Handle( // HomePage + RootDir for Logo / Index file etc.
		"/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// Admin Router
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)
	// Basic User Router
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpsByID)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdatePassword)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	// Auth Router
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	// Create Server assign Handler / Address
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// CLI Directory
	log.Printf("| Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("| Printing Directory:")
	log.Println("| Home API ---")
	log.Println("| Home Screen   | N/A  | http://localhost:8080/app/")
	log.Println("| Server Logo   | N/A  | http://localhost:8080/app/assets/logo.png")
	log.Println("| Server Health | GET  | http://localhost:8080/api/healthz")
	log.Println("| Admin API ---")
	log.Println("| List Users    | N/A") // Build this
	log.Println("| Check Metrics | GET  | http://localhost:8080/admin/metrics")
	log.Println("| Reset All     | POST | http://localhost:8080/admin/reset")
	log.Println("| Chirp / User API ---")
	log.Println("| Create User   | POST | http://localhost:8080/api/users")
	log.Println("| Create Chirp  | POST | http://localhost:8080/api/chirps")
	log.Println("| List Chirps   | GET  | http://localhost:8080/api/chirps")
	log.Println("| GetChirp by ID| GET  | http://localhost:8080/api/chirps/{chirpID}")
	log.Println("| Login User    | POST | http://localhost:8080/api/login")
	log.Println("| Auth API ---")
	log.Println("| Revoke Token  | POST | http://localhost:8080/api/revoke")
	log.Println("| Refresh Token | POST | http://localhost:8080/api/refresh")
	log.Println("| Reminder to check browser caching and restart the server after any changes")
	log.Fatal(srv.ListenAndServe())
}
