package main

import (
	"SharkScopeParser/config"
	"SharkScopeParser/rest"
	"SharkScopeParser/sharkscope"
	"SharkScopeParser/store"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var err error
	config.Cfg, err = config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	sharkscope.Inizializate()

	d, err := store.NewStore(config.Cfg.SqlConn, "./migrations")
	h := rest.API{
		DB: d,
	}

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

	e.Run(":8081")
}
