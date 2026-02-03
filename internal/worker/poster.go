package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/yoonhyunwoo/crab-trap/pkg/moltbook"
)

type Poster struct {
	client      *moltbook.Client
	submolt     string
	pauseBetween time.Duration
	serverURL   string
}

func NewPoster(apiKey, submolt, serverURL string) *Poster {
	return &Poster{
		client:      moltbook.NewClient(apiKey),
		submolt:     submolt,
		pauseBetween: 5 * time.Minute,
		serverURL:   serverURL,
	}
}

func (p *Poster) PostPrompt(prompt string, templateName string) error {
	title := fmt.Sprintf("Prompt Injection Test: %s", templateName)
	
	_, err := p.client.CreatePost(moltbook.CreatePostRequest{
		Submolt: p.submolt,
		Title:   title,
		Content: prompt,
	})

	if err != nil {
		if rateErr, ok := err.(moltbook.RateLimitError); ok {
			log.Printf("Rate limited: %v, waiting...", rateErr.Error())
			time.Sleep(30 * time.Minute)
			return p.PostPrompt(prompt, templateName)
		}
		return fmt.Errorf("failed to post prompt: %w", err)
	}

	log.Printf("Successfully posted prompt: %s", templateName)
	return nil
}

func (p *Poster) PostPrompts(prompts map[string]string) error {
	count := 0
	for name, prompt := range prompts {
		if err := p.PostPrompt(prompt, name); err != nil {
			log.Printf("Failed to post prompt %s: %v", name, err)
			continue
		}

		count++

		if count < len(prompts) {
			log.Printf("Pausing for %v before next post...", p.pauseBetween)
			time.Sleep(p.pauseBetween)
		}
	}

	return nil
}

func (p *Poster) PostAllTemplates(generator *Generator) error {
	prompts := generator.GenerateAll()

	for _, prompt := range prompts {
		title := fmt.Sprintf("Prompt Injection Test (Random) - %s", time.Now().Format("2006-01-02 15:04:05"))
		
		if err := p.PostWithRetry(title, prompt); err != nil {
			return err
		}
	}

	return nil
}

func (p *Poster) PostWithRetry(title, content string) error {
	maxRetries := 3
	var lastErr error
	var postURL string

	for i := 0; i < maxRetries; i++ {
		resp, err := p.client.CreatePost(moltbook.CreatePostRequest{
			Submolt: p.submolt,
			Title:   title,
			Content: content,
		})

		if err == nil && resp.Success {
			postURL = resp.Post.URL
			break
		}

		lastErr = err

		if rateErr, ok := err.(moltbook.RateLimitError); ok {
			waitTime := time.Duration(rateErr.RetryAfterMinutes) * time.Minute
			if waitTime == 0 {
				// 30 min cooldown + jitter (0 to +5 min)
				jitter := time.Duration(rand.Intn(5)) * time.Minute
				waitTime = 30*time.Minute + jitter
			} else {
				// Use API provided time + jitter (0 to +5 min)
				jitter := time.Duration(rand.Intn(5)) * time.Minute
				waitTime = waitTime + jitter
			}
			log.Printf("Rate limited, waiting %v before retry %d/%d", waitTime, i+1, maxRetries)
			time.Sleep(waitTime)
		} else {
			log.Printf("Error posting (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(1 * time.Minute)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
	}

	p.notifyServer(title, postURL)

	return nil
}

func (p *Poster) notifyServer(title, url string) {
	if p.serverURL == "" || url == "" {
		return
	}

	record := map[string]string{
		"title": title,
		"url":   url,
	}

	body, err := json.Marshal(record)
	if err != nil {
		return
	}

	_, err = http.Post(p.serverURL+"/post", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to notify server of post: %v", err)
	}
}

func (p *Poster) SetPauseBetween(duration time.Duration) {
	p.pauseBetween = duration
}
