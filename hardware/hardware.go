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
	led := gpio.NewLedDriver(r, "led", "12")
	button := gpio.NewButtonDriver(r, "afkbtn", "5")

	work := func() {
		on := false
		gobot.On(button.Event("push"), func(data interface{}) {
			AFK = !AFK
			if AFK {
				ib.Privmsg(conf.C.User, "AFK Button Pushed : AFK Mode Activated")
			} else {
				ib.Privmsg(conf.C.User, "AFK Button Pushed : AFK Mode Deactivated")
			}
		})
		gobot.Every(500*time.Millisecond, func() {
			if Unread {
				on = true
				if err = led.Toggle(); err != nil {
					log.Fatal(err)
				}
			} else {
				if on {
					if err = led.Off(); err != nil {
						log.Fatal(err)
					}
				}
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led, button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
