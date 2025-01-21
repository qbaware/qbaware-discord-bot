package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	// Define the bot's commands.
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with the API latency.",
		},
		{
			Name:        "help",
			Description: "Provides information about available commands.",
		},
	}

	// Define the bot command handlers.
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			helpMessage := "Available commands:\n"
			for _, cmd := range commands {
				helpMessage += fmt.Sprintf("`%s`: %s\n", cmd.Name, cmd.Description)
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: helpMessage,
				},
			})
			if err != nil {
				log.Println("Error responding to interaction: ", err)
			}
		},
	}
)

func main() {
	// Extract the API token from the environment.
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("Invalid token. Please set the DISCORD_BOT_TOKEN environment variable.")
	}

	// Check for connectivity to the Discord APIs.
	resp, err := http.Get("https://discord.com/api/v9/invites/discord-developers")
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal("Uh oh, looks like this node is currently being blocked by Discord")
	}
	resp.Body.Close()

	// Create the bot.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// Register handlers for the supported bot commands.
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	// Create a websocket connection with Discord's servers.
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
	}

	// Register supported commands with Discord's servers.
	for _, cmd := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd)
		if err != nil {
			log.Fatalf("Cannot create '%s' command: %v", cmd.Name, err)
		}
	}

	log.Println("Bot is now running...")
	select {}
}
