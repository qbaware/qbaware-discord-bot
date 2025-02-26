package main

import (
	"log"
	"net/http"
	"os"

	"qbaware-discord-bot/internal/discord"
	"qbaware-discord-bot/internal/github"
)

func main() {
	// Extract envornment variables.
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("Please set the DISCORD_BOT_TOKEN environment variable")
	}
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("Please set the GITHUB_WEBHOOK_SECRET environment variable")
	}

	// Initialize the Discord connection.
	dc := discord.NewConnection(token)
	err := dc.Init()
	if err != nil {
		log.Fatal("Error initializing Discord connection: " + err.Error())
	}
	log.Println("Discord connection initialized successfully")

	// Add GitHub webhook HTTP server.
	ghReleaseNotificationHandler := github.NewReleaseWebhookHandler([]func(github.RepoReleaseDetails){
		func(rrd github.RepoReleaseDetails) {
			dc.SendNewReleaseNotification(rrd.Repo, rrd.Name, rrd.Version, rrd.URL, rrd.Body)
		},
	})
	ghStarNotificationHandler := github.NewStarWebhookHandler([]func(github.RepoStarDetails){
		func(rsd github.RepoStarDetails) {
			dc.SendNewStarNotification(rsd.Repo, rsd.RepoURL, rsd.StarringUser, rsd.TotalStars)
		},
	})
	ghWebhookHandlers := []github.WebhookHandler{
		ghReleaseNotificationHandler,
		ghStarNotificationHandler,
	}
	http.HandleFunc("/gh-webhook", github.NewWebhookHandler(webhookSecret, ghWebhookHandlers))

	// Start HTTP server on the main thread.
	log.Println("Starting GitHub webhook server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting HTTP server: ", err)
	}
}
