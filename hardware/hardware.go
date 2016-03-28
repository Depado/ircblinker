package hardware

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

// Unread is the variable used to communicate between the ircbot and the hardware
var Unread = false

// Start is the function that starts the gobot instance
func Start() {
	var err error
	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
	led := gpio.NewLedDriver(r, "led", "12")

	work := func() {
		on := false
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
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
