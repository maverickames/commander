package commander

/*
This package provides a way to pass any command to the Commander for scheduling.
*/

import (
	"fmt"
	"time"
)

type Command interface {
	Start()
	ScheduledTime() time.Time
	TaskId() string
	Completed() bool
}

type Commander struct {
	Commands map[string]Command
	Disptach chan Command
	Done     chan bool
}

// Returns a new Commander
func NewCommander() Commander {
	return Commander{
		Commands: make(map[string]Command),
		Disptach: make(chan Command),
		Done:     make(chan bool),
	}
}

// Run will run in the background to manage all task.
func (cmdRec *Commander) Run() {
	defer fmt.Println("Exiting Commander") // this is just for testing.
	for {
		select {
		case cmd := <-cmdRec.Disptach:

			cmdRec.Commands[cmd.TaskId()] = cmd
			// I know I could prob make this a function.
			// Should I? I feel like it could stay here for more verboseness
			// but I can underand wanting to move it out as well
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

// Add task to the dispatcher
func (cmdRec *Commander) Add(cmd Command) bool {
	cmdRec.Disptach <- cmd
	return true
}

// Get single task
func (cmdRec *Commander) GetJob(cmdId string) interface{} {
	return cmdRec.Commands[cmdId]
}

// Return all tasks currently managed by the Commander
func (cmdRec *Commander) GetJobs(all bool) (cmdList []interface{}) {

	for _, cmd := range cmdRec.Commands {
		if all {
			cmdList = append(cmdList, cmd)
		} else if !cmd.Completed() {
			cmdList = append(cmdList, cmd)
		}
	}
	return
}

// Remove task after its completed its life cycle.
func (cmdRec *Commander) DelJob(cmdId string) (cmdList interface{}) {
	cmdList = cmdRec.Commands[cmdId]
	delete(cmdRec.Commands, cmdId)
	return
}

// Halts commander. Any task schedule will still complete.
func (cmdRec *Commander) Halt() {
	cmdRec.Done <- true
}
