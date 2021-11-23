package co2loader

import (
	"fmt"
	"time"

	"github.com/kaepa3/mhz19"
)

type Config struct {
	IsMock   bool
	Value    int
	Interval int
}

func CreateCo2Loader(conf Config, done <-chan struct{}) (<-chan int, <-chan string) {
	co2Chan := make(chan int)
	errChan := make(chan string)
	go Polling(done, co2Chan, errChan, conf)
	return co2Chan, errChan
}

func createRobot(conf Config) mhz19.Co2Dataloader {
	var loader mhz19.Co2Dataloader
	if conf.IsMock {
		loader = &mhz19.MockMHZ19{Value: conf.Value}
	} else {
		loader = &mhz19.MHZ19{}
	}
	return loader
}

func Polling(done <-chan struct{}, co2Chan chan<- int, errChan chan<- string, conf Config) {
	co2Robot := createRobot(conf)
	if err := co2Robot.Connect(); err != nil {
		errChan <- err.Error()
		return
	}
polling:
	for {
		select {
		case <-done:
			break polling
		default:
		}
		val, err := co2Robot.ReadCO2()
		if err != nil {
			errChan <- err.Error()
		} else {
			co2Chan <- val
		}
		time.Sleep(time.Second * time.Duration(conf.Interval))
	}
	fmt.Println("co2loader end")
}
