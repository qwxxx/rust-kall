package sharkscope

import (
	"SharkScopeParser/config"
	"SharkScopeParser/global"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type TournamentInfo struct {
	Players []global.Player
	Stake   float64
}

var appname string = ""
var password string = ""
var appkey string = ""

func Inizializate() {
	appname = config.Cfg.AppName
	password = config.Cfg.Password
	appkey = config.Cfg.AppKey
}

func request(url string) ([]byte, error) {

	h := md5.New()
	h.Write([]byte(password))

	md5passstr := hex.EncodeToString(h.Sum(nil))

	h = md5.New()
	h.Write([]byte(md5passstr))
	h.Write([]byte(appkey))
	md5key := hex.EncodeToString(h.Sum(nil))
	fmt.Printf("")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Username", "andreiiann@mail.ru")
	req.Header.Set("User-Agent", "Mozilla")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Password", string(md5key[:]))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	b, err := io.ReadAll(res.Body)
	return b, err
}

func GetActiveTournemants() []string {
	ans := []string{}
	js, err := request(`https://www.sharkscope.com/api/` + appname + `/networks/WPN/activeTournaments?Filter=Entrants:2~*;Stake:USD80~250;Type:ST,6MX;Class:SNG`)
	if err != nil {
		return ans
	}
	type RegisteringTournamentS struct {
		ID string `json:"@id"`
	}
	type ActiveTournamentsResponseOne struct {
		Response struct {
			RegisteringTournamentsResponse struct {
				RegisteringTournaments struct {
					RegisteringTournament RegisteringTournamentS `json:"RegisteringTournament"`
				} `json:"RegisteringTournaments"`
			} `json:"RegisteringTournamentsResponse"`
		} `json:"Response"`
	}
	type ActiveTournamentsResponseMany struct {
		Response struct {
			RegisteringTournamentsResponse struct {
				RegisteringTournaments struct {
					RegisteringTournament []RegisteringTournamentS `json:"RegisteringTournament"`
				} `json:"RegisteringTournaments"`
			} `json:"RegisteringTournamentsResponse"`
		} `json:"Response"`
	}
	responseOne := ActiveTournamentsResponseOne{}
	err = json.Unmarshal(js, &responseOne)
	if err != nil {
		responseMany := ActiveTournamentsResponseMany{}
		err = json.Unmarshal(js, &responseMany)
		if err != nil {
			return ans
		}
		for _, t := range responseMany.Response.RegisteringTournamentsResponse.RegisteringTournaments.RegisteringTournament {
			ans = append(ans, t.ID)
		}
	} else {
		ans = append(ans, responseOne.Response.RegisteringTournamentsResponse.RegisteringTournaments.RegisteringTournament.ID)
	}

	return ans
}

func PlayerGamesList(playerName string, network string, startDate, endDate int64) ([]string, error) {
	tournaments := []string{}
	js, err := request(`https://www.sharkscope.com/api/` + appname + `/networks/` + network + `/players/` + playerName + `/completedTournaments?player&filter=Entrants:6~6;Date:` + fmt.Sprint(startDate) + `~` + fmt.Sprint(endDate) + `&order=player,1~9999`)
	if err != nil {
		return tournaments, err
	}
	fmt.Println(string(js))

	player := global.PlayerGamesRes{}

	err = json.Unmarshal(js, &player)
	if err != nil {
		return nil, err
	}
	if player.Response.Success == "false" {
		return nil, fmt.Errorf("request failed: %v", player.Response.ErrorResponse.Error.Message)
	}

	if player.Response.PlayerResponse.PlayerView.Player.CompletedTournament == nil {
		return tournaments, err
	}
	tournamentsMap := player.Response.PlayerResponse.PlayerView.Player.CompletedTournament.(map[string]interface{})
	tour, ok := tournamentsMap["Tournament"].(map[string]interface{})
	if !ok {
		tournamentsMass := tournamentsMap["Tournament"].([]interface{})
		for _, tour := range tournamentsMass {
			val, ok := tour.(map[string]interface{})["@totalEntrants"]
			fmt.Println(val)
			if !ok {
				continue
			}
			valstr, ok := val.(string)
			if !ok || valstr != "6" {
				continue
			}
			tournaments = append(tournaments, fmt.Sprint(tour.(map[string]interface{})["@id"]))
		}
	} else {
		id, ok := tour["@id"]
		if ok {
			tournaments = append(tournaments, fmt.Sprint(id))
		}
	}

	return tournaments, nil
}
func TournamentList(tournamentId string, network string) (TournamentInfo, error) {
	ans := TournamentInfo{Players: []global.Player{}}

	//key := strings.ToLower(md5passstr) + appkey

	b, err := request(`https://www.sharkscope.com/api/` + appname + `/networks/` + network + `/tournaments/` + tournamentId)
	if err != nil {
		return ans, err
	}
	fmt.Println(string(b))
	o := global.TournamentRes{}
	err = json.Unmarshal(b, &o)
	if err != nil {
		fmt.Println(err.Error())
		return ans, err
	}

	if o.Response.Success == "false" {
		return ans, fmt.Errorf("Tournament not found")
	}

	activeEntries := o.Response.TournamentResponse.ATurn.TournamentEntry
	for _, e := range activeEntries {
		blocked := false
		icons, ok := e.Player.Icon.([]interface{})
		if ok {
			for _, in := range icons {
				m, ok := in.(map[string]any)
				if !ok {
					continue
				}
				value, ok := m["@type"]
				if ok && value == "blocked" {
					blocked = true
					break
				}
			}
		} else {
			icon, ok := e.Player.Icon.(map[string]interface{})
			if ok {
				t, ok := icon["@type"]
				if !ok {
					continue
				}
				tstr, ok := t.(string)
				blocked = ok && tstr == "blocked"
			}
		}
		ans.Players = append(ans.Players, global.Player{e.Player.PlayerName, blocked})
	}
	ans.Stake, _ = strconv.ParseFloat(o.Response.TournamentResponse.ATurn.Stake, 64)

	if len(ans.Players) == 0 {
		completedEntries := o.Response.TournamentResponse.CTurn.TournamentEntry

		for _, e := range completedEntries {
			ans.Players = append(ans.Players, global.Player{e.PlayerName, false})
		}
		ans.Stake, _ = strconv.ParseFloat(o.Response.TournamentResponse.CTurn.Stake, 64)
	}

	return ans, err

}

func GetInfo(period global.ReportPeriod) (*global.MessageStat, error) {
	fromDate := int64(0)
	toDate := time.Now().Unix()
	switch period {
	case global.Day:
		fromDate = time.Now().Add(-time.Hour * 24).Unix()
		break
	case global.Week:
		fromDate = time.Now().Add(-time.Hour * 24 * 7).Unix()
		break
	case global.Month:
		fromDate = time.Now().Add(-time.Hour * 24 * 30).Unix() //TODO february and 30/31 days
		break
	}
	fmt.Println(fromDate)
	fmt.Println(toDate)
	js, err := request(`https://www.sharkscope.com/api/andreiiann/networks/Player%20Group/players/Edge%203%25?filter=Entrants:5~*;Date:` + fmt.Sprint(fromDate) + `~` + fmt.Sprint(toDate) + `;Class:SNG&Currency=USD`)
	if err != nil {
		return nil, err
	}
	s := global.Stat{}
	err = json.Unmarshal(js, &s)
	if err != nil {
		return nil, err
	}

	res := global.MessageStat{}
	for _, m := range s.Response.PlayerResponse.PlayerView.PlayerGroup.Statistics.Statistic {
		switch m.Id {
		case "Count":
			res.NumOfTournament, _ = strconv.Atoi(m.Num)
		case "Profit":
			res.Profit, _ = strconv.ParseFloat(m.Num, 64)
		case "AvStake":
			res.ABI, _ = strconv.ParseFloat(m.Num, 64)
		case "TotalROI":
			res.TotalROI, _ = strconv.ParseFloat(m.Num, 64)
		case "AvROI":
			res.AvgROI, _ = strconv.ParseFloat(m.Num, 64)
		case "Rake":
			res.Profit, _ = strconv.ParseFloat(m.Num, 64)

		}
	}

	return &res, nil
}
