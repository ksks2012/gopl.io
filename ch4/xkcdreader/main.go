// Practice 4.12: Write a program that fetches the first ten comics from the xkcd API (https://xkcd.com/json.html) and saves them as JSON files in a directory named "var". If the files already exist, it should load them instead of fetching them again. The program should also allow filtering of the comics based on a substring in the image URL.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Comic struct {
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Link       string `json:"link"`
	News       string `json:"news"`
	Num        int    `json:"num"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
}

func saveJsonToFile() ([]*Comic, error) {
	urlTemplate := "https://xkcd.com/%d/info.0.json"
	savePath := "./var"
	var comics []*Comic

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		fmt.Printf("%s does not exist, creating it...\n", savePath)
		if err := os.MkdirAll(savePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %v", err)
		}
	}

	for i := 1; i <= 10; i++ { // Loop through the first 10 comics
		filePath := fmt.Sprintf("%s/xkcd_%d.json", savePath, i)
		if _, err := os.Stat(filePath); err == nil {
			file, err := os.Open(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
			}
			defer file.Close()
			var comic Comic
			if err := json.NewDecoder(file).Decode(&comic); err != nil {
				return nil, fmt.Errorf("failed to decode JSON from file %s: %v", filePath, err)
			}
			comics = append(comics, &comic)
			fmt.Printf("File %s already exists, loaded.\n", filePath)
			continue
		}
		var result Comic
		resp, err := http.Get(fmt.Sprintf(urlTemplate, i))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch comic %d: %v", i, err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode JSON from response for comic %d: %v", i, err)
		}

		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %v", filePath, err)
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to encode JSON to file %s: %v", filePath, err)
		}
		file.Close()
		fmt.Printf("Saved JSON to %s\n", filePath)
		comics = append(comics, &result)
		time.Sleep(time.Duration(0.5 * float64(time.Second))) // Sleep to avoid hitting the API too fast
	}
	return comics, nil
}

func main() {
	comicData, err := saveJsonToFile()
	if err != nil {
		log.Fatalf("Error fetching or saving comics: %v", err)
	}

	fmt.Println("\n--- All Comics Downloaded/Loaded ---")
	for _, comic := range comicData {
		fmt.Printf("Comic #%d: %s\n", comic.Num, comic.Title)
		fmt.Printf("Image URL: %s\n", comic.Img)
		fmt.Printf("Alt text: %s\n", comic.Alt)
		fmt.Println()
	}

	var filter string
	fmt.Print("Enter a substring to filter comic image URLs (case-sensitive): ")
	_, err = fmt.Scanln(&filter)
	if err != nil && err != io.EOF {
		log.Fatalf("Error reading filter input: %v", err)
	}

	fmt.Println("\n--- Filtered Comics (Image URL Contains Substring) ---")
	if filter == "" {
		fmt.Println("No filter applied (empty input). Displaying all comics again based on image URL.")
	}

	for _, comic := range comicData {
		if filter != "" && !strings.Contains(comic.Img, filter) {
			continue
		}
		fmt.Printf("Comic #%d: %s\n", comic.Num, comic.Title)
		fmt.Printf("Image URL: %s\n", comic.Img)
		fmt.Printf("Alt text: %s\n", comic.Alt)
		fmt.Println()
	}
}
