package github

import (
	"encoding/json"
	"log"
)

// RepoStarDetails contains details about a GitHub star event.
type RepoStarDetails struct {
	Repo         string
	RepoURL      string
	StarringUser string
	TotalStars   string
}

// GitHub star webhook payload.
type StarEvent struct {
	Action     string `json:"action"`
	StarredAt  string `json:"starred_at"`
	Repository struct {
		FullName string `json:"full_name"`
		Stars    string `json:"stargazers_count"`
	} `json:"repository"`
	URL    string `json:"html_url"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

// Ensure that StarWebhookHandler implements the WebhookHandler interface.
var _ WebhookHandler = StarWebhookHandler{}

// StarWebhookHandler handles GitHub star webhook payloads.
type StarWebhookHandler struct {
	messageSenders []func(RepoStarDetails)
}

// NewStarWebhookHandler returns a new StarWebhookHandler.
func NewStarWebhookHandler(messageSenders []func(RepoStarDetails)) StarWebhookHandler {
	return StarWebhookHandler{
		messageSenders: messageSenders,
	}
}

// Handle processes the GitHub star webhook payload.
// It extracts the new star event details and sends messages
// via all provided notification funcs.
func (r StarWebhookHandler) Handle(payload []byte) {
	var event StarEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Error unmarshalling GitHub star event: %v", err)
		return
	}

	if event.Action != "created" {
		log.Panicf("Unsupported GitHub star event action: %s", event.Action)
		return
	}

	for _, sendNotification := range r.messageSenders {
		sendNotification(RepoStarDetails{
			Repo:         event.Repository.FullName,
			RepoURL:      event.URL,
			StarringUser: event.Sender.Login,
			TotalStars:   event.Repository.Stars,
		})
	}
}
