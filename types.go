package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
)

type global struct {
	wg       *sync.WaitGroup  // For checking that everything has indeed shut down
	shutdown chan bool        // To make sure everything can shut down
	bot      *tgbotapi.BotAPI // The actual bot
	c        config           // Configuration file
	db       *sql.DB          // DB connection
}

type config struct {
	Apikey       string  // Telegram API key
	Admins       []int64 // Bot admins
	Mysql_user   string  // SQL username
	Mysql_passwd string  // Passwd
	Mysql_dbname string
}
