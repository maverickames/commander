package commander

/*
This package provides a way to pass any command to the Commander for scheduling.
*/

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Command interface {
	Start() error
	// TaskId() string
	ScheduledTime() time.Time
	Completed() bool
}

type Response struct {
	CmdId   string
	RunTime time.Time
	Err     error
}

type task struct {
	ctx context.Context
	cmd Command
}

type Commander struct {
	taskId    TaskId
	Disptach  chan task
	Responder chan Response
}

type TaskId string

// Returns a new Commander
func NewCommander() *Commander {
	return &Commander{
		taskId:    "id",
		Disptach:  make(chan task),
		Responder: make(chan Response),
	}
}

// Run will run in the background to manage all task.
func (cmdRec *Commander) Run() {
	defer close(cmdRec.Responder)
	for {
		select {
		case tsk, ok := <-cmdRec.Disptach:
			if !ok {
				return
			}

			// var id TaskId = "id"
			ctxVal := tsk.ctx.Value(cmdRec.taskId)
			if ctxVal == "" {
				fmt.Println("context value not found value:", ctxVal)
			}

			// Run or queue task.
			if tsk.cmd.ScheduledTime().Before(time.Now()) {
				err := tsk.cmd.Start()
				cmdRec.Responder <- Response{
					CmdId:   ctxVal.(string),
					RunTime: time.Now(),
					Err:     err,
				}
			} else {
				go func() {
					for {
						select {
						case <-time.After(tsk.cmd.ScheduledTime().Sub(time.Now())):
							err := tsk.cmd.Start()
							cmdRec.Responder <- Response{
								CmdId:   ctxVal.(string),
								RunTime: time.Now(),
								Err:     err,
							}
							return
						case <-tsk.ctx.Done():
							cmdRec.Responder <- Response{
								CmdId:   ctxVal.(string),
								RunTime: time.Now(),
								Err:     errors.New("Context Cancel Request"),
							}
							return
						}
					}
				}()
			}
		}
	}
}

// Add task to the dispatcher
func (cmdRec *Commander) Add(ctx context.Context, cmd Command) bool {
	cmdRec.Disptach <- task{
		ctx: ctx,
		cmd: cmd,
	}
	return true
}

// Halts commander. Any task schedule will still complete.
func (cmdRec *Commander) Halt() {
	close(cmdRec.Disptach)
}
