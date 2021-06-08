package commander

/*
This package provides a way to pass any command to the Commander for scheduling.
*/

import (
	"fmt"
	"time"
)

type Command interface {
	// Task() bool
	Start()
	ScheduledTime() time.Time
	TaskId() string
	Completed() (time.Time, bool)
}

type Commander struct {
	Commands map[string]Command
	Disptach chan Command
	Done     chan bool
}

func NewCommander() Commander {
	return Commander{
		Commands: make(map[string]Command),
		Disptach: make(chan Command),
		Done:     make(chan bool),
	}
}

func (cmdRec *Commander) Run() {
	defer fmt.Println("Exiting Commander")
	for {
		select {
		case cmd := <-cmdRec.Disptach:

			cmdRec.Commands[cmd.TaskId()] = cmd
			go func() {
				if cmd.ScheduledTime().Before(time.Now()) {
					cmd.Start()
					return
				}
				time.Sleep(cmd.ScheduledTime().Sub(time.Now()))
				cmd.Start()

			}()
		case <-cmdRec.Done:
			fmt.Println("done")
			return
		}
	}
}

func (cmdRec *Commander) Add(cmd Command) bool {
	cmdRec.Disptach <- cmd
	return true
}

func (cmdRec *Commander) Halt() {
	cmdRec.Done <- true
}
