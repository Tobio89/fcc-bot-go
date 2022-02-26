package commandHandlers

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/config"
	disc "github.com/BruceJi7/fcc-bot-go/discordHelpers"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/commands/commandHandlers/collab"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/commands/commandHandlers/erase"

	"github.com/bwmarrin/discordgo"
)

func AdminCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options

	interactionID := i.Interaction.ID
	interactionChannel, _ := disc.GetChannelByID(s, i.ChannelID)
	interactionMember := i.Member

	interactionMemberIsAdmin, err := disc.IsAdmin(s, config.GuildID, interactionMember.User.ID)
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
			// Triggered single erase mode

			erase.SingleErase(s, i, interactionChannel, interactionID, interactionMember)

		} else {
			// Multiple erase mode:

			erase.MultiErase(s, i, options, interactionChannel, interactionID, interactionMember)

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
			logmessage := fmt.Sprintf(disc.Log.Forcelog+"By User %s: %s", interactionMember.User.Username, logString)
			disc.SendLog(s, logmessage)
			fmt.Println("Force log: ", logString)
		}
	}
}
func CollabCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	interactionChannel, _ := disc.GetChannelByID(s, i.ChannelID)

	data := i.ApplicationCommandData()
	switch data.Name {
	case "collabwith":
		collab.CollabInvite(s, i, interactionChannel, data.Options)
	}
}
