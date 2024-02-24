package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BruceJi7/fcc-bot-go/app/msg"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type ChannelCfg struct {
	guild             string
	logs              string
	intros            string
	rfr               string
	learningResources string
}

type BotCfg struct {
	token string
	id    string
}

type Config struct {
	server   ChannelCfg
	bot      BotCfg
	roles    Roles
	meta     BotMeta
	database DatabaseCfg
}

type Bot struct {
	Session  *discordgo.Session
	Cfg      *Config
	Utils    *Utils
	Events   *Events
	Commands *Commands
}

type Roles struct {
	verified string
}

type BotMeta struct {
	startupViaCron bool
	startupTime    time.Time
}

type DatabaseCfg struct {
	dbPath string
}

type Database struct {
	conn *sql.DB
	Cfg  *DatabaseCfg
}

var cronStartupFlag *bool

func init() {
	var prodModeFlag = flag.Bool("p", false, "Use dev environment file")
	cronStartupFlag = flag.Bool("c", false, "Started up by Cron")

	flag.Parse()

	if *prodModeFlag {
		err := godotenv.Load("prod.env")
		if err != nil {
			panic("Could not load prod env file")
		} else {
			fmt.Println("Using production envs...")
		}
	} else {
		err := godotenv.Load("dev.env")
		if err != nil {
			panic("Could not load env file")
		} else {
			fmt.Println("Using development envs...")
		}
	}
}

func main() {

	cfg := &Config{
		server: ChannelCfg{
			guild:             os.Getenv("GUILD_ID"),
			logs:              os.Getenv("LOG_CHANNEL"),
			intros:            os.Getenv("INTRO_CHANNEL"),
			rfr:               os.Getenv("RFR_POST"),
			learningResources: os.Getenv("LEARNING_RESOURCE_CHANNEL"),
		},
		bot: BotCfg{
			token: os.Getenv("BOT_TOKEN"),
			id:    os.Getenv("BOT_ID"),
		},
		roles: Roles{
			verified: os.Getenv("ROLE_VERIFIED"),
		},
		meta: BotMeta{
			startupViaCron: *cronStartupFlag,
			startupTime:    time.Now(),
		},
		database: DatabaseCfg{
			dbPath: os.Getenv("DB_PATH"),
		},
	}

	db, err := sql.Open("sqlite3", cfg.database.dbPath)
	if err != nil {
		fmt.Println("Error opening database connection:")
		panic(err)
	}

	defer db.Close()

	database := &Database{
		conn: db,
		Cfg:  &cfg.database,
	}

	err = database.configureDatabase()
	if err != nil {
		fmt.Println("Error configuring database:")
		panic(err)
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
	fccbot.Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsAllWithoutPrivileged | discordgo.IntentMessageContent

	// Add handlers before open
	fccbot.Events.Initialize()
	fccbot.Commands.Initialize()

	fccbot.Start()

	fccbot.SendLog(msg.LogDatabase, fmt.Sprintf("DB setup with path: %s", cfg.database.dbPath))

	// Create channel, hold it open
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Report shutdown
	fccbot.SendLogAndPing(msg.LogShutdown, "Bot was shut down")
	// Cleanly close down the Discord session.
	fccbot.Session.Close()

}

// Start the bot
func (b *Bot) Start() {

	err := b.Session.Open() // Open the websocket
	if err != nil {
		fmt.Println("Error initialising websocket:")
		panic(err)
	}
}

// Send a message to the log channel
func (b *Bot) SendLog(logPrefix, logMessage string) {

	loc, _ := time.LoadLocation("Asia/Seoul")
	t := time.Now().In(loc)
	formattedTime := fmt.Sprintf("[%d/%02d/%02d T%02d:%02d:%02d]",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	msgString := fmt.Sprintf("%s %s %s", logPrefix, formattedTime, logMessage)
	b.Session.ChannelMessageSend(b.Cfg.server.logs, fmt.Sprintf("`%s`", msgString))
	fmt.Println(msgString)
}

func (b *Bot) SendLogAndPing(logPrefix, logMessage string) {
	ping, err := b.Utils.BotLogPing()
	if err != nil {
		b.SendLog(msg.LogError, err.Error())
	}

	b.Session.ChannelMessageSend(b.Cfg.server.logs, fmt.Sprintf("FYI %s:", ping))
	b.SendLog(logPrefix, logMessage)
}

func (b *Bot) SendMessageToChannel(channelName string, message string) {
	destChannel, err := b.Utils.GetChannelByName(channelName)
	if err != nil {
		b.SendLog(msg.LogError, "Whilst sending msg to channel:")
		b.SendLog(msg.LogError, err.Error())
	} else {
		b.Session.ChannelMessageSend(destChannel.ID, message)
	}
}

func (d *Database) configureDatabase() error {
	pragmaConfig := `
		PRAGMA busy_timeout = 5000;
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;
	`
	_, err := d.conn.Exec(pragmaConfig)
	if err != nil {
		return fmt.Errorf("failed configuring database: %w", err)
	}
	return nil
}
