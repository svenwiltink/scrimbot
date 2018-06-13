package main

import (
	"flag"
	"fmt"
	"github.com/svenwiltink/scrimbot/bot"
	"github.com/svenwiltink/scrimbot/config"
	"github.com/svenwiltink/scrimbot/scrimpost"
	"github.com/svenwiltink/scrimbot/scrimpost/database/bolt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configLocation := flag.String("c", "config.json", "config file location")
	configStruct, err := config.LoadConfig(*configLocation)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := bolt.Load(configStruct.BoltDatabase)
	scrimpost.RegisterDatabase(db)
	instance := bot.CreateNewBot(configStruct)
	err = instance.Start()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	instance.Close()

}
