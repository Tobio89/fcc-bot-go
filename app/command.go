package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	bot *Bot
}

func (c *Commands) Initialize() {
	c.create()
	c.bot.Session.AddHandler(c.AdminCommandGroup)
}

func (c *Commands) create() {

	_, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, EraseCommand)
	if err != nil {
		fmt.Println("Error adding erase command:")
		fmt.Println(err)
	} else {
		fmt.Println("Erase command added")
	}
	_, err = c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, ForceLogCommand)
	if err != nil {
		fmt.Println("Error adding forcelog command:")
		fmt.Println(err)
	} else {
		fmt.Println("Forcelog command added")
	}
	// _, err = b.Session.ApplicationCommandCreate(c.bot.Cfg.id, c.bot.Cfg.er.guild, CollaborationInviteCommand)
	// if err != nil {
	// 	fmt.Println("Error adding collab invitation command:")
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Collab command added")
	// }
}

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

func (c *Commands) AdminCommandGroup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options

	interactionID := i.Interaction.ID
	interactionChannel, _ := c.bot.Utils.GetChannelByID(i.ChannelID)
	interactionMember := i.Member

	interactionMemberIsAdmin, err := c.bot.Utils.IsUserAdmin(interactionMember.User.ID)
	if err != nil {
		fmt.Println("Error on evaluating admin permissions:")
		fmt.Println(err)
	} else {
		if !interactionMemberIsAdmin {
			return
		}
	}

	switch data.Name {
	case "erase":

		if len(options) == 0 {
			c.SingleErase(i, interactionChannel, interactionID, interactionMember)
		} else {
			c.MultiErase(i, options, interactionChannel, interactionID, interactionMember)
		}

	case "forcelog":
		fmt.Println("Force Log")

		err := s.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Log made in log channel", Flags: 1 << 6},
			})

		if err != nil {
			fmt.Println("Error responding to command Forcelog")
			fmt.Println(err)
		} else {
			logString := options[0].StringValue()
			c.bot.SendLog(fmt.Sprintf("By User %s: %s", interactionMember.User.Username, logString))
			fmt.Println("Force log: ", logString)
		}
	}
}

func (c *Commands) SingleErase(i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member) {

	err := c.bot.Session.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Messages Erased", Flags: 1 << 6},
		})

	if err != nil {
		fmt.Println("Error responding to command Erase")
		fmt.Println(err)
	} else {
		fmt.Println("Trigger Erase Command")
		deleteErr := c.DeleteMessages(1, interactionChannel.ID, interactionID)
		if deleteErr != nil {
			logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
			c.bot.SendLog(logMessage)
			fmt.Println("Error deleting one message")
			fmt.Println(deleteErr)
		} else {
			logMessage := fmt.Sprintf("User %s | channel %s", interactionMember.User.Username, interactionChannel.Name)
			c.bot.SendLog(logMessage)
		}
	}

}

func (c *Commands) MultiErase(i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member) {

	eraseAmount := options[0].IntValue()
	fmt.Println(eraseAmount)
	err := c.bot.Session.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Messages Erased", Flags: 1 << 6},
		})
	if err != nil {
		fmt.Println("Error responding to command Erase")
		fmt.Println(err)
	} else {
		fmt.Println("Trigger Multiple Erase Command: ", eraseAmount)
		deleteErr := c.DeleteMessages(int(eraseAmount), interactionChannel.ID, interactionID)
		if deleteErr != nil {
			logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
			c.bot.SendLog(logMessage)
			fmt.Println("Error deleting messages")
			fmt.Println(deleteErr)
		} else {
			logMessage := fmt.Sprintf("User %s | %d messages | channel %s", interactionMember.User.Username, eraseAmount, interactionChannel.Name)
			c.bot.SendLog(logMessage)
		}

	}
}

func (c *Commands) DeleteMessages(howMany int, channel string, messageID string) error {

	messages, err := c.bot.Session.ChannelMessages(channel, howMany, messageID, "", "")
	if err != nil {
		fmt.Println("Error getting messages to delete")
		return err
	}
	var messageIDs []string

	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	messageIDs = append(messageIDs, messageID)

	err = c.bot.Session.ChannelMessagesBulkDelete(channel, messageIDs)
	if err != nil {
		return err
	}

	return nil
}
