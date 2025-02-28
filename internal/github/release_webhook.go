package github

import (
	"encoding/json"
	"log"
)

// RepoReleaseDetails contains details about a GitHub release message.
type RepoReleaseDetails struct {
	Repo    string
	Name    string
	Version string
	URL     string
	Body    string
}

// GitHub release webhook payload.
type ReleaseEvent struct {
	Action     string `json:"action"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Release struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	} `json:"release"`
}

// Ensure that ReleaseWebhookHandler implements the WebhookHandler interface.
var _ WebhookHandler = ReleaseWebhookHandler{}

// ReleaseWebhookHandler handles GitHub release webhook payloads.
type ReleaseWebhookHandler struct {
	sendMsgFuncs []func(RepoReleaseDetails)
}

// NewReleaseWebhookHandler returns a new ReleaseWebhookHandler.
func NewReleaseWebhookHandler(messageSenders []func(RepoReleaseDetails)) ReleaseWebhookHandler {
	return ReleaseWebhookHandler{
		sendMsgFuncs: messageSenders,
	}
}

// Handle processes the GitHub release webhook payload.
// It extracts the new release details and sends messages
// via all provided notification funcs.
func (r ReleaseWebhookHandler) Handle(payload []byte) {
	var event ReleaseEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Error unmarshalling GitHub release event: %v", err)
		return
	}

	// Only process "published" release events.
	if event.Action != "published" {
		log.Printf("Unsupported GitHub release action: %s, skipping processing", event.Action)
		return
	}

	releaseRepo := event.Repository.FullName
	releaseName := event.Release.Name
	releaseVersion := event.Release.TagName
	releaseURL := event.Release.HTMLURL
	releaseBody := event.Release.Body

	for _, sendNotification := range r.sendMsgFuncs {
		sendNotification(RepoReleaseDetails{
			Repo:    releaseRepo,
			Name:    releaseName,
			Version: releaseVersion,
			URL:     releaseURL,
			Body:    releaseBody,
		})
	}
}
