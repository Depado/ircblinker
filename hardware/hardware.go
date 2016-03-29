package hardware

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/thoj/go-ircevent"

	"github.com/Depado/ircblinker/conf"
)

// Unread is the variable used to communicate between the ircbot and the hardware
var Unread = false

// AFK is a variable that tells that the user is AFK or not
var AFK = false

// Start is the function that starts the gobot instance
func Start(ib *irc.Connection) {
	var err error
	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
	hlLed := gpio.NewLedDriver(r, "hlled", "12")
	afkLed := gpio.NewLedDriver(r, "afkled", "11")
	serviceLed := gpio.NewLedDriver(r, "serviceled", "13")
	button := gpio.NewButtonDriver(r, "afkbtn", "5")

	work := func() {
		on := false
		gobot.On(button.Event("push"), func(data interface{}) {
			AFK = !AFK
			if AFK {
				ib.Privmsg(conf.C.User, "AFK Button Pushed : AFK Mode Activated")
				if err = afkLed.On(); err != nil {
					log.Fatal(err)
				}
			} else {
				ib.Privmsg(conf.C.User, "AFK Button Pushed : AFK Mode Deactivated")
				if err = afkLed.Off(); err != nil {
					log.Fatal(err)
				}
			}
		})
		gobot.Every(500*time.Millisecond, func() {
			if Unread {
				on = true
				if err = hlLed.Toggle(); err != nil {
					log.Fatal(err)
				}
			} else {
				if on {
					if err = hlLed.Off(); err != nil {
						log.Fatal(err)
					}
				}
			}
		})
		up, err := AllServicesUp()
		if err != nil {
			log.Println("Errored when querying status :", err)
		}
		if !up {
			if err = serviceLed.On(); err != nil {
				log.Fatal(err)
			}
		} else {
			if err = serviceLed.Off(); err != nil {
				log.Fatal(err)
			}
		}
		gobot.Every(30*time.Minute, func() {
			up, err := AllServicesUp()
			if err != nil {
				log.Println("Errored when querying status :", err)
			}
			if !up {
				if err = serviceLed.On(); err != nil {
					log.Fatal(err)
				}
			} else {
				if err = serviceLed.Off(); err != nil {
					log.Fatal(err)
				}
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{hlLed, afkLed, serviceLed, button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
