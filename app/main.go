package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type ChannelCfg struct {
	guild  string
	logs   string
	intros string
	rfr    string
}

type BotCfg struct {
	token string
	id    string
}

type Config struct {
	server ChannelCfg
	bot    BotCfg
}

type Bot struct {
	Session  *discordgo.Session
	Cfg      *Config
	Utils    *Utils
	Events   *Events
	Commands *Commands
}

func init() {
	err := godotenv.Load("local.env")
	if err != nil {
		panic("Could not load env file")
	}
}

func main() {

	cfg := &Config{
		server: ChannelCfg{
			guild:  "900729260531134474",
			logs:   "904628203430219787",
			intros: "1056444002754383942",
			rfr:    "915490456983457833",
		},
		bot: BotCfg{
			token: os.Getenv("TOKEN"),
			id:    "900734438416191508",
		},
	}

	dg, err := discordgo.New("Bot " + cfg.bot.token)
	if err != nil {
		fmt.Println("Error starting up:")
		panic(err)
	}

	var fccbot = &Bot{
		Session: dg,
		Cfg:     cfg,
	}
	fccbot.Utils = &Utils{bot: fccbot}
	fccbot.Events = &Events{bot: fccbot}
	fccbot.Commands = &Commands{bot: fccbot}
	// Add intents/permissions
	fccbot.Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged

	// Add handlers before open
	fccbot.Events.Initialize()
	fccbot.Commands.Initialize()

	fccbot.Start()

	// Create channel, hold it open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	fccbot.Session.Close()

}

// Start the bot
func (b *Bot) Start() {

	err := b.Session.Open() // Open the websocket
	if err != nil {
		fmt.Println("Error initialising websocket:")
		panic(err)
	} else {
		fmt.Println("FCCBot started up correctly\n(ctrl-c to exit)")
	}
}

// Send a message to the log channel
func (b *Bot) SendLog(logMessage string) {
	b.Session.ChannelMessageSend(b.Cfg.server.logs, logMessage)
}

func (b *Bot) SendMessageToChannel(channelName string, message string) {
	destChannel, err := b.Utils.GetChannelByName(channelName)
	if err != nil {
		fmt.Println("Error finding channel " + channelName)
		fmt.Println(err)
	} else {
		b.Session.ChannelMessageSend(destChannel.ID, message)
	}
}
