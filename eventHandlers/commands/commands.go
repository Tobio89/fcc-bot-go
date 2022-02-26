package commands

import "github.com/bwmarrin/discordgo"

var EraseCommand = &discordgo.ApplicationCommand{
	Name:        "erase",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Erase messages in a channel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "multiple",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Description: "Specify amount to erase",
		},
	},
}

var ForceLogCommand = &discordgo.ApplicationCommand{
	Name:        "forcelog",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Force Bot to Log Something",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "message",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify log message",
			Required:    true,
		},
	},
}
var CollaborationInviteCommand = &discordgo.ApplicationCommand{
	Name:        "collabwith",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Invite Someone to Collaborate",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Type:        discordgo.ApplicationCommandOptionUser,
			Description: "Specify who to invite",
			Required:    true,
		},
	},
}
