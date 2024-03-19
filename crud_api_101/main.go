package main

import (
	"encoding/json" // For encoding and decoding JSON data.
	"fmt"           // Print stuff
	"log"           //  For logging messages to standard error.
	"math/rand"     // Generate random numbers for movie IDs
	"net/http"      // For building HTTP servers and clients.
	"strconv"       //  For converting strings to numeric types and vice versa.

	"github.com/gorilla/mux" // For building flexible and powerful HTTP routers
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) // ... is used to unpact the slice and pass each element as a separate argument
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000)) // Convert random num	ID to string
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatio/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			// Delete existing movie
			movies = append(movies[:index], movies[index+1:]...)
			// Create and add a new movie with same data as sent by user
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {

	movies = append(movies, Movie{ID: "1", Isbn: "46321", Title: "The Batman", Director: &Director{FirstName: "Christopher", LastName: "Nolan"}})
	movies = append(movies, Movie{ID: "2", Isbn: "58031", Title: "Golmaal", Director: &Director{FirstName: "Rohit", LastName: "Shetty"}})
	movies = append(movies, Movie{ID: "3", Isbn: "58031", Title: "Golmaal 2", Director: &Director{FirstName: "Rohit", LastName: "Shetty"}})

	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting port at 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))

}
