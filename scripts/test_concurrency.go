package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostpayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func updatePost(postID int, p UpdatePostpayload, wg *sync.WaitGroup) {
	defer wg.Done()

	//Construct the URL for the update endpoint
	url := fmt.Sprintf("http://localhost:8080/v1/posts/%d", postID)
	fmt.Println("Updating post with URL: ", url)

	// Create the JSON payload
	b, _ := json.Marshal(p)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error creating request:")
		return
	}

	// Set headers as needed, for example:
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:")
		return
	}
	defer resp.Body.Close()

	fmt.Println("Update response status:", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	// Assuming the post ID to update is 1
	postID := 5

	// Simulate User A and User B updating the same post concurrently
	wg.Add(2)
	content := "NEW CONTENT FROM USER B"
	title := "NEW TITLE FROM USER A"

	go updatePost(postID, UpdatePostpayload{Title: &title}, &wg)
	go updatePost(postID, UpdatePostpayload{Content: &content}, &wg)
	wg.Wait()

}
