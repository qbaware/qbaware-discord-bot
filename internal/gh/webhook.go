package gh

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// GitHub release webhook payload.
type GitHubReleaseEvent struct {
	Action  string `json:"action"`
	Release struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	} `json:"release"`
}

// verifySignature verifies the GitHub webhook signature.
func verifySignature(payload []byte, signature string, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature[7:]), []byte(expectedMAC))
}

// HandleNewReleaseWebhook handles the GitHub webhook for new releases.
func HandleNewReleaseWebhook(w http.ResponseWriter, r *http.Request, webhookSecret string, sendNotification func(releaseName string, releaseVersion string, releaseURL string, releaseBody string)) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Verify GitHub webhook signature.
	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(payload, signature, webhookSecret) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse the event type.
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "release" {
		w.WriteHeader(http.StatusOK) // Acknowledge but ignore non-release events.
		return
	}

	// Parse the event.
	var event GitHubReleaseEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		http.Error(w, "Error parsing webhook payload", http.StatusBadRequest)
		return
	}

	// Only process "published" release events.
	if event.Action != "published" {
		w.WriteHeader(http.StatusOK)
		return
	}

	releaseName := event.Release.Name
	releaseVersion := event.Release.TagName
	releaseURL := event.Release.HTMLURL
	releaseBody := event.Release.Body

	sendNotification(releaseName, releaseVersion, releaseURL, releaseBody)

	w.WriteHeader(http.StatusOK)
}
