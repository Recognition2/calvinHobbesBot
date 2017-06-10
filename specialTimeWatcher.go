package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os/exec"
	"time"
)

func timeWatcher() {
	defer g.wg.Done()
	defer logWarn.Println("Shutting time watcher down")

outer:
	for {
		n := time.Minute // Check whether people are to be notified every n minutes
		t := time.After(n)
		select {
		case <-g.shutdown:
			break outer
		case <-t:
			go checkNotifications()
		}
	}
}

func checkNotifications() {
	// Get current time, accurate to the minute
	cTime := time.Now().Hour()*100 + time.Now().Minute()

	subs, err := g.db.Query(`
		SELECT tID, currentStrip
		FROM subscriptons
		WHERE active
		AND warnTime=?
	`, cTime)

	if err != nil {
		logErr.Println(err)
	}

	// Declare variables to put the values into
	var tID int64
	var currentStrip int

	for subs.Next() {
		err = subs.Scan(&tID, &currentStrip)
		if err != nil {
			logErr.Println(err)
		}

		sendStrip(tID, currentStrip)
	}
}

func sendStrip(id int64, currentStrip int) {
	stripToSend := findStrip(currentStrip + 1)
}

func findStrip(id int) int {

}
