package reactForRole

import (
	"fmt"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/constant"
	disc "github.com/BruceJi7/fcc-bot-go/discordHelpers"

	"github.com/bwmarrin/discordgo"
)

func RFRAdd(s *discordgo.Session, member *discordgo.Member, emojiUsed string) {
	noLocationRole, _ := disc.GetRoleByName(s, "No-Location")

	//If the role matches one of the RFR roles
	if RFRRoleSelected, exists := constant.RFRMap[emojiUsed]; exists {

		role, err := disc.GetRoleByName(s, RFRRoleSelected)
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
		s.GuildMemberRoleAdd(config.GuildID, member.User.ID, role.ID)
		logMessage := fmt.Sprintf("%s User %s receives role %s", disc.Log.RoleAdded, member.User.Username, RFRRoleSelected)
		disc.SendLog(s, logMessage)

	}

	// Check if they have the no-location role, and remove it
	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == noLocationRole.ID {
			err := s.GuildMemberRoleRemove(config.GuildID, member.User.ID, noLocationRole.ID)
			if err != nil {
				fmt.Println("Whilst parsing reaction removed:")
				fmt.Println("Error removing no-location role")
				fmt.Println(err)
			}

		}
	}
}

func RFRRemove(s *discordgo.Session, member *discordgo.Member, emojiUsed string) {
	noLocationRole, _ := disc.GetRoleByName(s, "No-Location")

	// If the role matches one of the RFR roles
	// RFRRoleSelected == role that the reaction was for
	if RFRRoleSelected, exists := constant.RFRMap[emojiUsed]; exists {

		// Get full role object for RFR role used
		role, err := disc.GetRoleByName(s, RFRRoleSelected)
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

			err = s.GuildMemberRoleRemove(config.GuildID, member.User.ID, role.ID)
			if err != nil {
				fmt.Println("Whilst parsing reaction removed:")
				fmt.Println("Error removing role")
				fmt.Println(err)
				return
			}
			fmt.Println("Successfully removed ", RFRRoleSelected)
			logMessage := fmt.Sprintf("%s User %s loses role %s", disc.Log.RoleAdded, member.User.Username, RFRRoleSelected)
			disc.SendLog(s, logMessage)
		}

		// If the user has none of the RFR roles, give them 'No-Location'

		shouldAddNoLocation := true
		for _, usersRoleID := range member.Roles { // Check over the roles the user has

			roleUserHas, _ := disc.GetRoleByID(s, usersRoleID)

			if roleUserHas.Name == RFRRoleSelected {

				continue
			} else {

				// Scan the list of RFR roles for this role
				for _, RFRRole := range constant.RFRRoles {
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
			logMessage := fmt.Sprintf("%s User %s has no location-based roles, gains No-Location", disc.Log.RoleAdded, member.User.Username)
			disc.SendLog(s, logMessage)
			fmt.Println("Add no location")
			s.GuildMemberRoleAdd(config.GuildID, member.User.ID, noLocationRole.ID)
		}

	}
}

func OnlineChatRoleAdd(s *discordgo.Session, member *discordgo.Member) {

	OnlineChatSubscriptionRole, _ := disc.GetRoleByName(s, constant.GatherRoleName)

	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == OnlineChatSubscriptionRole.ID {
			// They already have the role, EJECT
			return
		}
	}

	err := s.GuildMemberRoleAdd(config.GuildID, member.User.ID, OnlineChatSubscriptionRole.ID)
	if err != nil {
		fmt.Println("Whilst parsing reaction added:")
		fmt.Println("Error removing role")
		fmt.Println(err)
	}
	logMessage := fmt.Sprintf("%s User %s subscribes to Gather updates", disc.Log.RoleAdded, member.User.Username)
	disc.SendLog(s, logMessage)

}

func OnlineChatRoleRemove(s *discordgo.Session, member *discordgo.Member) {

	OnlineChatSubscriptionRole, _ := disc.GetRoleByName(s, constant.GatherRoleName)

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

	err := s.GuildMemberRoleRemove(config.GuildID, member.User.ID, OnlineChatSubscriptionRole.ID)
	if err != nil {
		fmt.Println("Whilst parsing reaction removed:")
		fmt.Println("Error removing role")
		fmt.Println(err)
		return
	}
	logMessage := fmt.Sprintf("%s User %s removes subscription to Gather updates", disc.Log.RoleAdded, member.User.Username)
	disc.SendLog(s, logMessage)
}
