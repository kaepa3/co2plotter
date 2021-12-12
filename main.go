package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/kaepa3/co2plotter/pkg/blinkled"
	"github.com/kaepa3/co2plotter/pkg/co2loader"
	"github.com/kaepa3/co2plotter/pkg/oled"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// scaan for terminal
var sc = bufio.NewScanner(os.Stdin)

// main runc
func main() {
	initlog("log.txt")
	log.Info("start")
	robot := blinkled.CreteBlinkingRobot(38)
	go robot.Start()

	done := make(chan struct{})
	conf := co2loader.Config{IsMock: false, Value: 2, Interval: 1}
	co2Chan, errChan := co2loader.CreateCo2Loader(conf, done)

	log.Info("loop")
	go doneLoop(done)
	mainLoop(done, co2Chan, errChan)

	if err := robot.Stop(); err != nil {
		log.Error(err)
	}
	log.Info("app end")
}

// get co2 data and display
func mainLoop(done <-chan struct{}, co2Chan <-chan int, errChan <-chan string) {
	oled := oled.CreateOled()
Polling:
	for {
		select {
		case v := <-co2Chan:
			oled.Display("co2:" + strconv.Itoa(v))
		case v := <-errChan:
			log.Println("err:", v)
		case <-done:
			break Polling
		}
	}
	oled.Display("end")
}

// loop for get text and app end
func doneLoop(done chan struct{}) {
	for {
		if sc.Scan() {
			t := sc.Text()
			fmt.Println("input:" + t)
			if t == "e" {
				break
			}
		}
	}
	close(done)
}

// logging init
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
