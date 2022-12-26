package discord

import (
	"SharkScopeParser/config"
	"SharkScopeParser/global"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ImportantChannelID string = ``
var StatChannelID string = ``
var channelID string = `1051584641158631424`
var token string = config.Cfg.Token

type Discord struct {
	session *discordgo.Session
}

func Create() (Discord, error) {
	session, err := discordgo.New(token)
	return Discord{session: session}, err
}
func (c *Discord) SendTest() {
	m, _ := c.session.ChannelMessages(channelID, 3, "", "", "")
	cc, err := c.session.MessageThreadStart(channelID, m[2].ID, "ТЕСТОВЫЙ ПОТОК", 0)

	fmt.Println(err)
	c.session.ChannelMessageSend(cc.ID, "Тестовый ответ")

}
func (c *Discord) SendReplyWithUpdated(messageID string, response global.CalculateTournamentResponse) {
	text := fmt.Sprintf("Итоговая оценка турнира %v: %d\n\n", response.Id, response.TotalScore)
	for i, v := range response.Players {
		if i == 0 || i == 1 {
			text += fmt.Sprintf("Место %v (ITM) - Player: %s\n", i+1, v.Name)
		} else {
			text += fmt.Sprintf("Место %v - Player: %s\n", i+1, v.Name)
		}
	}
	m, err := c.session.MessageThreadStart(channelID, messageID, "Турнир закончен", 0)
	if err != nil {
		return
	}
	c.session.ChannelMessageSend(m.ID, text)
	fmt.Println(err)
}
func (c *Discord) SendImportant() {
	for {
		lastReportDateMsk := time.Now().UTC().Add(time.Hour * 3)
		currentDateMsk := time.Now().UTC().Add(time.Hour * 3)

		if currentDateMsk.Day() != lastReportDateMsk.Day() && currentDateMsk.Hour() >= 13 && int(time.Now().Weekday()) == 4 {
			lastReportDateMsk = currentDateMsk
			c.session.ChannelMessageSend(ImportantChannelID, ":robot: @everyone\nНапоминаю, вы автоматически зарегистрировались через билет 95$ в турнире (кто по еженедельном ЛБ на эту ступеньку попал)\nЕсли не планируете играть именно этот турнир - отрегистрируйтесь")
		}
	}

}
func (c *Discord) SendTournamentInfo(response global.CalculateTournamentResponse) string {
	text := fmt.Sprintf(`Турнир: %d
$%d Hyper Turbo 6-max
Оценка: %d
Занято мест: %d/6`, response.Id, int64(response.Stake), response.TotalScore, response.PlayersCount)
	m, err := c.session.ChannelMessageSend(channelID, text)
	if err != nil {
		return ""
	}
	return m.ID
}

func (c *Discord) SendStat(stat *global.MessageStat, name string) {
	text := fmt.Sprintf("%v\nКоличество турниров: %v\nПрофит: %v\nАБИ: %v\nОбщий рой: %v\nСредний рой: %v\nОбщий рейк: %v", name, stat.NumOfTournament, stat.Profit, stat.ABI, stat.TotalROI, stat.AvgROI, stat.TotalReik)
	c.session.ChannelMessageSend(StatChannelID, text)
}
