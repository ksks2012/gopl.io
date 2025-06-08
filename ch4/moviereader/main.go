// Practice 4.12: Write a program that fetches the first ten comics from the xkcd API (https://xkcd.com/json.html) and saves them as JSON files in a directory named "var". If the files already exist, it should load them instead of fetching them again. The program should also allow filtering of the comics based on a substring in the image URL.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type Movie struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	ImdbRating string   `json:"imdbRating"`
	ImdbVotes  string   `json:"imdbVotes"`
	ImdbID     string   `json:"imdbID"`
	Type       string   `json:"Type"`
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"`
}

func searchMovie(searchTitle string, searchYear string) ([]*Movie, error) {
	apikey := os.Getenv("OMDB_APIKEY")
	if apikey == "" {
		return nil, fmt.Errorf("OMDB_APIKEY environment variable not set")
	}

	params := url.Values{}
	params.Add("apikey", apikey)
	if searchTitle != "" {
		params.Add("t", searchTitle)
	}
	if searchYear != "" {
		params.Add("y", searchYear)
	}

	fullURL := "http://www.omdbapi.com/?" + params.Encode()
	fmt.Println("Fetching movie data from:", fullURL)

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyGoClient/1.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return []*Movie{&movie}, nil
}

func main() {
	fmt.Print("Enter the movie title: ")
	var title string
	fmt.Scanln(&title)

	fmt.Print("Enter the year of the movie (press Enter to skip): ")
	var year string
	fmt.Scanln(&year)

	movies, err := searchMovie(title, year)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, movie := range movies {
		fmt.Printf("Title: %s\nYear: %s\nPoster: %s\n\n", movie.Title, movie.Year, movie.Poster)
	}

}
