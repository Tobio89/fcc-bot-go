package learningResources

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/constant"
	disc "github.com/BruceJi7/fcc-bot-go/discordHelpers"

	"github.com/bwmarrin/discordgo"
)

func LearningResourcePost(s *discordgo.Session, m *discordgo.MessageReactionAdd, learningDiscussionChannel *discordgo.Channel, learningResourcesChannel *discordgo.Channel) {

	message, err := s.ChannelMessage(learningDiscussionChannel.ID, m.MessageID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added")
		fmt.Println("Whilst handling learning-discussion reaction")
		fmt.Println("Error finding message")
		fmt.Println(err)
		return
	}

	hasBeenProcessed, bulbCount := parseLearningReactions(message.Reactions, constant.LearningEmoji)
	if hasBeenProcessed { // Bot already addressed this message
		return
	}

	if bulbCount >= config.LearningVoteRequirement { // If x bulbs (or more) (probably 5 lol)
		messageContents := fmt.Sprintf("Thanks, %s, who posted this resource: \n"+message.Content, message.Author.Mention())

		s.ChannelMessageSend(learningResourcesChannel.ID, messageContents)
		s.MessageReactionAdd(learningDiscussionChannel.ID, message.ID, constant.BotProcessedEmoji)
		logMessage := fmt.Sprintf(disc.Log.LearningPost+"%s's post was added to Learning Resources", message.Author)
		disc.SendLog(s, logMessage)
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
