package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Log debug information if DEBUG is enabled
	if os.Getenv("DEBUG") == "true" {
		log.Printf("--- DEBUG MODE ---")
	}

	// Parse templates
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title string
		}{
			Title: "WeeklyBet - Home",
		}
		tmpl.ExecuteTemplate(w, "base", data)
	})

	// HTMX API endpoints
	http.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "<p class='text-green-600 font-semibold'>Current time: %s</p>", currentTime)
	})

	// NFL Schedules API endpoint
	http.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		// Get API key from environment
		apiKey := os.Getenv("RAPIDAPI_KEY")
		if apiKey == "" {
			http.Error(w, "RAPIDAPI_KEY environment variable not set", http.StatusInternalServerError)
			return
		}

		// Create request to NFL API
		url := "https://tank01-nfl-live-in-game-real-time-statistics-nfl.p.rapidapi.com/getNFLGamesForWeek?week=1&seasonType=reg&season=2025"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Set headers
		req.Header.Add("x-rapidapi-host", "tank01-nfl-live-in-game-real-time-statistics-nfl.p.rapidapi.com")
		req.Header.Add("x-rapidapi-key", apiKey)

		// Make the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch NFL data", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// Read response body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		// Set content type and return JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.StatusCode)
		w.Write(body)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
