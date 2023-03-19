package main

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/app/msg"
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
	e.bot.Session.AddHandler(e.onNewMember)
}

func (e *Events) onReady(s *discordgo.Session, _ *discordgo.Ready) {
	logMessage := "Bot was turned on"
	e.bot.SendLog(msg.LogOnReady, logMessage)
}

func (e *Events) onNewMember(s *discordgo.Session, memberJoinEvent *discordgo.GuildMemberAdd) {

	greeting := msg.Opening.GetRandom()
	suggestion := msg.Suggestion.GetRandom()
	closing := msg.Closing.GetRandom()

	botWelcomeScript := fmt.Sprintf("%s, %s! Welcome to FCC Korea's discord server!\n*You'll need to introduce yourself here to complete your verification and get access to the full server :)*\n*여기서 자기소개하면 사용자 검증을 완료 될 겁니다*\nWe'd love to get to know you and find out where you are on your coding journey!\nOnce you're verified, %s check out the react-for-roles channel and let us know where you're based!\n%s", greeting, memberJoinEvent.Mention(), suggestion, closing)

	e.bot.Session.ChannelMessageSend(e.bot.Cfg.server.intros, botWelcomeScript)

	userNick := e.bot.Utils.MakeUserNickLogString(memberJoinEvent.User)
	e.bot.SendLog(msg.LogNewMember, fmt.Sprintf("User %s joined the server", userNick))
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

	if hasRole, err := e.bot.Utils.MemberHasRoleByRoleID(member, e.bot.Cfg.roles.verified); err != nil || hasRole {
		return
	}

	e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, e.bot.Cfg.roles.verified)

	userNick := e.bot.Utils.MakeUserNickLogString(member.User)
	e.bot.SendLog(msg.LogVerification, fmt.Sprintf("User %s became verified", userNick))
}

func (e *Events) parseReactionAdded(m *discordgo.MessageReactionAdd) {
	member, err := e.bot.Utils.GetMemberByID(m.UserID)
	if err != nil {
		e.bot.SendLog(msg.LogError, "Whilst parsing reaction add:")
		e.bot.SendLog(msg.LogError, err.Error())
		return
	}

	// Only verified users can use this feature
	if isVerified, _ := e.bot.Utils.MemberHasRoleByRoleID(member, e.bot.Cfg.roles.verified); !isVerified {
		return
	}

	emojiUsed := m.Emoji.MessageFormat()

	// If the reaction was on the RFR Post:
	if m.MessageID == e.bot.Cfg.server.rfr {
		if emojiUsed == OnlineMeetupEmoji {
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
		e.bot.SendLog(msg.LogError, "Whilst parsing reaction remove:")
		e.bot.SendLog(msg.LogError, err.Error())
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == e.bot.Cfg.server.rfr {
		if emojiUsed == OnlineMeetupEmoji {
			e.onlineChatRoleRemove(member)
		} else {
			e.rfrRemove(member, emojiUsed)
		}
	}
}

func (e *Events) rfrAdd(member *discordgo.Member, emojiUsed string) {

	userNick := e.bot.Utils.MakeUserNickLogString(member.User)

	//If the role matches one of the RFR roles
	if RFRRoleSelected, exists := RFRMap[emojiUsed]; exists {

		role, err := e.bot.Utils.GetRoleByID(RFRRoleSelected)
		if err != nil {
			e.bot.SendLog(msg.LogError, "Whilst parsing reaction add, getting role:")
			e.bot.SendLog(msg.LogError, err.Error())
			return
		}

		for _, userExistingRoleID := range member.Roles {
			if userExistingRoleID == role.ID {
				// They already have the role, EJECT
				return
			}
		}
		e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, role.ID)

		userNick := e.bot.Utils.MakeUserNickLogString(member.User)
		e.bot.SendLog(msg.LogRFR, fmt.Sprintf("User %s gains role %s", userNick, RFRRoleSelected))
	} else {
		e.bot.SendLog(msg.LogError, fmt.Sprintf("User %s used rando emoji: %s", userNick, emojiUsed))

		if emojiUsed == "<:srs:1065903555401240656>" {
			e.bot.SendLog(msg.LogError, "The full code (<:srs:1065903555401240656>) identifies the emoji")
		}
		if emojiUsed == ":srs:1065903555401240656" {
			e.bot.SendLog(msg.LogError, "The full code, with no <> (:srs:1065903555401240656) identifies the emoji")
		}
		if emojiUsed == ":srs:" {
			e.bot.SendLog(msg.LogError, "The code :srs: identifies the emoji")
		}
		if emojiUsed == "1065903555401240656" {
			e.bot.SendLog(msg.LogError, "The emoji's ID identifies the emoji ")
		}
	}
}

func (e *Events) rfrRemove(member *discordgo.Member, emojiUsed string) {

	// If the role matches one of the RFR roles
	// RFRRoleSelected == role that the reaction was for
	if RFRRoleSelected, exists := RFRMap[emojiUsed]; exists {

		// Get full role object for RFR role used
		role, err := e.bot.Utils.GetRoleByID(RFRRoleSelected)
		if err != nil {
			e.bot.SendLog(msg.LogError, "Whilst parsing reaction remove, getting role:")
			e.bot.SendLog(msg.LogError, err.Error())
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
				e.bot.SendLog(msg.LogError, "Whilst parsing reaction remove, removing role:")
				e.bot.SendLog(msg.LogError, err.Error())
				return
			}

			userNick := e.bot.Utils.MakeUserNickLogString(member.User)
			e.bot.SendLog(msg.LogRFR, fmt.Sprintf("User %s loses role %s", userNick, RFRRoleSelected))
		}
	}
}

func (e *Events) onlineChatRoleAdd(member *discordgo.Member) {

	OnlineChatSubscriptionRole, _ := e.bot.Utils.GetRoleByID(OnlineMeetupRoleID)

	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == OnlineChatSubscriptionRole.ID {
			// They already have the role, EJECT
			return
		}
	}

	err := e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, OnlineChatSubscriptionRole.ID)
	if err != nil {
		e.bot.SendLog(msg.LogError, "Whilst adding Gather role:")
		e.bot.SendLog(msg.LogError, err.Error())
	} else {

		userNick := e.bot.Utils.MakeUserNickLogString(member.User)
		e.bot.SendLog(msg.LogNewMember, fmt.Sprintf("User %s subscribes to Online Meetup updates", userNick))
	}
}

func (e *Events) onlineChatRoleRemove(member *discordgo.Member) {

	OnlineChatSubscriptionRole, _ := e.bot.Utils.GetRoleByID(OnlineMeetupRoleID)

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
		e.bot.SendLog(msg.LogError, "Whilst parsing gather role add:")
		e.bot.SendLog(msg.LogError, err.Error())
		return
	}
	userNick := e.bot.Utils.MakeUserNickLogString(member.User)
	e.bot.SendLog(msg.LogNewMember, fmt.Sprintf("User %s unsubscribes from Online Meetup updates", userNick))
}

func (e *Events) learningResourcePost(m *discordgo.MessageReactionAdd, learningDiscussionChannel *discordgo.Channel, learningResourcesChannel *discordgo.Channel) {

	message, err := e.bot.Session.ChannelMessage(learningDiscussionChannel.ID, m.MessageID)
	if err != nil {
		e.bot.SendLog(msg.LogError, "Whilst parsing learning resource, finding msg:")
		e.bot.SendLog(msg.LogError, err.Error())
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

		userNick := e.bot.Utils.MakeUserNickLogString(m.Member.User)
		e.bot.SendLog(msg.LogLearning, fmt.Sprintf("User %s's post was added to Learning Resources", userNick))
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
	"💺": "734321889585004595",
	"💗": "734321831095435315",
	"🚌": "781398242877112322",
	"🌄": "872308863548923974",
}

var RFRRoles = []string{
	"Seoul",
	"Ulsan",
	"Busan",
	"Overseas",
}

const BotProcessedEmoji = "✅"
const LearningEmoji = "💡"

const OnlineMeetupEmoji = "🍇"
const OnlineMeetupRoleID = "933240244596256808"

var CollabRoleMap = map[string]string{
	"project-shelf": "Project Shelf Collaborator",
}

const LearningVoteRequirement = 3
