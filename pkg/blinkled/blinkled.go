package blinkled

import (
	"strconv"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func CreteBlinkingRobot(pin int) *gobot.Robot {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, strconv.Itoa(pin))

	work := func() {
		gobot.Every(2*time.Second, func() {
			led.Toggle()
		})
	}

	return gobot.NewRobot("blinkLED",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)
}
