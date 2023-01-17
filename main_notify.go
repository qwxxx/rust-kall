package main

import (
	"SharkScopeParser/config"
	"SharkScopeParser/discord"
	"SharkScopeParser/global"
	"SharkScopeParser/rest"
	"SharkScopeParser/sharkscope"
	"SharkScopeParser/store"
	"log"
	"strconv"
	"time"
)

func main() {
	var err error
	config.Cfg, err = config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	sharkscope.Inizializate()
	d, err := store.NewStore(config.Cfg.SqlConn, "./migrations")
	if err != nil {
		log.Fatal(err)
	}
	ds, err := discord.Create()
	if err != nil {
		log.Fatalln(err)
	}

	discord.SetToken()
	go ds.SendImportant()
	go func() {
		lastReportDateMsk := time.Now().UTC().Add(time.Hour * 3)
		for range time.Tick(time.Minute * 1) {

			currentDateMsk := time.Now().UTC().Add(time.Hour * 3)

			if currentDateMsk.Day() != lastReportDateMsk.Day() && currentDateMsk.Hour() >= 13 {

				_, _, day := time.Now().Date()
				if day == 1 {
					lastReportDateMsk = currentDateMsk
					s, _ := sharkscope.GetInfo(global.Month)
					ds.SendStat(s, "Месячный отчет")
				}
				if int(currentDateMsk.Weekday()) == 6 {
					lastReportDateMsk = currentDateMsk
					s, _ := sharkscope.GetInfo(global.Week)
					ds.SendStat(s, "Недельный отчет")
				}

				lastReportDateMsk = currentDateMsk
				s, _ := sharkscope.GetInfo(global.Day)
				ds.SendStat(s, "Дневной отчет")

			}
		}

	}()
	go AutoFindActiveTournaments(ds, d)
	ds.SendTest()
}

func AutoFindActiveTournaments(ds discord.Discord, DB store.Storage) {
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
		oldTornamentsIds := newTournamentIds

		go func(DB store.Storage, ds discord.Discord) {
			ticker := time.NewTicker(10 * time.Minute)
			select {
			case <-ticker.C:
				find(oldTornamentsIds, newTournamentIds, DB, ds)
			}
		}(DB, ds)

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
				calculated, err := rest.CalculateTournamentRaw("WPN", new_id, DB)
				if err != nil || (calculated.TotalScore < 8 && calculated.PlayersCount != 6) {
					continue
				}
				if calculated.PlayersCount != 6 && calculated.TotalScore >= 9 {
					mid := ds.SendTournamentInfo(calculated)
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
				calculated, err := rest.CalculateTournamentRaw("WPN", old_id, DB)
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
				ds.SendReplyWithUpdated(tournament.MessageId, tournament.Calculated)
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

func find(old, new []string, DB store.Storage, ds discord.Discord) {
	for _, oldId := range old {
		for _, newId := range new {
			if oldId == newId {
				calculated, err := rest.CalculateTournamentRaw("WPN", newId, DB)
				if err != nil || (calculated.TotalScore < 8 && calculated.PlayersCount != 6) {
					continue
				}
				if calculated.PlayersCount != 6 && calculated.TotalScore >= 9 {
					ds.SendTournamentInfo(calculated)
				}
			}
		}
	}
}
