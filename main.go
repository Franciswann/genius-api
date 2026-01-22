package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Genius API Response Structure (Simplified)
type GeniusResponse struct {
	Response struct {
		Hits []struct {
			Result struct {
				Title  string `json:"title"`
				Artist string `json:"primary_artist_names"`
				URL    string `json:"url"` // Link to the lyrics page
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Read from environment variables; prompt if not found.
	accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")

	if accessToken == "" {
		fmt.Println("Please set the environment variable: GENIUS_ACCESS_TOKEN")
		return
	}

	// add scanner
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter lyrics: ")

	// Read lyrics
	if scanner.Scan() {
		searchQuery := scanner.Text()

		// if user didn't enter, exit
		if strings.TrimSpace(searchQuery) == "" {
			fmt.Println("Plz enter something useful")
			return
		}

		fmt.Printf("\n Searching: %s...\n", searchQuery)
		searchInGenius(searchQuery, accessToken)
	}
}

// searchQuery := "Found the puzzle piece and feel completed"

func searchInGenius(query string, token string) {
	// 1. Create URL and encode Chinese characters
	apiURL := fmt.Sprintf("https://api.genius.com/search?q=%s", url.QueryEscape(query))

	// 2. Create Request
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// 3. Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	// 4. Read and parse JSON
	body, _ := io.ReadAll(resp.Body)
	var result GeniusResponse

	// // 在 json.Unmarshal 之前加入這行，看看原始的「水」長什麼樣子
	// fmt.Println("Raw JSON:", string(body))

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Parse JSON failed", err)
		return
	}

	// 5. print result
	fmt.Println("-----------------")
	if len(result.Response.Hits) == 0 {
		fmt.Println("Couldn't find related song.")
		return
	}

	for i, hit := range result.Response.Hits {
		res := hit.Result
		fmt.Printf("[%d] 🎵 Title: %s\n", i+1, res.Title)
		fmt.Printf("    🎤 Artist: %s\n", res.Artist)
		fmt.Printf("    🔗 Lyrics link: %s\n\n", res.URL)
		// fmt.Println("-----------------")

		// fmt.Printf("apiURL: %s\n", apiURL)
		// fmt.Printf("result: %s\n", result)

	}

}
