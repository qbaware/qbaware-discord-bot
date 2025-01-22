package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// SendNewReleaseNotification sends a new release notification to the specified Discord channel.
func (d *Connection) SendNewReleaseNotification(releaseNotificationChannel string, releaseName string, releaseVersion string, releaseURL string, releaseBody string) error {
	embed := &discordgo.MessageEmbed{
		Title:       "New Release: " + releaseName,
		Description: releaseBody,
		URL:         releaseURL,
		Color:       0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Version",
				Value:  releaseVersion,
				Inline: true,
			},
		},
	}

	_, err := d.s.ChannelMessageSendEmbed(releaseNotificationChannel, embed)
	if err != nil {
		return errors.New("error sending release notification: " + err.Error())
	}

	return nil
}
