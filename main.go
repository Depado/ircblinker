package main

import (
	"crypto/tls"
	"log"
	"strings"
	"time"

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
		m := e.Message()
		s := strings.Fields(m)

		if to != conf.C.Channel {
			return
		}
		if utils.StringInSlice(e.Nick, conf.C.HL) {
			hardware.Unread = false
			hardware.AFK = false
			return
		}
		if utils.Any(conf.C.HL, s) {
			hardware.Unread = true
			if hardware.AFK {
				ib.Privmsgf(conf.C.User, "%s mentionned you at %s", e.Nick, time.Now())
				ib.Privmsg(conf.C.User, m)
			}
		}
	})

	go hardware.Start(ib)
	ib.Loop()
}
