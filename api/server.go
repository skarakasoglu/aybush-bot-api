package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skarakasoglu/aybush-bot-api/repository"
	"log"
)

type server struct{
	discordRepository repository.DiscordRepository

	certFile string
	keyFile string
}

func NewServer(discordRepository repository.DiscordRepository, certFile string, keyFile string) *server{
	return &server{
		discordRepository: discordRepository,
		certFile: certFile,
		keyFile: keyFile,
	}
}


func (srv *server) Start() error {
	router := gin.Default()

	apiv1 := NewApiV1(srv.discordRepository)

	v1 := router.Group("/v1")
	{
		v1.GET("/leaderboard", apiv1.getLeaderboardByEpisode)
		v1.GET("/episode", apiv1.getActiveEpisodes)
	}

	go func() {
		err := router.RunTLS(":443", srv.certFile, srv.keyFile)
		if err != nil {
			log.Printf("Error on running the router: %v", err)
		}
	}()

	return nil
}