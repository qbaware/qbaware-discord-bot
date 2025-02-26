package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

const (
	signatureHeader = "X-Hub-Signature-256"
)

// WebhookHandler is an interface for handling GitHub webhook payloads.
type WebhookHandler interface {
	Handle(payload []byte)
}

// NewWebhookHandler returns a new HTTP handler that handles all supported GitHub webhooks.
func NewWebhookHandler(webhookSecret string, handlers []WebhookHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		signature := r.Header.Get(signatureHeader)
		if !verifySignature(payload, signature, webhookSecret) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		// For all registered handlers, handle the webhook payload.
		// Theoretically, this payload should be handled by only one handler.
		for _, handler := range handlers {
			handler.Handle(payload)
		}

		w.WriteHeader(http.StatusOK)
	}
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
