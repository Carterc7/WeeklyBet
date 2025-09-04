package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
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

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
