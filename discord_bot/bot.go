package main

import (
	"os"
	"fmt"
	"time"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/carlescere/scheduler"

	//sql
	"github.com/jmoiron/sqlx"
	_"github.com/go-sql-driver/mysql"

	"github.com/Tayu0404/attendance_rec/discord_bot/modules"
)

//command
var (
	cmndAttendance    = "a!attendance"

	//test command
	cmndSendMessage   = "a!send"
	cmndDeleteMessage = "a!delete"
	cmndReactionadd   = "a!addreaction"
	cmndReactions     = "a!reaction"
)

//Discord Reaction
type Reaction struct {
	UserID    string
	MessageID string
	Emoji     string
	Time      time.Time
}

type ReacCheck struct {
	UserID    string
	MessageID string
	ChannelID string
	Time      time.Time
}

//Discord Reaction map
var reac = make(map[string]Reaction)
var reacCheckList = make(map[string]ReacCheck)

//db connect
var db, _ = sqlx.Connect("mysql",
	"attendance_rec:@tcp(db:3306)/attendance_rec_db")

func main() {
	var (
		Token = os.Getenv("BOT_TOKEN")
		stopBot = make(chan bool)
	)

	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Error Bot Logging in")
		log.Println(err)
	}
 
	//Discord add handler
	discord.AddHandler(onMessageCreate)
	discord.AddHandler(onMessageReactionAdd)
	
	//scheduler start
	scheduler.Every(1).Minutes().Run(reactionTimeout)

	//websocket open
	err = discord.Open()
	if err != nil {
		log.Println("Error Websocket Open : ", err)
	}

	fmt.Println("Listening...")
	<-stopBot
	return
}

//Discord WSAPI event
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%20s %20s > %s\n", time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	
	//bot measures
	if m.Author.Bot {
		return
	}

	switch {
		//bot commands
		case strings.HasPrefix(m.Content, fmt.Sprintf(cmndSendMessage)):
			sendMessage(s, m.ChannelID, "sendMessage test")
			reactionTimeout()
		case strings.HasPrefix(m.Content, fmt.Sprintf(cmndDeleteMessage)):
			deleteMessage(s,m.ChannelID, m.ID)
		case strings.HasPrefix(m.Content, fmt.Sprintf(cmndReactionadd)):
			messageReactionAdd(s, m.ChannelID, m.ID, "🍣")
		case strings.HasPrefix(m.Content, fmt.Sprintf(cmndAttendance)):
			msg := sendMessage(s, m.ChannelID, "🤒 : Sick\n😴 : Oversleeping\n💼 : Other")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "🤒")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "😴")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "💼")
			
			reacCheckList[msg.ID] = ReacCheck{m.Author.ID, msg.ID, msg.ChannelID, time.Now()}
	}
}

func onMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	fmt.Printf("%20s %20s > %s\n", time.Now().Format(time.Stamp), r.UserID, r.Emoji.Name)
	reac[r.MessageID] = Reaction{r.UserID, r.MessageID, r.Emoji.Name, time.Now()}
	fmt.Println(len(reacCheckList))
	if len(reacCheckList) != 0 {
		for k, _ := range reacCheckList {
			 reactionCheck(reacCheckList[k].UserID, reacCheckList[k].MessageID, k)
		}
	}
}

//Discord bot
func sendMessage(s *discordgo.Session, c string, msg string) *discordgo.Message {
	m, err := s.ChannelMessageSend(c, msg)
	log.Println(">>>" + msg)
	if err != nil {
		log.Println("Error sendMessage : ", err)
	}
	return m
}

func deleteMessage(s *discordgo.Session,c string, m string) {
	err := s.ChannelMessageDelete(c, m)
	if err != nil {
		log.Println("Error deleteMessage : ", err)
	}
}

func messageReactionAdd(s *discordgo.Session, c string, m string, emojiID string) {
	err := s.MessageReactionAdd(c, m, emojiID)
	if err != nil {
		log.Println("Error messageReactionAdd : ", err)
	}
}

//Reaction chack
func reactionCheck(u string, m string, key string) {
	for k, _ := range reac {
		if reac[k].MessageID != m {
			continue
		}
		if reac[k].UserID == u {
			delete(reacCheckList, key)
			switch {
				case "🤒" == reac[k].Emoji:
					t := time.Now().Format("20060102")
					module.AddToDB(db, u, t, "Sick")
				case "😴" == reac[k].Emoji:
					t := time.Now().Format("20060102")
					module.AddToDB(db, u, t, "Oversleeping")
				case "💼" == reac[k].Emoji:
					t := time.Now().Format("20060102")
					module.AddToDB(db, u, t, "Other")
			}
		}
	}
}

//Reaction timeout
func reactionTimeout() {
	for k, _ := range reac {
		t := time.Since(reac[k].Time)
		c := time.Duration(300000000000)
		fmt.Println(c <= t)
		if c <= t {
			delete(reac,k)
			fmt.Println(reac)
		}
	}

	for k, _ := range reacCheckList {
		t := time.Since(reacCheckList[k].Time)
		c := time.Duration(300000000000)
		if c <= t {
			delete(reacCheckList, k)
			fmt.Println(reacCheckList)
		}
	}
}
