package discord

import (
	"errors"
	"fmt"
	"log"
	"net/http"

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

// Connection represents a websocket connection to the Discord API.
type Connection struct {
	s *discordgo.Session
}

// NewConnection creates a new Discord connection.
func NewConnection(token string) *Connection {
	// Create the bot.
	dg, _ := discordgo.New("Bot " + token)

	return &Connection{
		s: dg,
	}
}

// Init initializes the Discord connection and the Discord bot.
func (d *Connection) Init() error {
	// Check for connectivity to the Discord APIs.
	resp, err := http.Get("https://discord.com/api/v9/invites/discord-developers")
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("blocked by Discord")
	}
	resp.Body.Close()

	// Register handlers for the supported bot commands.
	d.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	// Create a websocket connection with Discord's servers.
	err = d.s.Open()
	if err != nil {
		return errors.New("error opening Discord session: " + err.Error())
	}

	// Register supported commands with Discord's servers.
	for _, cmd := range commands {
		_, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, "", cmd)
		if err != nil {
			return errors.New("error creating command: " + err.Error())
		}
	}

	log.Println("Bot is now running...")

	return nil
}
