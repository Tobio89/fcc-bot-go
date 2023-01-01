package main

import (
	"github.com/bwmarrin/discordgo"
)

type Events struct {
	bot *Bot
}

func (e *Events) AddEventHandlers() {
	e.bot.Session.AddHandler(e.onReady)
	e.bot.Session.AddHandler(e.onMessageSent)
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

// func (e *Events) onNewMember(s *discordgo.Session, memberJoinEvent *discordgo.GuildMemberAdd) {

// 	greeting := greetings.Opening.GetRandom()
// 	suggestion := greetings.Suggestion.GetRandom()
// 	secondSuggestion := greetings.Suggestion.GetRandom()
// 	closing := greetings.Closing.GetRandom()

// 	e.bot.SendMessageToChannel(
// 		"introductions",
// 		fmt.Sprintf("%s, %s! %s introduce yourself, tell us your coding story.\n %s check out the react-for-roles channel and let us know where you're based!\n %s", greeting, memberJoinEvent.Mention(), suggestion, secondSuggestion, closing),
// 	)
// 	e.bot.SendLog(fmt.Sprintf("New Member: " + memberJoinEvent.User.Username))
// }

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
