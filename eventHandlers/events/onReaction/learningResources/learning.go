package learningResources

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/constant"

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
		messageContents := message.Content
		s.ChannelMessageSend(learningResourcesChannel.ID, messageContents)
		s.MessageReactionAdd(learningDiscussionChannel.ID, message.ID, constant.BotProcessedEmoji)

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
