package helpers

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jenujari/go-srv-bootstrap/config"
)

type Commander struct {
	interuppted bool

	done           chan bool
	interrupt      chan os.Signal
	FatalErrorChan chan error

	CTX    context.Context
	cancel context.CancelFunc

	WG *sync.WaitGroup
}

func NewCommander() *Commander {
	cmd := new(Commander)
	cmd.CTX, cmd.cancel = context.WithCancel(context.Background())
	cmd.WG = new(sync.WaitGroup)
	cmd.done = make(chan bool)
	cmd.FatalErrorChan = make(chan error)
	cmd.interrupt = make(chan os.Signal)
	cmd.interuppted = false

	return cmd
}

func (cmder *Commander) WaitForFinish() {
	go handleInterrupt(cmder)
	go waitGroupDone(cmder)
	go watchError(cmder)
	gracefullExit(cmder)
}

func handleInterrupt(cmder *Commander) {
	// system inturrpt signal or terminate signal will be passed on interrupt channel 
	signal.Notify(cmder.interrupt, syscall.SIGINT, syscall.SIGTERM)

	for range cmder.interrupt {
		if cmder.interuppted {
			config.Log.Println("\nInterrupt signal already captured working on closing the process.")
			continue
		}
		cmder.interuppted = true
		cmder.cancel()
		config.Log.Println("Interuppt signal captured.")
	}
}

func watchError(cmder *Commander) {
	err := <-cmder.FatalErrorChan
	config.Log.Println("Fatal error captured :: ", err)
	cmder.cancel()
}

func waitGroupDone(cmder *Commander) {
	cmder.WG.Wait()
	cmder.done <- true
}

func gracefullExit(cmder *Commander) {
	<-cmder.done
	config.Log.Println("Gracefull exit")
	os.Exit(0)
}