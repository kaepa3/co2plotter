package main

import (
	"fmt"
	"time"

	"github.com/kebhr/mhz19"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	initlog("log.txt")
	log.Info("start")
	robot := creteBlinkingRobot()
	go robot.Start()

	go func() {
		co2Robot := mhz19.MHZ19{}
		if err := co2Robot.Connect(); err != nil {
			log.Error(err)
			return
		}
		fmt.Println("start")
		for cnt := 0; cnt < 8; cnt++ {
			val, err := co2Robot.ReadCO2()
			if err != nil {
				log.Error(err)
			} else {
				logstr := fmt.Sprintf("co2:%d", val)
				fmt.Println(logstr)
				log.Info(logstr)
			}
			time.Sleep(time.Second * 1)
		}
		fmt.Println("end")
	}()

	time.Sleep(time.Minute * 1)
	if err := robot.Stop(); err != nil {
		log.Error(err)
	}
	log.Info("app end")
}

func creteBlinkingRobot() *gobot.Robot {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "37")

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
