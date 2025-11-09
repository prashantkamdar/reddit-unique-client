package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

type RedditPost struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Subreddit   string  `json:"subreddit"`
	Author      string  `json:"author"`
	Score       int     `json:"score"`
	NumComments int     `json:"num_comments"`
	Created     float64 `json:"created_utc"`
	URL         string  `json:"url"`
	Permalink   string  `json:"permalink"`
	Thumbnail   string  `json:"thumbnail"`
	SelfText    string  `json:"selftext"`
	IsVideo     bool    `json:"is_video"`
}

type RedditResponse struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
		After string `json:"after"`
	} `json:"data"`
}

type PostsResponse struct {
	Posts []RedditPost `json:"posts"`
	Error string       `json:"error,omitempty"`
}

func main() {
	// Serve static files
	http.Handle("/", http.FileServer(http.FS(staticFiles)))

	// API endpoint to fetch user posts
	http.HandleFunc("/api/user/", handleUserPosts)

	port := ":8080"
	log.Printf("Reddit Client starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleUserPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract username from URL path
	username := r.URL.Path[len("/api/user/"):]
	if username == "" {
		json.NewEncoder(w).Encode(PostsResponse{Error: "Username is required"})
		return
	}

	log.Printf("Fetching posts for user: %s", username)

	posts, err := fetchAllUserPosts(username)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		json.NewEncoder(w).Encode(PostsResponse{Error: err.Error()})
		return
	}

	// Deduplicate posts by title (handles crossposts)
	uniquePosts := deduplicatePosts(posts)
	log.Printf("User %s: %d posts (%d duplicates removed)", username, len(uniquePosts), len(posts)-len(uniquePosts))

	// Sort chronologically (newest first)
	sort.Slice(uniquePosts, func(i, j int) bool {
		return uniquePosts[i].Created > uniquePosts[j].Created
	})

	json.NewEncoder(w).Encode(PostsResponse{Posts: uniquePosts})
}

func fetchAllUserPosts(username string) ([]RedditPost, error) {
	var allPosts []RedditPost
	after := ""
	limit := 100

	// Fetch multiple pages to get more posts
	for i := 0; i < 10; i++ { // Limit to 10 pages (1000 posts max)
		url := fmt.Sprintf("https://www.reddit.com/user/%s/submitted.json?limit=%d", username, limit)
		if after != "" {
			url += "&after=" + after
		}

		posts, nextAfter, err := fetchUserPostsPage(url)
		if err != nil {
			return nil, err
		}

		allPosts = append(allPosts, posts...)

		if nextAfter == "" || len(posts) == 0 {
			break
		}
		after = nextAfter

		// Be respectful to Reddit's rate limits
		time.Sleep(500 * time.Millisecond)
	}

	return allPosts, nil
}

func fetchUserPostsPage(url string) ([]RedditPost, string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// Set User-Agent as Reddit requires it
	req.Header.Set("User-Agent", "RedditClient/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("reddit API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var redditResp RedditResponse
	if err := json.Unmarshal(body, &redditResp); err != nil {
		return nil, "", err
	}

	var posts []RedditPost
	for _, child := range redditResp.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts, redditResp.Data.After, nil
}

func deduplicatePosts(posts []RedditPost) []RedditPost {
	// Group posts by title to identify duplicates (crosspost detection)
	// Keep the earliest post (by created_utc) from each group
	postMap := make(map[string]*RedditPost)

	for i := range posts {
		post := &posts[i]

		// Use title as the primary deduplication key
		// This catches crossposts even when they're posted as different types
		// (e.g., original as self-post, crosspost as link to image)
		key := post.Title

		existing, exists := postMap[key]
		if !exists {
			// First time seeing this title
			postMap[key] = post
		} else {
			// We've seen this title before - keep the earlier one
			if post.Created < existing.Created {
				postMap[key] = post
			}
		}
	}

	// Convert map back to slice
	var unique []RedditPost
	for _, post := range postMap {
		unique = append(unique, *post)
	}

	return unique
}
