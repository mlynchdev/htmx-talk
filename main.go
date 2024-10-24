package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	tmdbBaseURL  = "https://api.themoviedb.org/3"
	imageBaseURL = "https://image.tmdb.org/t/p/w200"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/active-search", activeSearchHandler)

	// Serve static files (CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func activeSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/active-search.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		query := r.FormValue("query")
		results, err := searchMovies(query)
		if err != nil {
			http.Error(w, "Error fetching movies", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/search-results.html"))
		tmpl.Execute(w, results)
	}
}

type Movie struct {
	Title     string
	PosterURL string
}

func searchMovies(query string) ([]Movie, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY environment variable is not set")
	}

	// Search for movies
	searchURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", tmdbBaseURL, apiKey, url.QueryEscape(query))
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResults struct {
		Results []struct {
			ID         int    `json:"id"`
			Title      string `json:"title"`
			PosterPath string `json:"poster_path"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return nil, err
	}

	var movies []Movie
	for _, result := range searchResults.Results {
		

		movie := Movie{
			Title:    result.Title,
		}

		if result.PosterPath != "" {
			movie.PosterURL = imageBaseURL + result.PosterPath
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
