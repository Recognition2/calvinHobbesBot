package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func messageMonitor() {
	defer g.wg.Done()
	logInfo.Println("Starting message monitor")
	defer logWarn.Println("Stopping message monitor")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 300
	updates, err := g.bot.GetUpdatesChan(u)
	if err != nil {
		logErr.Printf("Update failed: %v\n", err)
	}

outer:
	for {
		select {
		case <-g.shutdown:
			break outer
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			if update.Message.IsCommand() {
				handleMessage(update.Message)
			}
		}
	}
}

func commandIsForMe(t string) bool {
	command := strings.SplitN(t, " ", 2)[0] // Return first substring before space, this is entire command

	i := strings.Index(command, "@") // Position of @ in command
	if i == -1 {                     // Not in command
		return true // Assume command is for everybody, including this bot
	}

	return strings.ToLower(command[i+1:]) == strings.ToLower(g.bot.Self.UserName)
}

func handleMessage(m *tgbotapi.Message) {
	if !commandIsForMe(m.Text) {
		return
	}

	switch strings.ToLower(m.Command()) {
	case "start", "pause", "stop":
		handleStart(m)
	case "id":
		handleGetID(m)
	case "help":
		handleHelp(m)
	case "hi":
		handleHi(m)
	}
}

func handleStart(m *tgbotapi.Message) {

}

func handleHelp(m *tgbotapi.Message) {
	msg := "This bot warns you at special times. Add a time at which you want to be warned every day using '/add'"
	g.bot.Send(tgbotapi.NewMessage(m.Chat.ID, msg))
}

func handleHi(m *tgbotapi.Message) {
	g.bot.Send(tgbotapi.NewMessage(m.Chat.ID, "Hi!"))

}

func handleGetID(cmd *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(cmd.Chat.ID, fmt.Sprintf("Hi, %s %s, your Telegram user ID is given by %d", cmd.From.FirstName, cmd.From.LastName, cmd.From.ID))
	_, err := g.bot.Send(msg)
	if err != nil {
		logErr.Println(err)
	}
}
