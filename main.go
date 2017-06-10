package main

import (
	"database/sql"
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	logErr  = log.New(os.Stderr, "[ERRO] ", log.Ldate+log.Ltime+log.Ltime+log.Lshortfile)
	logWarn = log.New(os.Stdout, "[WARN] ", log.Ldate+log.Ltime)
	logInfo = log.New(os.Stdout, "[INFO] ", log.Ldate+log.Ltime)
	g       = global{shutdown: make(chan bool)}
)

func main() {
	/////////////
	// STARTUP
	//////////////

	// Parse settings file
	_, err := toml.DecodeFile("settings.toml", &g.c)
	if err != nil {
		logErr.Println(err)
		return
	}

	// Create new bot
	g.bot, err = tgbotapi.NewBotAPI(g.c.Apikey)
	if err != nil {
		logErr.Println(err)
	}

	logInfo.Printf("Running as @%s", g.bot.Self.UserName)

	// Create waitgroup, for synchronized shutdown
	var wg sync.WaitGroup
	g.wg = &wg

	// DB Initialization
	user := g.c.Mysql_user
	passwd := g.c.Mysql_passwd
	dbname := g.c.Mysql_dbname

	db, err := sql.Open("mysql", user+":"+passwd+"@/"+dbname) // DOES NOT open a connection
	if err != nil {
		logErr.Println(err)
	}

	err = db.Ping() // Validating DSN data
	if err != nil {
		logErr.Printf("Failed opening db connection: %v\n", err)
		close(g.shutdown)
		return
	}
	defer db.Close()

	g.db = db
	logInfo.Println("Database connection opened!")

	wg.Add(2)
	go messageMonitor() // Receiving messages
	go timeWatcher()    // Sending pictures

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT)

	logInfo.Println("All routines have been started, awaiting kill signal")

	///////////////
	// SHUTDOWN
	///////////////

	// Program will hang here
	select {
	case <-sigs:
		close(g.shutdown)
	case <-g.shutdown:
	}
	println()
	logInfo.Println("Shutdown signal received. Waiting for goroutines")

	// Shutdown after all goroutines have exited
	g.wg.Wait()
	logWarn.Println("Shutting down")
}
