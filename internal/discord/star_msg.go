package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// SendNewStarNotification sends a new star notification to a Discord channel that matches the given repository name.
func (d *Connection) SendNewStarNotification(repoFullName string, repoURL string, starringUser string, totalStars string) error {
	channelID, ok := ChannelMapping[repoFullName]
	if !ok {
		return errors.New("no channel mapping found for repository: " + repoFullName)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "New Star! :star2:",
		Description: fmt.Sprintf("User '%s' starred the repository", starringUser),
		URL:         repoURL,
		Color:       0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Stars",
				Value:  totalStars,
				Inline: true,
			},
		},
	}

	_, err := d.s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return errors.New("error sending release notification: " + err.Error())
	}

	return nil
}
