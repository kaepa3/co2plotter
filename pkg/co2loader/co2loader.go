package co2loader

import (
	"fmt"
	"time"

	"github.com/kaepa3/mhz19"
)

func CreateCo2Loader(done <-chan struct{}) (<-chan int, <-chan string) {
	co2Chan := make(chan int)
	errChan := make(chan string)
	go func() {
		co2Robot := mhz19.MHZ19{}
		if err := co2Robot.Connect(); err != nil {
			errChan <- err.Error()
			return
		}
		fmt.Println("start")
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
			time.Sleep(time.Second * 1)
		}
		fmt.Println("co2loader end")
	}()

	return co2Chan, errChan
}
