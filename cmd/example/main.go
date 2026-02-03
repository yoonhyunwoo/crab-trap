package main

import (
	"fmt"
	"log"

	"github.com/yoonhyunwoo/crab-trap/pkg/moltbook"
)

func main() {

	apiKey := "YOUR_API_KEY"

	client := moltbook.NewClient(apiKey)

	me, err := client.GetMe()
	if err != nil {
		log.Fatalf("Failed to get profile: %v", err)
	}

	fmt.Printf("Hello, %s! Karma: %d\n", me.Name, me.Karma)

	posts, err := client.GetPosts(moltbook.GetPostsOptions{
		Sort:  "hot",
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Failed to get posts: %v", err)
	}

	fmt.Printf("Found %d posts\n", len(posts.Posts))
	for i, post := range posts.Posts {
		fmt.Printf("%d. %s by %s (↑%d)\n", i+1, post.Title, post.Author.Name, post.Upvotes)
	}

	newPost, err := client.CreatePost(moltbook.CreatePostRequest{
		Submolt: "general",
		Title:   "Hello from Go SDK!",
		Content: "This is my first post using the Moltbook Go SDK.",
	})
	if err != nil {
		log.Printf("Failed to create post: %v", err)
	} else {
		fmt.Printf("Created post: %s\n", newPost.Post.Title)
	}

	searchResults, err := client.Search("how do agents handle memory", moltbook.SearchRequest{
		Type:  "all",
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	fmt.Printf("Found %d search results\n", searchResults.Count)
	for i, result := range searchResults.Results {
		title := "Comment"
		if result.Title != nil {
			title = *result.Title
		}
		fmt.Printf("%d. %s (similarity: %.2f)\n", i+1, title, result.Similarity)
	}

	feed, err := client.GetFeed(moltbook.FeedOptions{
		Sort:  "new",
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Failed to get feed: %v", err)
	}

	fmt.Printf("Your feed has %d posts\n", len(feed.Posts))
}
