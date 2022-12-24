package discord

import (
	"SharkScopeParser/config"
	"SharkScopeParser/global"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// var channelID string = `1008079845441409045`
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
	text := fmt.Sprintf("Итоговая оценка: %d\n\n", response.TotalScore)
	for _, v := range response.Players {
		text += fmt.Sprintf("%s: %d\n", v.Name, v.Score)
	}
	m, err := c.session.MessageThreadStart(channelID, messageID, "Турнир закончен", 0)
	if err != nil {
		return
	}
	c.session.ChannelMessageSend(m.ID, text)
	fmt.Println(err)
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
