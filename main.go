package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("Invalid token. Please set the DISCORD_BOT_TOKEN environment variable.")
	}

	resp, err := http.Get("https://discord.com/api/v9/invites/discord-developers")
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal("Uh oh, looks like this node is currently being blocked by Discord")
	}
	resp.Body.Close()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "ping" {
			apiLatency := s.HeartbeatLatency()
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("API Latency is %dms", apiLatency.Milliseconds()),
				},
			})
			if err != nil {
				log.Println("Error responding to interaction: ", err)
			}
		}
	})

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
	}

	log.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}
