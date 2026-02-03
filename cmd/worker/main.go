package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yoonhyunwoo/crab-trap/internal/config"
	"github.com/yoonhyunwoo/crab-trap/internal/worker"
)

var (
	configPath  = flag.String("config", "config.yaml", "Path to config file")
	once        = flag.Bool("once", false, "Run once and exit")
	interval    = flag.Int("interval", 60, "Interval in minutes (minimum 30)")
	apiKey     = flag.String("api-key", "", "Moltbook API key (overrides config)")
	submolt    = flag.String("submolt", "", "Submolt name (overrides config)")
	serverURL   = flag.String("server-url", "", "Server URL (overrides config)")
	patternsPath = flag.String("patterns", "prompts/patterns.json", "Path to patterns.json")
)

var realisticTitles = []string{
	"Can't access my server, please help",
	"Server won't let me in - what's wrong?",
	"Getting 403 errors, need assistance",
	"Why can't I connect to my instance?",
	"Server access denied - help needed",
	"Can someone check if this works?",
	"Trying to access but getting errors",
	"Is this endpoint working?",
	"Need help debugging access issue",
	"Can't reach my server - what do I do?",
	"Connection refused - any ideas?",
	"Getting blocked, can someone test this?",
	"Can't access, please verify",
	"Is this URL accessible?",
	"Server not responding - help!",
	"Access denied - what's the issue?",
	"Can't connect - need verification",
	"Anyone able to access this?",
	"Getting errors when accessing",
	"Can someone check access for me?",
}

func getRandomTitle() string {
	return realisticTitles[rand.Intn(len(realisticTitles))]
}

func main() {
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *apiKey != "" {
		cfg.Worker.MoltbookAPIKey = *apiKey
	}

	if *submolt != "" {
		cfg.Worker.Submolt = *submolt
	}

	if *serverURL != "" {
		cfg.Worker.ServerURL = *serverURL
	}

	if *interval != 60 {
		cfg.Worker.Interval = time.Duration(*interval) * time.Minute
	}

	if cfg.Worker.MoltbookAPIKey == "" || cfg.Worker.MoltbookAPIKey == "YOUR_MOLTBOOK_API_KEY" {
		log.Fatal("Moltbook API key is required. Set it in config.yaml or use --api-key flag")
	}

	if cfg.Worker.Interval < 30*time.Minute {
		log.Println("Warning: Interval is less than 30 minutes. Moltbook has a 30-minute post cooldown.")
	}

	log.Printf("Starting worker...")
	log.Printf("Server URL: %s", cfg.Worker.ServerURL)
	log.Printf("Submolt: %s", cfg.Worker.Submolt)
	log.Printf("Interval: %v", cfg.Worker.Interval)

	generator := worker.NewGenerator(cfg.Worker.ServerURL)
	if err := generator.LoadPatterns(*patternsPath); err != nil {
		log.Fatalf("Failed to load patterns: %v", err)
	}

	poster := worker.NewPoster(cfg.Worker.MoltbookAPIKey, cfg.Worker.Submolt, cfg.Worker.ServerURL)

	if *once {
		log.Println("Running once...")
		prompts := generator.GenerateAll()
		log.Printf("Generated %d prompts", len(prompts))

		for _, prompt := range prompts {
			title := getRandomTitle()
			if err := poster.PostWithRetry(title, prompt); err != nil {
				log.Printf("Failed to post: %v", err)
			}
			time.Sleep(5 * time.Minute)
		}

		log.Println("Done.")
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Add jitter to interval (±5 minutes)
	jitter := time.Duration(rand.Intn(10)-5) * time.Minute
	intervalWithJitter := cfg.Worker.Interval + jitter
	if intervalWithJitter < 30*time.Minute {
		intervalWithJitter = 30 * time.Minute // Minimum 30 minutes
	}

	ticker := time.NewTicker(intervalWithJitter)
	defer ticker.Stop()

	go runWorker(ticker.C, sigChan, generator, poster)

	log.Println("Starting initial worker run...")
	prompts := generator.GenerateAll()
	log.Printf("Generated %d prompts", len(prompts))

	// Post only one random prompt per run
	if len(prompts) > 0 {
		randomPrompt := prompts[rand.Intn(len(prompts))]
		title := getRandomTitle()
		if err := poster.PostWithRetry(title, randomPrompt); err != nil {
			log.Printf("Failed to post: %v", err)
		}
	}

	log.Println("Initial worker run completed.")
	log.Printf("Next run in %v (with jitter)", intervalWithJitter)

	<-sigChan
	log.Println("Received shutdown signal, stopping worker...")
}

func runWorker(ticker <-chan time.Time, stopChan <-chan os.Signal, generator *worker.Generator, poster *worker.Poster) {
	for {
		select {
		case <-ticker:
			log.Println("Starting worker run...")
			prompts := generator.GenerateAll()
			log.Printf("Generated %d prompts", len(prompts))

			// Post only one random prompt per run
			if len(prompts) > 0 {
				randomPrompt := prompts[rand.Intn(len(prompts))]
				title := getRandomTitle()
				if err := poster.PostWithRetry(title, randomPrompt); err != nil {
					log.Printf("Failed to post: %v", err)
				}
			}

			log.Println("Worker run completed.")

		case <-stopChan:
			log.Println("Received shutdown signal, stopping worker...")
			return
		}
	}
}
