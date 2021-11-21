package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	initlog("log.txt")
	log.Info("start")
	ledBlinking()
}

func ledBlinking() {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "7")

	work := func() {
		gobot.Every(2*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkLED",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}

func initlog(path string) {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)

	log.SetOutput(&lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     365,
		LocalTime:  true,
		Compress:   true,
	})
}
