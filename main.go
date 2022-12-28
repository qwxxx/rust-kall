package main

import (
	"SharkScopeParser/config"
	"SharkScopeParser/discord"
	"SharkScopeParser/rest"
	"SharkScopeParser/sharkscope"
	"SharkScopeParser/store"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	var err error
	config.Cfg, err = config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	discord.SetToken()
	sharkscope.Inizializate()

	d := store.NewStore()
	ds, err := discord.Create()
	if err != nil {
		log.Fatalln(err)
	}
	h := rest.API{
		DB: d,
		DS: &ds,
	}

	d.GetScore("", false, 0)
	d.UpdateScores()

	e := gin.Default()
	e.Use(h.CORSMiddleware())
	// Middleware
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	// Routes

	e.POST("/login", h.Authorization)

	api := e.Group("/api")

	api.POST("/config", h.LoadConfig)
	api.GET("/config", h.GetConfig)
	api.GET("/unknownNames", h.GetUnknownNames)
	api.POST("/unknownNames/clear", h.ClearUnknownNames)
	api.GET("/tournament", h.CalculateTournament)
	api.GET("/player", h.CalculatePlayer)
	api.GET("/state", h.State)
	api.GET("/restart", h.Restart)

	e.StaticFile("/", "./frontend/dist/index.html")
	e.Static("/static", "./frontend/dist/")
	e.Static("/css", "./frontend/dist/css")
	e.Static("/js", "./frontend/dist/js")
	e.Static("/assets", "./frontend/dist/assets")
	e.Static("/files", "./files")

	go ds.SendImportant()
	go func() {
		lastReportDateMsk := time.Now().UTC().Add(time.Hour * 3)
		for range time.Tick(time.Minute * 1) {

			currentDateMsk := time.Now().UTC().Add(time.Hour * 3)

			if currentDateMsk.Day() != lastReportDateMsk.Day() && currentDateMsk.Hour() >= 13 {
				s, _ := sharkscope.GetInfo()
				_, _, day := time.Now().Date()
				if day == 1 {
					lastReportDateMsk = currentDateMsk
					ds.SendStat(s, "Месячный отчет")
				} else if int(currentDateMsk.Weekday()) == 6 {
					lastReportDateMsk = currentDateMsk
					ds.SendStat(s, "Недельный отчет")
				} else {
					lastReportDateMsk = currentDateMsk
					ds.SendStat(s, "Дневной отчет")
				}

			}
		}

	}()
	go h.AutoFindActiveTournaments()
	e.Run(":8081")
}
