package main

import (
	"os"
	"os/signal"
	"sync"
	
	"./room"

	"github.com/Balhazraell/logger"
)

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

func main() {
	logger.InitLogger()
	room.StartNewRoom(0)
	waitForCtrlC()
}
