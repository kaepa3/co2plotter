package main

import (
	"github.com/kaepa3/co2plotter/pkg/blinkled"
	"github.com/kaepa3/co2plotter/pkg/co2loader"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	initlog("log.txt")
	log.Info("start")
	robot := blinkled.CreteBlinkingRobot(37)
	go robot.Start()

	done := make(chan struct{})
	conf := co2loader.Config{IsMock: true, Value: 2, Interval: 1}
	co2Chan, errChan := co2loader.CreateCo2Loader(conf, done)

	counter := 0
	for counter < 10 {
		select {
		case v := <-co2Chan:
			log.Println("co2:", v)
			counter++
		case v := <-errChan:
			log.Println("err:", v)
			counter++
		}
	}
	if err := robot.Stop(); err != nil {
		log.Error(err)
	}
	log.Info("app end")
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
