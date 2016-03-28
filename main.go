package main

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/Depado/ircblinker/conf"
	"github.com/Depado/ircblinker/hardware"
	"github.com/Depado/ircblinker/utils"
	"github.com/thoj/go-ircevent"
)

func main() {
	var err error
	var ib *irc.Connection

	if err = conf.Load("conf.yml"); err != nil {
		log.Fatal(err)
	}
	ib = irc.IRC(conf.C.Nick, conf.C.Nick)
	if conf.C.TLS {
		ib.UseTLS = true
		if conf.C.InsecureTLS {
			ib.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
	if err = ib.Connect(conf.C.Server); err != nil {
		log.Fatal(err)
	}

	ib.AddCallback("001", func(e *irc.Event) {
		ib.Join(conf.C.Channel)
	})

	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		to := e.Arguments[0]
		m := strings.Fields(e.Message())

		if to != conf.C.Channel {
			return
		}
		if utils.StringInSlice(e.Nick, conf.C.HL) {
			hardware.Unread = false
			return
		}
		if utils.Any(conf.C.HL, m) {
			hardware.Unread = true
		}
	})

	go hardware.Start()
	ib.Loop()
}
