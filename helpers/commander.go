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
	interrupted bool

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
	cmd.interrupted = false

	return cmd
}

func (cmder *Commander) AddWorker(n int) {
	cmder.WG.Add(n)
}

func (cmder *Commander) CompleteOneWorker() {
	cmder.WG.Done()
}

func (cmder *Commander) WaitForFinish() {
	go cmder.handleInterrupt()
	go cmder.waitGroupDone()
	go cmder.watchError()
	cmder.gracefullyExit()
}

func (cmder *Commander) handleInterrupt() {
	// system inturrpt signal or terminate signal will be passed on interrupt channel
	signal.Notify(cmder.interrupt, syscall.SIGINT, syscall.SIGTERM)

	for range cmder.interrupt {
		if cmder.interrupted {
			config.GetLogger().Println("\nInterrupt signal already captured working on closing the process.")
			continue
		}
		cmder.interrupted = true
		cmder.cancel()
		config.GetLogger().Println("Interuppt signal captured.")
	}
}

func (cmder *Commander) watchError() {
	err := <-cmder.FatalErrorChan
	config.GetLogger().Println("Fatal error captured :: ", err)
	cmder.cancel()
}

func (cmder *Commander) waitGroupDone() {
	cmder.WG.Wait()
	cmder.done <- true
}

func (cmder *Commander) gracefullyExit() {
	<-cmder.done
	config.GetLogger().Println("Gracefully exit")
	os.Exit(0)
}
