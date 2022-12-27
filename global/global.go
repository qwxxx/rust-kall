package global

var UserPassword = `userPass`
var AdminPassword = `12345g`

var UnknownName = "<UNKNOWN>"

type CalculateTournamentResponsePlayer struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Unknown bool   `json:"unknown"`
	Blocked bool   `json:"blocked"`
}
type CalculateTournamentResponse struct {
	Players      []CalculateTournamentResponsePlayer `json:"players"`
	PlayersCount int                                 `json:"playersCount"`
	TotalScore   int                                 `json:"total_score"`
	Stake        float64                             `json:"stake"`
	Id           int64                               `json:"id"`
}
type CompletedTournamentEntry struct {
	PlayerName string `json:"@playerName"`
}
type Player struct {
	Name    string `json:"@name"`
	Blocked bool
}
type Icon struct {
	Type string `json:"@type"`
}
type ActiveTournamentEntry struct {
	Player struct {
		PlayerName string      `json:"@name"`
		Icon       interface{} `json:"Icon"`
	} `json:"Player"`
}
type CompletedTournament struct {
	Stake           string                     `json:"@stake"`
	TournamentEntry []CompletedTournamentEntry `json:"TournamentEntry"`
}
type ActiveTournament struct {
	Stake           string                  `json:"@stake"`
	TournamentEntry []ActiveTournamentEntry `json:"TournamentEntry"`
}

type TournamentRes struct {
	Response struct {
		Success            string `json:"@success"`
		TournamentResponse struct {
			CTurn CompletedTournament `json:"CompletedTournament"`
			ATurn ActiveTournament    `json:"ActiveTournament"`
		} `json:"TournamentResponse"`
	} `json:"Response"`
}

type PlayerGamesRes struct {
	Response struct {
		Success       string `json:"@success"`
		ErrorResponse struct {
			Error struct {
				Message string `json:"$"`
			} `json:"Error"`
		} `json:"ErrorResponse"`
		PlayerResponse struct {
			PlayerView struct {
				Player struct {
					CompletedTournament interface{} `json:"CompletedTournaments"`
				} `json:"Player"`
			} `json:"PlayerView"`
		} `json:"PlayerResponse"`
	} `json:"Response"`
}

type Stat struct {
	Response struct {
		PlayerResponse struct {
			PlayerView struct {
				PlayerGroup struct {
					Statistics struct {
						Statistic []struct {
							Id  string `json:"@id"`
							Num string `json:"$"`
						} `json:"Statistic"`
					} `json:"Statistics"`
				} `json:"PlayerGroup"`
			} `json:"PlayerView"`
		} `json:"PlayerResponse"`
	} `json:"Response"`
}

type MessageStat struct {
	NumOfTournament int
	Profit          float64
	ABI             float64
	TotalROI        float64
	AvgROI          float64
	TotalReik       float64
}
