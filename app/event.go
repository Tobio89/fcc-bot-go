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
	secondSuggestion := msg.Suggestion.GetRandom()
	closing := msg.Closing.GetRandom()

	botWelcomeScript := fmt.Sprintf("%s, %s! %s introduce yourself, tell us your coding story.\n %s check out the react-for-roles channel and let us know where you're based!\n %s", greeting, memberJoinEvent.Mention(), suggestion, secondSuggestion, closing)

	e.bot.Session.ChannelMessageSend(e.bot.Cfg.server.intros, botWelcomeScript)
	e.bot.SendLog(msg.LogNewMember, fmt.Sprintf(memberJoinEvent.User.Username))
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

	if hasRole, err := e.bot.Utils.UserHasRoleByRoleID(member, e.bot.Cfg.roles.verified); err != nil || hasRole {
		return
	}

	e.bot.Session.GuildMemberRoleAdd(e.bot.Cfg.server.guild, member.User.ID, e.bot.Cfg.roles.verified)
	e.bot.SendLog(msg.LogVerification, fmt.Sprintf("User %s became verified", member.User.Username))
}

func (e *Events) parseReactionAdded(m *discordgo.MessageReactionAdd) {
	emojiUsed := m.Emoji.MessageFormat()

	member, err := e.bot.Utils.GetMemberByID(m.UserID)
	if err != nil {
		e.bot.SendLog(msg.LogError, "Whilst parsing reaction add:")
		e.bot.SendLog(msg.LogError, err.Error())
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
		e.bot.SendLog(msg.LogError, "Whilst parsing reaction remove:")
		e.bot.SendLog(msg.LogError, err.Error())
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

	//If the role matches one of the RFR roles
	if RFRRoleSelected, exists := RFRMap[emojiUsed]; exists {

		role, err := e.bot.Utils.GetRoleByName(RFRRoleSelected)
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
		e.bot.SendLog(msg.LogRFR, fmt.Sprintf("User %s receives role %s", member.User.Username, RFRRoleSelected))
	}
}

func (e *Events) rfrRemove(member *discordgo.Member, emojiUsed string) {

	// If the role matches one of the RFR roles
	// RFRRoleSelected == role that the reaction was for
	if RFRRoleSelected, exists := RFRMap[emojiUsed]; exists {

		// Get full role object for RFR role used
		role, err := e.bot.Utils.GetRoleByName(RFRRoleSelected)
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
			e.bot.SendLog(msg.LogRFR, fmt.Sprintf("User %s loses role %s", member.User.Username, RFRRoleSelected))
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
		e.bot.SendLog(msg.LogError, "Whilst adding Gather role:")
		e.bot.SendLog(msg.LogError, err.Error())
	} else {
		e.bot.SendLog(msg.LogRFR, fmt.Sprintf("User %s subscribes to Gather updates", member.User.Username))
	}
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
		e.bot.SendLog(msg.LogError, "Whilst parsing gather role add:")
		e.bot.SendLog(msg.LogError, err.Error())
		return
	}
	e.bot.SendLog(msg.LogRFR, fmt.Sprintf("User %s removes subscription to Gather updates", member.User.Username))
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
		e.bot.SendLog(msg.LogLearning, fmt.Sprintf("%s's post was added to Learning Resources", message.Author))
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
