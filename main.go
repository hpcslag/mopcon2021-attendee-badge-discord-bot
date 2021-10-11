package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// before using the bot, please use this url paste to browser and accept it to channel or guild
// https://discord.com/api/oauth2/authorize?client_id=[put client id inside]&scope=bot&permissions=8

// Variables used for command line parameters
var (
	BotSecret    string
	GuildID      string
	XLSXFile     string
	BackupPeroid string
)

func init() {

	// -t token -g gulidid -f test.xlsx
	flag.StringVar(&BotSecret, "t", "", "Bot Token")
	flag.StringVar(&GuildID, "g", "", "855718269092233247")
	flag.StringVar(&XLSXFile, "f", "", "test.xlsx")
	flag.StringVar(&BackupPeroid, "period", "", "3600")
	flag.Parse()
}

func main() {

	ReadXLSXToMap(XLSXFile, &KeyPairMap)

	c := cron.New()
	c.AddFunc("@every "+BackupPeroid+"s", func() {
		SaveUsedCSV()
		log.Println("File Backup Done.")
	})
	go c.Start()
	// close cron
	defer c.Stop()

	RegisterBotFuncAndRun(DiscordAuth{
		BotSecret: BotSecret,
	}, messageCreate)
}

// check the used KKTIX Token.
type Record struct {
	User string
	Time string
}

var usedToken map[string]*Record = make(map[string]*Record)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// check message and messages guild role ids
	/*
		fmt.Println("Received msg: ", m.Content)
		fmt.Println("GuildID:", m.GuildID)
		st, _ := s.GuildRoles(m.GuildID)
		for _, v := range st {
			fmt.Println("RoleName and ID: ", v.Name, v.ID)
		}
	*/

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// https://github-wiki-see.page/m/bwmarrin/discordgo/wiki/FAQ
	isDM, err := ComesFromDM(s, m)
	if err != nil {
		panic(err)
	}

	// comefrom private messahe
	if isDM {

		// split to command
		commands := strings.Split(m.Content, " ")

		if len(commands) == 0 {
			s.ChannelMessageSend(m.ChannelID, "找不到這項指令，若您希望註冊您在 MOPCON 2021 的身分，請您輸入: 「kktix [您的票號]」。")
			return
		}

		if commands[0] == "kktix" {
			// check commands
			if len(commands) != 2 {
				s.ChannelMessageSend(m.ChannelID, "您還沒輸入正確的票號喔，您需要輸入: 「kktix [您的票號]」，範例: 「kktix 123456789」。")
				return
			}

			// block duplicate register
			if _, ok := usedToken[commands[1]]; ok {
				s.ChannelMessageSend(m.ChannelID, "這個票號已經被註冊過了，請尋求服務台頻道的協助。")
				return
			}

			// given badge and set used token
			if badge, ok := KeyPairMap[commands[1]]; ok {
				usedToken[commands[1]] = &Record{
					User: m.Author.ID,
					Time: time.Now().Format("2006-01-02 15:04:05"),
				}

				s.GuildMemberRoleAdd(GuildID, m.Author.ID, badge)

				s.ChannelMessageSend(m.ChannelID, "您的身分已經完成設定，歡迎您回到 MOPCON Discord 會場!")
			} else {
				s.ChannelMessageSend(m.ChannelID, "這個票號不存在，請尋求服務台頻道的協助。")
				return
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "這個指令不存在，您需要輸入: 「kktix [您的票號]」，範例: 「kktix 123456789」 以進行註冊。")
			return
		}
	} else {
		// receive global message
		s.ChannelMessageSend(m.ChannelID, "請您使用私訊的方式進行 KKTIX 驗票註冊，謝謝您!")
	}
}

func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func SaveUsedCSV() {

	csvData := [][]string{
		{"code", "user", "time"},
	}
	for key, val := range usedToken {
		csvData = append(csvData, []string{key, val.User, val.Time})
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	w.WriteAll(csvData)

	t := time.Now()
	os.WriteFile(t.Format("20060102150405")+".csv", b.Bytes(), 0644)

}
