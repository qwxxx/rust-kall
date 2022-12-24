package rest

import (
	"SharkScopeParser/discord"
	"SharkScopeParser/global"
	"SharkScopeParser/sharkscope"
	"SharkScopeParser/store"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type API struct {
	DB                       *store.Store
	DS                       *discord.Discord
	isPlayerCalculateRunning bool
}

func (h *API) ClearUnknownNames(c *gin.Context) {
	h.DB.UnknownNames = map[string]bool{}
	os.Remove("out.xlsx")
	h.DB.GetScore("", false, 0)
	c.Status(200)
}

func (h *API) AutoFindActiveTournaments() {
	type FindMode int
	const (
		slowMode FindMode = iota
		fastMode
	)
	tournamentIds := []string{}
	sentTournamentIds := map[string]string{}
	mode := fastMode
	lastReportDateMsk := time.Now().UTC().Add(time.Hour * 3)
	type DiscordEndTournament struct {
		MessageId  string
		Calculated global.CalculateTournamentResponse
	}
	discordEndTournaments := make([]DiscordEndTournament, 0)
	for {
		newTournamentIds := sharkscope.GetActiveTournemants()
		if len(newTournamentIds) == 0 {
			mode = slowMode
		} else {
			mode = fastMode
		}
		tempTournamentIds := make([]string, 0)
		for _, new_id := range newTournamentIds {
			isRlyNew := true
			for _, old_id := range tournamentIds {
				if old_id == new_id {
					isRlyNew = false
					break
				}
			}
			if isRlyNew {
				calculated, err := h.CalculateTournamentRaw("WPN", new_id)
				if err != nil || (calculated.TotalScore < 8 && calculated.PlayersCount != 6) {
					continue
				}
				if calculated.PlayersCount != 6 && calculated.TotalScore >= 9 {
					mid := h.DS.SendTournamentInfo(calculated)
					sentTournamentIds[strconv.FormatInt(calculated.Id, 10)] = mid
				}
			}
			tempTournamentIds = append(tempTournamentIds, new_id)
		}
		for _, old_id := range tournamentIds {
			isStayed := false
			for _, new_id := range newTournamentIds {
				if old_id == new_id {
					isStayed = true
					break
				}
			}
			if dsMessageId, ok := sentTournamentIds[old_id]; ok && !isStayed {
				calculated, err := h.CalculateTournamentRaw("WPN", old_id)
				if err != nil || calculated.PlayersCount < 6 {
					continue
				}
				discordEndTournaments = append(discordEndTournaments, DiscordEndTournament{MessageId: dsMessageId, Calculated: calculated})
				delete(sentTournamentIds, old_id)
			}
		}
		currentDateMsk := time.Now().UTC().Add(time.Hour * 3)
		if currentDateMsk.Day() != lastReportDateMsk.Day() && currentDateMsk.Hour() >= 13 {
			lastReportDateMsk = currentDateMsk
			for _, tournament := range discordEndTournaments {
				h.DS.SendReplyWithUpdated(tournament.MessageId, tournament.Calculated)
			}
			discordEndTournaments = make([]DiscordEndTournament, 0)
		}
		tournamentIds = tempTournamentIds
		switch mode {
		case slowMode:
			time.Sleep(time.Minute * 2)
			break
		case fastMode:
			time.Sleep(time.Second * 20)
			break
		}
	}
}
func (h *API) CalculateTournamentRaw(network string, tournamentId string) (global.CalculateTournamentResponse, error) {
	res := global.CalculateTournamentResponse{}
	tournamentInfo, err := sharkscope.TournamentList(tournamentId, network)
	if err != nil {
		return res, err
	}

	sum := 0

	for _, n := range tournamentInfo.Players {
		pl := global.CalculateTournamentResponsePlayer{}
		pl.Name = n.Name
		if n.Blocked {
			pl.Blocked = true
			if n.Name == global.UnknownName {
				pl.Score = 0
				pl.Name = ""
			} else {
				pl.Score = h.DB.UnknownPlusLocked
			}
		}

		switch pl.Name {
		case global.UnknownName:
			pl.Blocked = true
			pl.Score = 0
			pl.Name = ""
			pl.Unknown = false
		default:
			pl.Score = h.DB.GetScore(n.Name, pl.Blocked, tournamentInfo.Stake)
			pl.Unknown = !h.DB.IsKnownPlayer(n.Name)
		}
		sum += pl.Score
		res.Players = append(res.Players, pl)
	}

	for i := len(tournamentInfo.Players) + 1; i <= 6; i++ {
		pl := global.CalculateTournamentResponsePlayer{}
		pl.Name = fmt.Sprintf("%d место", i)
		pl.Score = h.DB.PlaceScores[i-1]
		sum += pl.Score
		res.Players = append(res.Players, pl)
	}
	res.Id, _ = strconv.ParseInt(tournamentId, 10, 64)
	res.TotalScore = sum
	res.Stake = tournamentInfo.Stake
	res.PlayersCount = len(tournamentInfo.Players)
	return res, nil
}
func (h *API) CalculatePlayer(c *gin.Context) {
	password := c.Query("password")
	if password != global.AdminPassword {
		c.Status(401)
	}
	h.isPlayerCalculateRunning = true
	type CalculatePlayerResponse struct {
		PlayerName              string  `json:"playerName"`
		AverageTournamentsScore float64 `json:"average_tournaments_score"`
		StartDate               string  `json:"start_date"`
		EndDate                 string  `json:"end_date"`
		TournamentsCount        int     `json:"tournaments_count"`
		Error                   int     `json:"error"`
	}

	network := c.Query("network")
	playerName := c.Query("playerName")
	startDateParam := c.Query("startDate")
	startDate, err := strconv.ParseInt(startDateParam, 10, 64)
	if err != nil {
		h.isPlayerCalculateRunning = false
		c.AbortWithError(500, err)
	}
	endDateParam := c.Query("endDate")
	endDate, err := strconv.ParseInt(endDateParam, 10, 64)
	if err != nil {
		h.isPlayerCalculateRunning = false
		c.AbortWithError(500, err)
	}

	res := CalculatePlayerResponse{}
	tournamentIds, err := sharkscope.PlayerGamesList(playerName, network, startDate, endDate)
	if err != nil && err.Error() == `request failed: Player not found or has opted out.` {
		res.Error = 1
	} else {
		if err != nil {
			h.isPlayerCalculateRunning = false
			c.AbortWithError(500, err)
			return
		}
		sumscore := 0
		for _, tournamentId := range tournamentIds {
			s, err := h.CalculateTournamentRaw(network, tournamentId)
			if err != nil {
				continue
			}
			sumscore += s.TotalScore
		}
		if len(tournamentIds) == 0 {
			res.Error = 2
		} else {
			res.AverageTournamentsScore = float64(sumscore) / float64(len(tournamentIds))
		}

	}
	endDate -= 60 * 60 * 24
	res.EndDate = strconv.FormatInt(endDate, 10)
	res.StartDate = strconv.FormatInt(startDate, 10)
	res.PlayerName = playerName
	res.TournamentsCount = len(tournamentIds)

	h.isPlayerCalculateRunning = false
	c.JSON(200, res)
}
func (h *API) CalculateTournament(c *gin.Context) {
	network := c.Query("network")
	tournamentId := c.Query("tournament_id")
	res, err := h.CalculateTournamentRaw(network, tournamentId)
	if err != nil {
		if "Tournament not found" == err.Error() {
			c.AbortWithError(404, err)

		}

		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, res)

}
func (h *API) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func (h *API) Authorization(c *gin.Context) {
	var err error
	var password string
	password = c.Query("password")
	user_id := 0
	switch password {
	case global.AdminPassword:
		user_id = 1
	case global.UserPassword:
		user_id = 2
	default:
		log.Printf("Get user error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{"user_id": user_id})
}
func (h *API) GetUnknownNames(c *gin.Context) {
	password := c.Query("password")
	if password != global.AdminPassword {
		c.Status(401)
	}
	c.File("out.xlsx")
}
func (h *API) GetConfig(c *gin.Context) {
	password := c.Query("password")
	if password != global.AdminPassword {
		c.Status(401)
	}
	c.File("in.xlsx")
}
func (h *API) LoadConfig(c *gin.Context) {
	password := c.Query("password")
	if password != global.AdminPassword {
		c.Status(401)
		return
	}
	type uploadFile struct {
		Way string `json:"way"`
	}
	var err error

	if mf, err := c.MultipartForm(); err == nil {
		defer func() {
			mf.RemoveAll()
		}()
	}

	form, _ := c.MultipartForm()
	files := form.File["files"]

	if len(files) == 0 {
		c.Status(500)
		return
	}
	file := files[0]
	src, err := file.Open()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer src.Close()

	if fl, ok := src.(*os.File); ok && fl != nil {
		defer os.Remove(fl.Name())
	}

	os.Remove("in.xlsx")
	dst, err := os.Create("in.xlsx")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		c.AbortWithError(500, err)
		return
	}
	repeatingNames := h.DB.UpdateScores()
	if len(repeatingNames) > 0 {
		c.AbortWithStatusJSON(400, repeatingNames)
	}
	c.Status(http.StatusCreated)
}
func (h *API) State(c *gin.Context) {
	c.JSON(200, map[string]interface{}{"isPlayerCalculateRunning": h.isPlayerCalculateRunning})
}
