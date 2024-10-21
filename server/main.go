package main

import (
	"log"
	"net/http"
	"os"
    "fmt"
	"github.com/gorilla/mux"
     "github.com/joho/godotenv"
	"github.com/Ghat0tkach/Gohacktober-Backend/internal/handlers"
	"github.com/Ghat0tkach/Gohacktober-Backend/config"
	"github.com/Ghat0tkach/Gohacktober-Backend/internal/github"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize GitHub client
	github.Init(cfg)

	// Initialize handlers
	handlers.Init(cfg)

	// Create a new router
	r := mux.NewRouter()

	// Set up routes
    r.HandleFunc("/api/hacktoberfest-contributions", handlers.GetHacktoberfestContributionsHandler)
	r.HandleFunc("/auth/github", handlers.GitHubAuthHandler)
	r.HandleFunc("/auth/github/callback", handlers.GitHubCallbackHandler)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to Gohacktober API")
	})
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println(`
    GGGG   OOO      H   H   AAAAA   CCCCC   K   K   TTTTT   OOO   BBBBB   EEEEE   RRRR  
   G      O   O     H   H   A   A   C       K  K      T    O   O  B    B  E       R   R
   G  GG  O   O     HHHHH   AAAAA   C       KKK       T    O   O  BBBBB   EEEE    RRRR 	
   G   G  O   O     H   H   A   A   C       K  K      T    O   O  B    B  E       R  R 
   GGGG    OOO      H   H   A   A   CCCCC   K   K     T     OOO   BBBBB   EEEEE   R   R

 BBBBB    AAAAA   CCCCC   K   K   EEEEE   N   N   DDDD  
 B    B   A   A   C       K  K    E       NN  N   D   D
 BBBBB    AAAAA   C       KKK     EEEE    N N N   D   D
 B    B   A   A   C       K  K    E       N  NN   D   D
 BBBBB    A   A   CCCCC   K   K   EEEEE   N   N   DDDD  
`)
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}