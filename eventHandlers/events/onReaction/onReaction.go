package onReaction

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/constant"
	disc "github.com/BruceJi7/fcc-bot-go/discordHelpers"

	"github.com/BruceJi7/fcc-bot-go/eventHandlers/events/onReaction/learningResources"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/events/onReaction/reactForRole"

	"github.com/bwmarrin/discordgo"
)

func ParseReactionAdded(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	emojiUsed := m.Emoji.MessageFormat()

	if config.TestMode {
		fmt.Println("Emoji used: ", emojiUsed)
	}

	member, err := disc.FetchMember(s, m.UserID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added:")
		fmt.Println("Error finding user")
		fmt.Println(err)
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == config.RFRPostID {
		if emojiUsed == constant.GatherEmoji {
			reactForRole.OnlineChatRoleAdd(s, member)
		} else {
			reactForRole.RFRAdd(s, member, emojiUsed)
		}
	} else {
		//If not, might be learning-related
		learningDiscussionChannel, _ := disc.GetChannelByName(s, "learning-discussion")
		learningResourcesChannel, _ := disc.GetChannelByName(s, "learning-resources")

		if m.ChannelID == learningDiscussionChannel.ID && emojiUsed == constant.LearningEmoji {

			learningResources.LearningResourcePost(s, m, learningDiscussionChannel, learningResourcesChannel)

		}

	}

}

func ParseReactionRemoved(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	emojiUsed := m.Emoji.MessageFormat()

	member, err := disc.FetchMember(s, m.UserID)
	if err != nil {
		fmt.Println("Whilst parsing reaction removed:")
		fmt.Println("Error finding user")
		fmt.Println(err)
		return
	}

	// If the reaction was on the RFR Post:
	if m.MessageID == config.RFRPostID {
		if emojiUsed == constant.GatherEmoji {
			reactForRole.OnlineChatRoleRemove(s, member)
		} else {
			reactForRole.RFRRemove(s, member, emojiUsed)
		}
	}
}
