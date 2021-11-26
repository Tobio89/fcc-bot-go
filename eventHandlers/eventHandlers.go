package eventHandlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/BruceJi7/fcc-bot-go/config"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/commands"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/commands/commandHandlers"
	"github.com/BruceJi7/fcc-bot-go/eventHandlers/events"
)

func AddEventHandlers(dg *discordgo.Session) {

	dg.AddHandler(events.OnReady)
	dg.AddHandler(events.OnNewMember)
	dg.AddHandler(events.OnReactionAdded)
	dg.AddHandler(events.OnReactionRemoved)

	dg.AddHandler(commandHandlers.AdminCommands)

}

func CreateCommands(dg *discordgo.Session) {

	_, err := dg.ApplicationCommandCreate(config.AppID, config.GuildID, commands.EraseCommand)
	if err != nil {
		fmt.Println("Error adding erase command:")
		fmt.Println(err)
	} else {
		fmt.Println("Erase command added")
	}
	_, err = dg.ApplicationCommandCreate(config.AppID, config.GuildID, commands.ForceLogCommand)
	if err != nil {
		fmt.Println("Error adding forcelog command:")
		fmt.Println(err)
	} else {
		fmt.Println("Forcelog command added")
	}

}
