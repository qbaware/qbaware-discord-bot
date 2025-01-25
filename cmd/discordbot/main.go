package main

import (
	"log"
	"net/http"
	"os"
	"qbaware-discord-bot/internal/discord"
	"qbaware-discord-bot/internal/gh"
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
	http.HandleFunc("/gh-webhook", func(w http.ResponseWriter, r *http.Request) {
		gh.HandleNewReleaseWebhook(w, r, webhookSecret, func(repoFullName string, releaseName string, releaseVersion string, releaseURL string, releaseBody string) {
			err := dc.SendNewReleaseNotification(repoFullName, releaseName, releaseVersion, releaseURL, releaseBody)
			if err != nil {
				log.Printf("Error sending release notification: %s", err.Error())
			}
		})
	})

	// Start HTTP server on the main thread.
	log.Println("Starting GitHub webhook server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting HTTP server: ", err)
	}
}
