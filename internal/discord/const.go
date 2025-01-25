package discord

// Repository constants.
const (
	RepoHomeassistantEldom   = "qbaware/homeassistant-eldom"
	RepoPyeldom              = "qbaware/pyeldom"
	RepoNilawayAction        = "qbaware/nilaway-action"
	RepoRenderRedeployAction = "qbaware/render-redeploy-action"
	RepoQbawareDiscordBot    = "qbaware/qbaware-discord-bot"
)

// Channel ID constants.
const (
	ChannelHomeassistantEldom   = "1312386863947972720"
	ChannelPyeldom              = "1312386488075288596"
	ChannelNilawayAction        = "1312387091744686080"
	ChannelRenderRedeployAction = "1312386976363839510"
	ChannelQbawareDiscordBot    = "1331855607539699713"
)

// ChannelMapping contains the mapping of repository names to their respective Discord channels.
var ChannelMapping = map[string]string{
	RepoHomeassistantEldom:   ChannelHomeassistantEldom,
	RepoPyeldom:              ChannelPyeldom,
	RepoNilawayAction:        ChannelNilawayAction,
	RepoRenderRedeployAction: ChannelRenderRedeployAction,
	RepoQbawareDiscordBot:    ChannelQbawareDiscordBot,
}
