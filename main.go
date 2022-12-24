package main

import (
	"SharkScopeParser/discord"
	"SharkScopeParser/rest"
	"SharkScopeParser/store"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
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
	api.GET("/state", h.CalculatePlayer)

	e.StaticFile("/", "./frontend/dist/index.html")
	e.Static("/static", "./frontend/dist/")
	e.Static("/css", "./frontend/dist/css")
	e.Static("/js", "./frontend/dist/js")
	e.Static("/assets", "./frontend/dist/assets")
	e.Static("/files", "./files")

	go h.AutoFindActiveTournaments()
	e.Run(":8081")
}
