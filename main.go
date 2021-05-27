package main

import (
	"flag"
	"github.com/skarakasoglu/aybush-bot-api/api"
	"github.com/skarakasoglu/aybush-bot-api/data"
	"github.com/skarakasoglu/aybush-bot-api/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	dbHost string
	dbPort int
	dbUsername string
	dbPassword string
	dbName string
	certFile string
	keyFile string
)

func init() {
	flag.StringVar(&dbHost, "db-ip-address", "", "database ip address")
	flag.IntVar(&dbPort, "db-port", 0, "database port")
	flag.StringVar(&dbUsername, "db-username", "", "database login username")
	flag.StringVar(&dbPassword, "db-password", "", "database login password")
	flag.StringVar(&dbName, "db-name", "", "database name")
	flag.StringVar(&certFile, "cert-file", "", "ssl certificate file")
	flag.StringVar(&keyFile, "key-file", "", "ssl private key file")
	flag.Parse()
}

func main() {
	log.Println(dbHost)

	db, err := data.NewDB(data.DatabaseCredentials{
		Host:         dbHost,
		Port:         dbPort,
		Username:     dbUsername,
		Password:     dbPassword,
		DatabaseName: dbName,
	})
	if err != nil {
		log.Printf("Error on creating database connection: %v", err)
		return
	}

	discordRepository := service.NewDiscordService(db)

	srv := api.NewServer(discordRepository, certFile, keyFile)
	srv.Start()

	log.Println("AYBUSH BOT WEB API is now running. Press CTRL + C to interrupt.")
	signalHandler := make(chan os.Signal)
	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	receivedSignal := <-signalHandler

	log.Printf("[AybushBotWebApi] %v signal received. Gracefully shutting down the application.", receivedSignal)
	log.Println("[AybushBotWebApi] Application exited.")
}
