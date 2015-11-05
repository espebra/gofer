package main

import (
	"crypto/tls"
	"log"
	"strconv"
	//"net/http"
	//"sync"
	//"time"
	//"fmt"
	"github.com/thoj/go-ircevent"
	"github.com/espebra/gofer/app/config"
)

func main() {
	var err error
	config := config.Configuration {}
	err = config.Read("config.json")
	if err != nil {
		log.Fatal(err)
	}

	i := irc.IRC(config.Nickname, config.Username)
	i.VerboseCallbackHandler = config.Verbose
	i.Debug = config.Debug
	i.UseTLS = true
	i.TLSConfig = &tls.Config{InsecureSkipVerify: config.TLSSkipVerify}

	err = i.Connect(config.Server + ":" + strconv.Itoa(config.Port))
	if err != nil {
		log.Fatal(err)
	}

	for c := range config.Channels {
		var channel = config.Channels[c].Name
		var key = config.Channels[c].Key
		i.AddCallback("001", func(e *irc.Event) { i.Join(channel + " " + key) })
	}

	//i.AddCallback("002", func(e *irc.Event) { i.Privmsg("#mehbar", "zoo") })

	i.Loop()
}
