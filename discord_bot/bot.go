package main

import (
	"os"
	"fmt"
	"time"
	"log"
	"strings"
	"strconv"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/carlescere/scheduler"

	//sql
	"github.com/jmoiron/sqlx"
	_"github.com/go-sql-driver/mysql"

	"github.com/Tayu0404/attendance_rec/discord_bot/modules"
	"github.com/Tayu0404/attendance_rec/discord_bot/calculation"
)

//command
var (
	cmndAttendance     = "a!attendance"
	cmndAttendanceRate = "a!rate"

	//test command
	cmndSendMessage    = "a!send"
	cmndDeleteMessage  = "a!delete"
	cmndReactionadd    = "a!addreaction"
	cmndReactions      = "a!reaction"
	cmndKuramubonRate  = "a!testk"
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
	Turget    string
	Day       string
	CaseID    string
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
	
	//scheduler.Every(2).Seconds().Run(calculation.Regularly)

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
	//option
	ou := false //-u
	od := false //-d

	if strings.Contains(m.Content, "-u"){
		ou = true
	}
	if strings.Contains(m.Content, "-d"){
		od = true
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
			var tui string
			var tsd string
			if ou {
				re := regexp.MustCompile(`-u\D*\d{18}`)
				ui := re.Find([]byte(m.Content))
				fmt.Println("ui:",string(ui))
				re = regexp.MustCompile(`\d{18}`)
				tui = string(re.Find(ui))
				fmt.Println("tui:",string(tui))
				mem, _ := s.GuildMembers(m.GuildID,"",1000)
				sf := true
				for i, _ := range mem {
					if mem[i].User.ID == tui {
						uc, _ := s.User(tui)
						if uc.Bot {
							return
						}
						sf = false
						continue
					}
				}
				if tui == "" {
					sendMessage(s, m.ChannelID, "Invalid argument")
					return
				}

				if sf {
					sendMessage(s, m.ChannelID, "no search user")
					return
				}
			}
			if od {
				re := regexp.MustCompile(`-d.*?\d{4}\d{1,2}\d{1,2}`)
				sd := re.Find([]byte(m.Content))
				re = regexp.MustCompile(`\d{4}\d{1,2}\d{1,2}`)
				tsd = string(re.Find(sd))
				fmt.Println(tsd)
				if tsd == "" {
					sendMessage(s, m.ChannelID, "Invalid argument")
					return
				}
				year, _ := strconv.Atoi(tsd[:4])
				month, _ := strconv.Atoi(tsd[4:6])
				day, _ := strconv.Atoi(tsd[6:8])
				err := isExist(year, month, day)
				if err != nil {
					sendMessage(s, m.ChannelID, "Invalid argument")
					return
				}
			}
			msg := sendMessage(s, m.ChannelID, "🤒 : Sick\n😴 : Oversleeping\n🚃 : Train delay\n💼 : Other\n🏫 : Official absence")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "🤒")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "😴")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "🚃")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "💼")
			messageReactionAdd(s, msg.ChannelID, msg.ID, "🏫")
			
			reacCheckList[msg.ID] = ReacCheck{m.Author.ID, msg.ID, msg.ChannelID, tui, tsd, "attendance", time.Now()}
		
		case strings.HasPrefix(m.Content, fmt.Sprintf(cmndAttendanceRate)):
			sd, ad, ar:= calculation.AttendanceRate(db, m.Author.ID)
			msg := fmt.Sprintf("<@%s> Attendance Rate \n School days : %d \n Absent days : %d \n Attendance rate : %g%%", m.Author.ID, sd, ad, ar)
			sendMessage(s, m.ChannelID, msg)
		
		case strings.HasPrefix(m.Content, fmt.Sprintf(cmndKuramubonRate)):
			sd, ad, ar := calculation.AttendanceRate(db, "269793922740518913")
			msg := fmt.Sprintf("Tester K Attendance Rate \n School days : %d \n Absent days : %d \n Attendance rate : %g%%", sd, ad, ar)
			sendMessage(s, m.ChannelID, msg)
	}
}

func onMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	fmt.Printf("%20s %20s > %s\n", time.Now().Format(time.Stamp), r.UserID, r.Emoji.Name)
	reac[r.MessageID] = Reaction{r.UserID, r.MessageID, r.Emoji.Name, time.Now()}
	if len(reacCheckList) != 0 {
		for k, _ := range reacCheckList {
			 reactionCheck(k, s)
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
func reactionCheck(key string, s *discordgo.Session) {
	u := reacCheckList[key].UserID
	m := reacCheckList[key].MessageID
	c := reacCheckList[key].ChannelID
	fmt.Println(u)
	for k, _ := range reac {
		if reac[k].MessageID != m {
			continue
		}
		fmt.Println(reac[k].UserID)
		fmt.Println(u)
		if reac[k].UserID == u {
			if reacCheckList[key].Turget != "" {
				u = reacCheckList[key].Turget
				fmt.Println("turget: ", u)
			}
			
			t := time.Now().Format("20060102")
			if reacCheckList[key].Day != "" {
				t = reacCheckList[key].Day
			}
			cd := module.CheckDate(db, u, t)
			if cd {
				fmt.Println(cd)
				sendMessage(s, c, "ToDo")
				return
			}

			switch {
				case "🤒" == reac[k].Emoji:
					err := module.AddToDB(db, u, t, "Sick")
					if err == nil {
						msg := fmt.Sprintf("Record \n User   : <@%s> \n Date   : %s \n Reason : Sick", u, t)
						sendMessage(s, c, msg)
					} else {
						sendMessage(s, c, "Error")
					}
					delete(reacCheckList, key)
				case "😴" == reac[k].Emoji:
					err := module.AddToDB(db, u, t, "Oversleeping")
					if err == nil {
						msg := fmt.Sprintf("Record \n User   : <@%s> \n Date   : %s \n Reason : Over Sleeping", u, t)
						sendMessage(s, c, msg)
					} else {
						sendMessage(s, c, "Error")
					}
					delete(reacCheckList, key)
				case "🚃" == reac[k].Emoji:
					err := module.AddToDB(db, u, t, "Train delay")
					if err == nil {
						msg := fmt.Sprintf("Record \n User   : <@%s> \n Date   : %s \n Reason : Train delay", u, t)
						sendMessage(s, c, msg)
					} else {
						sendMessage(s, c, "Error")
					}
					delete(reacCheckList, key)
				case "💼" == reac[k].Emoji:
					err := module.AddToDB(db, u, t, "Other")
					if err == nil {
						msg := fmt.Sprintf("Record \n User   : <@%s> \n Date   : %s \n Reason : Other", u, t)
						sendMessage(s, c, msg)
					} else {
						sendMessage(s, c, "Error")
					}
					delete(reacCheckList, key)
				case "🏫" == reac[k].Emoji:
					err := module.AddToDB(db, u, t, "Official absence")
					if err == nil {
						msg := fmt.Sprintf("Record \n User   : <@%s> \n Date   : %s \n Reason : Official absence", u, t)
						sendMessage(s, c, msg)
					} else {
						sendMessage(s, c, "Error")
					}
					delete(reacCheckList, key)
			}
		}
	}
}

//Reaction timeout
func reactionTimeout() {
	for k, _ := range reac {
		t := time.Since(reac[k].Time)
		c := time.Duration(300000000000)
		if c <= t {
			delete(reac,k)
		}
	}

	for k, _ := range reacCheckList {
		t := time.Since(reacCheckList[k].Time)
		c := time.Duration(300000000000)
		if c <= t {
			delete(reacCheckList, k)
		}
	}
}

func isExist(year, month, day int) (error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if date.Year() == year && date.Month() == time.Month(month) && date.Day() == day {
		return nil
	} else {
		return fmt.Errorf("%d-%d-%d is not exist", year, month, day)
	}
}
