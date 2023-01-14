package main

import (
	"errors"

	"github.com/BruceJi7/fcc-bot-go/app/msg"
	"github.com/bwmarrin/discordgo"
)

type Utils struct {
	bot *Bot
}

// Return discord channel struct from channel's name string
func (u *Utils) GetChannelByName(name string) (c *discordgo.Channel, err error) {
	channels, _ := u.bot.Session.GuildChannels(u.bot.Cfg.server.guild)
	for _, c := range channels {
		if c.Name == name {
			return c, nil
		}
	}
	return channels[0], errors.New("channel not found")
}

func (u *Utils) GetChannelByID(cID string) (c *discordgo.Channel, err error) {

	channels, _ := u.bot.Session.GuildChannels(u.bot.Cfg.server.guild)
	for _, c := range channels {
		if c.ID == cID {
			return c, nil
		}
	}
	return channels[0], errors.New("channel not found")
}

// Return discord role struct from role's name string
func (u *Utils) GetRoleByName(roleName string) (role *discordgo.Role, err error) {
	roles, _ := u.bot.Session.GuildRoles(u.bot.Cfg.server.guild)

	for _, role := range roles {
		if role.Name == roleName {
			return role, nil
		}
	}
	return roles[0], errors.New("role not found")
}

func (u *Utils) GetRoleByID(roleID string) (role *discordgo.Role, err error) {
	roles, _ := u.bot.Session.GuildRoles(u.bot.Cfg.server.guild)

	for _, role := range roles {
		if role.ID == roleID {
			return role, nil
		}
	}
	return roles[0], errors.New("role not found")
}

// Return boolean: does user have role, from role's name string
func (u *Utils) UserHasRoleByRoleName(member *discordgo.Member, roleToFind string) (bool, error) {

	role, err := u.GetRoleByName(roleToFind)
	if err != nil {
		return false, err
	}

	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == role.ID {
			return true, nil
		}
	}
	return false, nil
}

// Return boolean: does user have role, from role's ID string
func (u *Utils) UserHasRoleByRoleID(member *discordgo.Member, roleToFind string) (bool, error) {

	for _, userExistingRoleID := range member.Roles {
		if userExistingRoleID == roleToFind {
			return true, nil
		}
	}
	return false, nil
}

// Return discord member struct from user's ID string
func (u *Utils) GetMemberByID(userDetails string) (member *discordgo.Member, err error) {
	guildMembers, err := u.bot.Session.GuildMembers(u.bot.Cfg.server.guild, "", 1000)

	if err != nil {
		u.bot.SendLog(msg.LogError, "Whilst fetching member by ID:")
		u.bot.SendLog(msg.LogError, err.Error())
	}
	for _, member := range guildMembers {
		if member.User.ID == userDetails {
			return member, nil
		}
	}
	return guildMembers[0], errors.New("member not found")
}

func (u *Utils) MemberHasPermission(userID string, permission int64) (bool, error) {

	member, err := u.bot.Session.State.Member(u.bot.Cfg.server.guild, userID)
	if err != nil {
		if member, err = u.bot.Session.GuildMember(u.bot.Cfg.server.guild, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := u.bot.Session.State.Role(u.bot.Cfg.server.guild, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

func (u *Utils) IsUserAdmin(userID string) (bool, error) {
	return u.MemberHasPermission(userID, discordgo.PermissionAdministrator)
}
