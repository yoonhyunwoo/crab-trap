package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/your-username/moltbook-prompt-injector/pkg/moltbook"
)

type Poster struct {
	client      *moltbook.Client
	submolt     string
	pauseBetween time.Duration
}

func NewPoster(apiKey, submolt string) *Poster {
	return &Poster{
		client:      moltbook.NewClient(apiKey),
		submolt:     submolt,
		pauseBetween: 5 * time.Minute,
	}
}

func (p *Poster) PostPrompt(prompt string, templateName string) error {
	title := fmt.Sprintf("🦞 Prompt Injection Test: %s", templateName)
	
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
		title := fmt.Sprintf("🦞 Prompt Injection Test (Random) - %s", time.Now().Format("2006-01-02 15:04:05"))
		
		if err := p.postWithRetry(title, prompt); err != nil {
			return err
		}
	}

	return nil
}

func (p *Poster) postWithRetry(title, content string) error {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		_, err := p.client.CreatePost(moltbook.CreatePostRequest{
			Submolt: p.submolt,
			Title:   title,
			Content: content,
		})

		if err == nil {
			return nil
		}

		lastErr = err

		if rateErr, ok := err.(moltbook.RateLimitError); ok {
			waitTime := time.Duration(rateErr.RetryAfterMinutes) * time.Minute
			log.Printf("Rate limited, waiting %v before retry %d/%d", waitTime, i+1, maxRetries)
			time.Sleep(waitTime)
		} else {
			log.Printf("Error posting (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(1 * time.Minute)
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func (p *Poster) SetPauseBetween(duration time.Duration) {
	p.pauseBetween = duration
}
