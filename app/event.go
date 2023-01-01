package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Events struct {
	bot *Bot
}

func (e *Events) Initialize() {
	e.bot.Session.AddHandler(e.onReady)
	e.bot.Session.AddHandler(e.onMessageSent)
	e.bot.Session.AddHandler(e.onReactionAdded)
	e.bot.Session.AddHandler(e.onReactionRemoved)
}

func (e *Events) onReady(s *discordgo.Session, _ *discordgo.Ready) {
	logMessage := "Bot is here"
	e.bot.SendLog(logMessage)
}

func (e *Events) onMessageSent(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == e.bot.Cfg.bot.id {
		return
	}
	e.handleIntroductionVerification(m)
}

func (e *Events) onReactionAdded(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	} else {
		e.parseReactionAdded(m)
	}
}

func (e *Events) onReactionRemoved(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID {
		return
	} else {
		e.parseReactionRemoved(m)
	}
}

func (e *Events) handleIntroductionVerification(m *discordgo.MessageCreate) {
	if m.ChannelID != e.bot.Cfg.server.intros {
		return
	}

	member, err := e.bot.Utils.GetMemberByID(m.Author.ID)
	if err != nil {
		return
	}

	if hasRole, err := e.bot.Utils.UserHasRole(member, "verified"); err != nil || hasRole {
		return
	}

	e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, "1056454967772323861")
	e.bot.SendLog("Added verified")
}

func (e *Events) parseReactionAdded(m *discordgo.MessageReactionAdd) {
	emojiUsed := m.Emoji.MessageFormat()

	member, err := e.bot.Utils.GetMemberByID(m.UserID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added:")
		fmt.Println("Error finding user")
		fmt.Println(err)
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == e.bot.Cfg.server.rfr {
		if emojiUsed == GatherEmoji {
			e.onlineChatRoleAdd(member)
		} else {
			e.rfrAdd(member, emojiUsed)
		}
	} else {
		//If not, might be learning-related
		learningDiscussionChannel, _ := e.bot.Utils.GetChannelByName("learning-discussion")
		learningResourcesChannel, _ := e.bot.Utils.GetChannelByName("learning-resources")

		if m.ChannelID == learningDiscussionChannel.ID && emojiUsed == LearningEmoji {
			e.learningResourcePost(m, learningDiscussionChannel, learningResourcesChannel)
		}

	}
}

func (e *Events) parseReactionRemoved(m *discordgo.MessageReactionRemove) {
	emojiUsed := m.Emoji.MessageFormat()

	member, err := e.bot.Utils.GetMemberByID(m.UserID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added:")
		fmt.Println("Error finding user")
		fmt.Println(err)
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == e.bot.Cfg.server.rfr {
		if emojiUsed == GatherEmoji {
			e.onlineChatRoleRemove(member)
		} else {
			e.rfrRemove(member, emojiUsed)
		}
	}
}

func (e *Events) rfrAdd(member *discordgo.Member, emojiUsed string) {
	noLocationRole, _ := e.bot.Utils.GetRoleByName("No-Location")

	//If the role matches one of the RFR roles
	if RFRRoleSelected, exists := RFRMap[emojiUsed]; exists {

		role, err := e.bot.Utils.GetRoleByName(RFRRoleSelected)
		if err != nil {
			fmt.Println("Whilst parsing reaction added:")
			fmt.Println("Error finding role")
			fmt.Println(err)
			return
		}

		for _, userExistingRoleID := range member.Roles {
			if userExistingRoleID == role.ID {
				// They already have the role, EJECT
				return
			}
		}
		e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, role.ID)
		e.bot.SendLog(fmt.Sprintf("User %s receives role %s", member.User.Username, RFRRoleSelected))
	}

	// Check if they have the no-location role, and remove it
	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == noLocationRole.ID {
			err := e.bot.Session.GuildMemberRoleRemove(e.bot.Cfg.server.guild, member.User.ID, noLocationRole.ID)
			if err != nil {
				fmt.Println("Whilst parsing reaction removed:")
				fmt.Println("Error removing no-location role")
				fmt.Println(err)
			}
		}
	}
}

func (e *Events) rfrRemove(member *discordgo.Member, emojiUsed string) {
	noLocationRole, _ := e.bot.Utils.GetRoleByName("No-Location")

	// If the role matches one of the RFR roles
	// RFRRoleSelected == role that the reaction was for
	if RFRRoleSelected, exists := RFRMap[emojiUsed]; exists {

		// Get full role object for RFR role used
		role, err := e.bot.Utils.GetRoleByName(RFRRoleSelected)
		if err != nil {
			fmt.Println("Whilst parsing reaction removed:")
			fmt.Println("Error finding role")
			fmt.Println(err)
			return
		}

		// If the user actually has that role, remove it.
		shouldRemove := false
		for _, userExistingRoleID := range member.Roles {
			if userExistingRoleID == role.ID {
				shouldRemove = true
				break
			}
		}

		if shouldRemove {

			err = e.bot.Session.GuildMemberRoleRemove(e.bot.Cfg.server.guild, member.User.ID, role.ID)
			if err != nil {
				fmt.Println("Whilst parsing reaction removed:")
				fmt.Println("Error removing role")
				fmt.Println(err)
				return
			}
			fmt.Println("Successfully removed ", RFRRoleSelected)
			e.bot.SendLog(fmt.Sprintf("User %s loses role %s", member.User.Username, RFRRoleSelected))
		}

		// If the user has none of the RFR roles, give them 'No-Location'

		shouldAddNoLocation := true
		for _, usersRoleID := range member.Roles { // Check over the roles the user has

			roleUserHas, _ := e.bot.Utils.GetRoleByID(usersRoleID)

			if roleUserHas.Name == RFRRoleSelected {

				continue
			} else {

				// Scan the list of RFR roles for this role
				for _, RFRRole := range RFRRoles {
					// If RFR list contains the current role we're checking for
					if RFRRole == roleUserHas.Name {
						shouldAddNoLocation = false
						break
					}
				}

				if !shouldAddNoLocation {
					break
				}
			}
		}
		// If none of the location-based (RFR) roles
		// Add No-location role
		if shouldAddNoLocation {
			e.bot.SendLog(fmt.Sprintf("User %s has no location-based roles, gains No-Location", member.User.Username))
			fmt.Println("Add no location")
			e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, noLocationRole.ID)
		}

	}
}

func (e *Events) onlineChatRoleAdd(member *discordgo.Member) {

	OnlineChatSubscriptionRole, _ := e.bot.Utils.GetRoleByName(GatherRoleName)

	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == OnlineChatSubscriptionRole.ID {
			// They already have the role, EJECT
			return
		}
	}

	err := e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, OnlineChatSubscriptionRole.ID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added:")
		fmt.Println("Error removing role")
		fmt.Println(err)
	}
	e.bot.SendLog(fmt.Sprintf("User %s subscribes to Gather updates", member.User.Username))
}

func (e *Events) onlineChatRoleRemove(member *discordgo.Member) {

	OnlineChatSubscriptionRole, _ := e.bot.Utils.GetRoleByName(GatherRoleName)

	shouldRemove := false
	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == OnlineChatSubscriptionRole.ID {
			shouldRemove = true
			break
		}
	}
	if !shouldRemove {
		// Leave the function, there is no role to remove
		return
	}

	err := e.bot.Session.GuildMemberRoleRemove(e.bot.Cfg.server.guild, member.User.ID, OnlineChatSubscriptionRole.ID)
	if err != nil {
		fmt.Println("Whilst parsing reaction removed:")
		fmt.Println("Error removing role")
		fmt.Println(err)
		return
	}
	e.bot.SendLog(fmt.Sprintf("User %s removes subscription to Gather updates", member.User.Username))
}

func (e *Events) learningResourcePost(m *discordgo.MessageReactionAdd, learningDiscussionChannel *discordgo.Channel, learningResourcesChannel *discordgo.Channel) {

	message, err := e.bot.Session.ChannelMessage(learningDiscussionChannel.ID, m.MessageID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added")
		fmt.Println("Whilst handling learning-discussion reaction")
		fmt.Println("Error finding message")
		fmt.Println(err)
		return
	}

	hasBeenProcessed, bulbCount := parseLearningReactions(message.Reactions, LearningEmoji)
	if hasBeenProcessed { // Bot already addressed this message
		return
	}

	if bulbCount >= LearningVoteRequirement { // If x bulbs (or more) (probably 5 lol)
		messageContents := fmt.Sprintf("%s\nThanks, %s, who posted this resource: \n"+message.Content, e.makeMessageLink(message.Reference()), message.Author.Mention())
		e.bot.Session.ChannelMessageSend(learningResourcesChannel.ID, messageContents)
		e.bot.Session.MessageReactionAdd(learningDiscussionChannel.ID, message.ID, BotProcessedEmoji)
		e.bot.SendLog(fmt.Sprintf("%s's post was added to Learning Resources", message.Author))
	}
}

func parseLearningReactions(reactions []*discordgo.MessageReactions, emoji string) (bool, int) {

	hasBotResponded := false
	bulbCount := 0
	for _, r := range reactions {
		if r.Me {
			hasBotResponded = true
			break
		}
		if r.Emoji.MessageFormat() == emoji {
			bulbCount = r.Count
		}
	}
	return hasBotResponded, bulbCount
}

func (e *Events) makeMessageLink(reference *discordgo.MessageReference) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", e.bot.Cfg.server.guild, reference.ChannelID, reference.MessageID)
}

var RFRMap = map[string]string{
	"💺": "Seoul-Based",
	"💗": "Ulsan-Based",
	"🚌": "Busan-Based",
	"🌄": "Overseas-Based",
}

var RFRRoles = []string{
	"Seoul-Based",
	"Ulsan-Based",
	"Busan-Based",
	"Overseas-Based",
}

const BotProcessedEmoji = "✅"
const LearningEmoji = "💡"

const GatherEmoji = "🍇"
const GatherRoleName = "Gather-Attendees"

var CollabRoleMap = map[string]string{
	"project-shelf": "Project Shelf Collaborator",
}

const LearningVoteRequirement = 3
