package main

import (
	"fmt"
	"net/url"

	"github.com/BruceJi7/fcc-bot-go/app/msg"
	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	bot *Bot
}

func (c *Commands) Initialize() {
	c.create()
	c.bot.Session.AddHandler(c.AdminCommandGroup)
	c.bot.Session.AddHandler(c.RegularCommandGroup)
}

func (c *Commands) create() {

	allSuccessful := true

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, EraseCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding erase command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, StealthEraseCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding erase-quietly command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, IntroCleanupCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding intro-cleanup command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, ForceLogCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding forcelog command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, LearningResourceCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding learning resource command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, RemindCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding remind command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, VerifyCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding verify command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, DeVerifyCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding DeVerify command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if allSuccessful {
		c.bot.SendLog(msg.LogOnReady, "All commands successfully added")
	}
}

var EraseCommand = &discordgo.ApplicationCommand{
	Name:        "erase",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Erase messages in a channel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "reason",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify reason for erasing",
			Required:    true,
		},
		{
			Name:        "starting-point",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify starting post ID",
			Required:    true,
		},
		{
			Name:        "ending-point",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify post ID to erase until, or 'end' to erase til end",
		},
	},
}

var StealthEraseCommand = &discordgo.ApplicationCommand{
	Name:        "erase-quietly",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Erase messages in a channel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "starting-point",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify starting post id",
			Required:    true,
		},
		{
			Name:        "ending-point",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify post ID to erase until, or 'end' to erase til end",
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

var IntroCleanupCommand = &discordgo.ApplicationCommand{
	Name:        "clear-introductions",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Bot clears Bot-made posts in Introductions Channel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "starting-point",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify starting post ID",
			Required:    true,
		},
	},
}

var LearningResourceCommand = &discordgo.ApplicationCommand{
	Name:        "learning-resource",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Submit a useful learning resource",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "resource-url",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "A valid url",
			Required:    true,
		},
		{
			Name:        "resource-description",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "A description of the resource. What language is it for? What can we learn from it?",
			Required:    true,
		},
	},
}

var RemindCommand = &discordgo.ApplicationCommand{
	Name:        "remind",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Post a reminder of how to do something:",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "topic",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "what to remind the user of",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "learning",
					Value: "learning",
				},
				{
					Name:  "nickname",
					Value: "nickname",
				},
			},
		},
	},
}
var VerifyCommand = &discordgo.ApplicationCommand{
	Name:        "verify",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Manually add verified role to user:",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Type:        discordgo.ApplicationCommandOptionUser,
			Description: "who to verify",
			Required:    true,
		},
	},
}

var DeVerifyCommand = &discordgo.ApplicationCommand{
	Name:        "deverify",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Manually remove verified role from user:",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Type:        discordgo.ApplicationCommandOptionUser,
			Description: "who to de-verify",
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
		c.bot.SendLog(msg.LogError, "Whilst evaluating admin privileges:")
		c.bot.SendLog(msg.LogError, err.Error())
		return
	} else {
		if !interactionMemberIsAdmin {
			return
		}
	}

	switch data.Name {
	case "erase", "erase-quietly":

		reason := ""
		startingPostID := ""
		untilPostID := ""

		for _, opt := range options {
			if opt.Name == "reason" {
				reason = opt.StringValue()
			} else if opt.Name == "starting-point" {
				startingPostID = opt.StringValue()
			} else if opt.Name == "ending-point" {
				if opt.StringValue() == "end" {
					untilPostID = interactionID
				} else {
					untilPostID = opt.StringValue()
				}
			}
		}

		if untilPostID != "" {
			if reason != "" {
				c.MultiEraseWithReason(i, interactionChannel, interactionID, interactionMember, startingPostID, untilPostID, reason)
			} else {
				c.MultiEraseNoReason(i, interactionChannel, interactionID, interactionMember, startingPostID, untilPostID)
			}
		} else {
			if reason != "" {
				c.SingleEraseWithReason(i, interactionChannel, interactionID, startingPostID, interactionMember, reason)
			} else {
				c.SingleEraseNoReason(i, interactionChannel, interactionID, startingPostID, interactionMember)
			}
		}

	case "forcelog":

		c.bot.Utils.SendResponse(i, "Log made in Log Channel")

		logString := options[0].StringValue()
		c.bot.SendLog(msg.CommandForceLog, fmt.Sprintf("By User %s: %s", interactionMember.User.Username, logString))

	case "remind":
		switch options[0].StringValue() {
		case "nickname":
			{
				response := fmt.Sprintf("From  %s <3", interactionMember.Mention())
				c.bot.Session.ChannelMessageSend(interactionChannel.ID, response)
				response = "Check out this post to find out how to change your nickname!:\nhttps://discord.com/channels/726648668907765842/770532252768927745/955293674152009748"
				c.bot.Session.ChannelMessageSend(interactionChannel.ID, response)
				c.bot.SendLog(msg.CommandRemind, fmt.Sprintf("User %s requested reminder '%s'", c.bot.Utils.GetMemberNickOrUsername(*interactionMember), options[0].StringValue()))
				c.bot.Utils.SendResponse(i, "Remind message sent")

			}
		case "learning":
			{
				response := fmt.Sprintf("From  %s <3", interactionMember.Mention())
				c.bot.Session.ChannelMessageSend(interactionChannel.ID, response)
				response = "Type `/learning-resource` to make a post to the #learning-resources channel!\nYou'll need to add the URL of the resource, and a description of it."
				c.bot.Session.ChannelMessageSend(interactionChannel.ID, response)
				c.bot.SendLog(msg.CommandRemind, fmt.Sprintf("User %v requested reminder '%s'", c.bot.Utils.GetMemberNickOrUsername(*interactionMember), options[0].StringValue()))
				c.bot.Utils.SendResponse(i, "Remind message sent")

			}
		default:
			c.bot.SendLog(msg.LogError, fmt.Sprintf("User %v requested reminder '%s'", c.bot.Utils.GetMemberNickOrUsername(*interactionMember), options[0].StringValue()))
			c.bot.Utils.SendResponse(i, "That remind topic doesn't exist")

		}

	case "verify":
		{
			c.ManualVerify(i, options[0].UserValue(c.bot.Session))
		}

	case "deverify":
		{
			c.ManualDeVerify(i, options[0].UserValue(c.bot.Session))
		}

	case "clear-introductions":
		{
			startingPostID := options[0].StringValue()
			c.ClearIntroductions(i, startingPostID)
		}
	}

}

func (c *Commands) RegularCommandGroup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options
	interactionMember := i.Member

	switch data.Name {
	case "learning-resource":
		resourceUrl := options[0].StringValue()
		resourceDescription := options[1].StringValue()

		if _, err := url.ParseRequestURI(resourceUrl); err != nil {
			c.bot.Utils.SendResponse(i, "Whoops! Looks like your URL was invalid")

		} else {

			messageContents := fmt.Sprintf("Thanks, %s, who posted this resource:\n%s\nDescription: %s", interactionMember.User.Mention(), resourceUrl, resourceDescription)
			c.bot.Session.ChannelMessageSend(c.bot.Cfg.server.learningResources, messageContents)
			c.bot.SendLog(msg.LogLearning, fmt.Sprintf("%s submitted a Learning Resource via the bot", interactionMember.User.Username))

			c.bot.Utils.SendResponse(i, "Thanks for posting a Learning Resource!")
		}

	}
}

func (c *Commands) SingleEraseNoReason(i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID, startingPostID string, interactionMember *discordgo.Member) {

	c.bot.Utils.SendResponse(i, "Message erased")

	{
		deleteErr := c.DeleteSingleMessage(interactionChannel.ID, interactionID, startingPostID)
		if deleteErr != nil {
			c.bot.SendLog(msg.LogError, "Whilst attempting to delete:")
			logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
			c.bot.SendLog(msg.LogError, logMessage)
		} else {
			logMessage := fmt.Sprintf("User %s | channel %s | quiet", interactionMember.User.Username, interactionChannel.Name)
			c.bot.SendLog(msg.CommandErase, logMessage)
		}
	}
}

func (c *Commands) SingleEraseWithReason(i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID, startingPostID string, interactionMember *discordgo.Member, reason string) {

	deleteErr := c.DeleteSingleMessage(interactionChannel.ID, interactionID, startingPostID)
	if deleteErr != nil {
		c.bot.SendLog(msg.LogError, "Whilst attempting to delete:")
		logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
		c.bot.SendLog(msg.LogError, logMessage)
	} else {
		logMessage := fmt.Sprintf("User %s | channel %s | reason \"%s\"", interactionMember.User.Username, interactionChannel.Name, reason)
		c.bot.SendLog(msg.CommandErase, logMessage)
		eraseReasonMessage := fmt.Sprintf("User %s erased messages in this channel, giving the reason:\n*%s*", interactionMember.User.Username, reason)
		c.bot.Session.ChannelMessageSend(interactionChannel.ID, eraseReasonMessage)
	}

	responseContent := "Message Erased"
	if deleteErr != nil {
		responseContent = "Whoops! Failed to erase messages. See log channel for more information"
	}

	c.bot.Utils.SendResponse(i, responseContent)

}

func (c *Commands) MultiEraseNoReason(i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member, eraseFromStartingPostID, eraseUntilPostID string) {

	deleteErr := c.DeleteMultipleMessages(eraseFromStartingPostID, eraseUntilPostID, interactionChannel.ID)
	if deleteErr != nil {
		logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
		c.bot.SendLog(msg.LogError, "Whilst attempting to delete:")
		c.bot.SendLog(msg.LogError, logMessage)
	} else {
		logMessage := fmt.Sprintf("User %s | channel %s | quiet", interactionMember.User.Username, interactionChannel.Name)
		c.bot.SendLog(msg.CommandErase, logMessage)
	}

	responseContent := "Messages Erased"
	if deleteErr != nil {
		responseContent = "Whoops! Failed to erase messages. See log channel for more information"
	}

	c.bot.Utils.SendResponse(i, responseContent)

}

func (c *Commands) MultiEraseWithReason(i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member, eraseFromStartingPostID, eraseUntilPostID, reason string) {

	deleteErr := c.DeleteMultipleMessages(eraseFromStartingPostID, eraseUntilPostID, interactionChannel.ID)
	if deleteErr != nil {
		logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
		c.bot.SendLog(msg.LogError, "Whilst attempting to delete:")
		c.bot.SendLog(msg.LogError, logMessage)
	} else {
		logMessage := fmt.Sprintf("User %s | channel %s | reason \"%s\"", interactionMember.User.Username, interactionChannel.Name, reason)
		c.bot.SendLog(msg.CommandErase, logMessage)
		eraseReasonMessage := fmt.Sprintf("User %s erased messages in this channel, giving the reason:\n*%s*", interactionMember.User.Username, reason)
		c.bot.Session.ChannelMessageSend(interactionChannel.ID, eraseReasonMessage)
	}

	responseContent := "Messages Erased"
	if deleteErr != nil {
		responseContent = "Whoops! Failed to erase messages. See log channel for more information"
	}

	c.bot.Utils.SendResponse(i, responseContent)

}

func (c *Commands) ClearIntroductions(i *discordgo.InteractionCreate, startingID string) {

	if i.ChannelID != c.bot.Cfg.server.intros {
		c.bot.SendLog(msg.LogError, "Attempted introduction clear in wrong channel")
		c.bot.Utils.SendResponse(i, "Use this command in introductions channel.")
		return
	}

	messages, err := c.bot.Session.ChannelMessages(c.bot.Cfg.server.intros, 100, "", startingID, "")
	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst attempting to clear introductions:")
		c.bot.SendLog(msg.LogError, err.Error())
		c.bot.Utils.SendResponse(i, "Failed to clear introductions. See log channel for more information")
		return
	}
	messageIDs := []string{}

	for _, m := range messages {
		if m.Author.ID == c.bot.Cfg.bot.id {
			messageIDs = append(messageIDs, m.ID)
		}
	}

	err = c.bot.Session.ChannelMessagesBulkDelete(c.bot.Cfg.server.intros, messageIDs)
	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst attempting to clear introductions:")
		c.bot.SendLog(msg.LogError, err.Error())
		c.bot.Utils.SendResponse(i, "Failed to clear introductions. See log channel for more information")

		return
	}

	responseContent := fmt.Sprintf("%d introduction messages cleared.", len(messageIDs))
	c.bot.Utils.SendResponse(i, responseContent)
	c.bot.SendLog(msg.CommandClearIntros, responseContent)
}

func (c *Commands) ManualVerify(i *discordgo.InteractionCreate, user *discordgo.User) {

	var responseContent string
	var userNick string = c.bot.Utils.MakeUserNickLogString(user)

	// Get member from server's member list
	member, err := c.bot.Utils.GetMemberByID(user.ID)
	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst attempting manual verification:")
		c.bot.SendLog(msg.LogError, err.Error())
		responseContent = fmt.Sprintf("Failed to get member %s: have they left?", userNick)

		c.bot.Utils.SendResponse(i, responseContent)

		return
	}
	// Check if they have the role already
	if hasRole, err := c.bot.Utils.MemberHasRoleByRoleID(member, c.bot.Cfg.roles.verified); err != nil || hasRole {
		responseContent = fmt.Sprintf("User %s is already verified", userNick)
		c.bot.SendLog(msg.CommandVerify, responseContent)
		c.bot.Utils.SendResponse(i, responseContent)

		return
	}

	c.bot.Session.GuildMemberRoleAdd(c.bot.Cfg.server.guild, member.User.ID, c.bot.Cfg.roles.verified)
	c.bot.SendLog(msg.CommandVerify, fmt.Sprintf("User %s became verified manually", userNick))

	responseContent = fmt.Sprintf("User %s verified", userNick)
	c.bot.Utils.SendResponse(i, responseContent)
}

func (c *Commands) ManualDeVerify(i *discordgo.InteractionCreate, user *discordgo.User) {

	var responseContent string
	var userNick string = c.bot.Utils.MakeUserNickLogString(user)

	// Get member from server's member list
	member, err := c.bot.Utils.GetMemberByID(user.ID)
	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst attempting manual de-verification:")
		c.bot.SendLog(msg.LogError, err.Error())

		responseContent = fmt.Sprintf("Failed to get member %s: have they left?", userNick)
		c.bot.Utils.SendResponse(i, responseContent)

		return
	}
	// Check if they have the role already
	if hasRole, err := c.bot.Utils.MemberHasRoleByRoleID(member, c.bot.Cfg.roles.verified); err != nil || !hasRole {
		responseContent = fmt.Sprintf("User %s is already not verified", userNick)
		c.bot.SendLog(msg.CommandDeVerify, fmt.Sprintf("User %s is already not verified", userNick))
		c.bot.Utils.SendResponse(i, responseContent)

		return
	}

	c.bot.Session.GuildMemberRoleRemove(c.bot.Cfg.server.guild, member.User.ID, c.bot.Cfg.roles.verified)
	c.bot.SendLog(msg.CommandDeVerify, fmt.Sprintf("User %s became de-verified", userNick))

	responseContent = fmt.Sprintf("User %s de-verified", userNick)
	c.bot.Utils.SendResponse(i, responseContent)
}

// Command Utilities

func (c *Commands) DeleteSingleMessage(channel, messageID, targetMessageID string) error {

	var messageIDs []string

	messageIDs = append(messageIDs, messageID)
	messageIDs = append(messageIDs, targetMessageID)

	err := c.bot.Session.ChannelMessagesBulkDelete(channel, messageIDs)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) DeleteMultipleMessages(eraseFromStartingPostID, eraseUntilPostID, channel string) error {

	messages, err := c.bot.Session.ChannelMessages(channel, 100, "", eraseFromStartingPostID, "")
	if err != nil {
		return err
	}
	messageIDs := []string{}

	for i := len(messages) - 1; i > 0; i-- {
		m := messages[i]
		messageIDs = append(messageIDs, m.ID)
		if m.ID == eraseUntilPostID {
			break
		}
	}
	messageIDs = append(messageIDs, eraseFromStartingPostID)

	err = c.bot.Session.ChannelMessagesBulkDelete(channel, messageIDs)
	if err != nil {
		return err
	}

	return nil
}
