package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// SendNewReleaseNotification sends a new release notification to a Discord channel that matches the given repository name.
func (d *Connection) SendNewReleaseNotification(repoFullName string, releaseName string, releaseVersion string, releaseURL string, releaseBody string) error {
	channelID, ok := ChannelMapping[repoFullName]
	if !ok {
		return errors.New("no channel mapping found for repository: " + repoFullName)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "New Release :tada:",
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

	_, err := d.s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return errors.New("error sending release notification: " + err.Error())
	}

	return nil
}
