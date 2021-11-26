package events

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/constant"
	disc "github.com/BruceJi7/fcc-bot-go/discordHelpers"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/events/onReaction"

	"github.com/bwmarrin/discordgo"
)

var MSG_TO_WATCH string = ""

func OnReady(s *discordgo.Session, _ *discordgo.Ready) {
	logMessage := fmt.Sprintf(disc.Log.Init)
	disc.SendLog(s, logMessage)
}

func OnNewMember(s *discordgo.Session, memberJoinEvent *discordgo.GuildMemberAdd) {

	r := rand.New(rand.NewSource(time.Now().Unix()))
	greeting := constant.RandomGreeting(r)
	suggestion := constant.RandomSuggestion(r)
	secondSuggestion := constant.RandomSuggestion(r)
	closing := constant.RandomClosing(r)

	botWelcomeScript := fmt.Sprintf("%s, %s! %s introduce yourself, tell us your coding story.\n %s check out the react-for-roles channel and let us know where you're based!\n %s", greeting, memberJoinEvent.Mention(), suggestion, secondSuggestion, closing)

	welcomeChannel, err := disc.GetChannelByName(s, "off-topic")
	if err != nil {
		fmt.Println("Error finding off-topic channel")
		fmt.Println(err)
	} else {
		s.ChannelMessageSend(welcomeChannel.ID, botWelcomeScript)
	}
	s.GuildMemberRoleAdd(config.GuildID, memberJoinEvent.Member.User.ID, "739417002921427084")
	logMessage := fmt.Sprintf(disc.Log.NewMember + memberJoinEvent.User.Username)
	disc.SendLog(s, logMessage)
}

func OnReactionAdded(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	} else {
		onReaction.ParseReactionAdded(s, m)
	}

}
func OnReactionRemoved(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID {
		return
	} else {
		onReaction.ParseReactionRemoved(s, m)
	}

}
