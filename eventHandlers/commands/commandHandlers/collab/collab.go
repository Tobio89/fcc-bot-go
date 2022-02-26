package collab

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/constant"
	disc "github.com/BruceJi7/fcc-bot-go/discordHelpers"
	"github.com/bwmarrin/discordgo"
)

func CollabInvite(s *discordgo.Session, i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, options []*discordgo.ApplicationCommandInteractionDataOption) {

	// Check if command was used in correct location
	canUseCommandHere := false
	for channel := range constant.CollabRoleMap {
		if channel == interactionChannel.Name {
			canUseCommandHere = true
			break
		}
	}
	// Reject command if not in the right room
	if !canUseCommandHere {
		err := s.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "You can only use this command whilst in a private collaboration channel", Flags: 1 << 6},
			})
		if err != nil {
			fmt.Println("Error responding to command collaboration")
			fmt.Println(err)
		}
		return
	}

	targetMember, err := disc.FetchMember(s, options[0].UserValue(s).ID)

	// Get the role that matches the current channel
	collaborationRole, err := disc.GetRoleByName(s, constant.CollabRoleMap[interactionChannel.Name])
	alreadyHasRole, err := disc.UserHasRole(s, targetMember, collaborationRole.Name)
	if err != nil {
		fmt.Println("Error finding role:")
		fmt.Println(err)
	} else {

		if !alreadyHasRole {

			s.GuildMemberRoleAdd(config.GuildID, targetMember.User.ID, collaborationRole.ID)
			responseMessage := fmt.Sprintf("You invited user %s to collaborate in %s", targetMember.Nick, interactionChannel.Name)
			err := s.InteractionRespond(i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: responseMessage, Flags: 1 << 6},
				})
			if err != nil {
				fmt.Println("Error responding to command collabwith")
				fmt.Println(err)
			}
		} else {
			responseMessage := fmt.Sprintf("User %s is already a collaborator in %s", targetMember.Nick, interactionChannel.Name)
			err := s.InteractionRespond(i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: responseMessage, Flags: 1 << 6},
				})
			if err != nil {
				fmt.Println("Error responding to command collabwith")
				fmt.Println(err)
			}

		}
	}

}
