package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maverickames/commander"
)

type command struct {
	id commander.TaskId
	// WSConn         *websocket.Conn `json:"-"`
	WSConn         *os.File  `json:"-"` // This is jsut temporary until its tie into the websockets.
	CmdId          string    // Command GUID
	SchedTime      time.Time // Planned Runtime
	ExcTime        time.Time // Actuall Execution Time
	DispatchedTime time.Time // Dispatched and queued
	Cmd            string    // Json data
	//response       string    // Response from Sever //Still paying with this idea. Need more time to work it out.
}

func (app *App) addTask(task command) context.CancelFunc {
	app.task[task.CmdId] = task
	task.id = "id"
	ctx, cancel := context.WithCancel(context.Background())
	ctxVal := context.WithValue(ctx, task.id, task.CmdId)
	var newTask commander.Command = task
	app.cmdr.Add(ctxVal, newTask)
	return cancel
}

func (app *App) handleResponder() {
	for taskResp := range app.cmdr.Responder {
		fmt.Println(taskResp)
		fmt.Println("Updating Task Details for: ", taskResp.CmdId)
	}
}

func (cmd command) Start() error {
	jsonData, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	_, err = cmd.WSConn.Write(jsonData)
	return err
}

func (cmd command) ScheduledTime() time.Time {
	return cmd.SchedTime
}

func (cmd command) TaskId() string {
	return cmd.CmdId
}

// Completed returns bool true if not nil assignment and exctime.
func (cmd command) Completed() bool {
	return !cmd.ExcTime.IsZero()
}

// Get single task
func (app *App) getTask(cmdId string) command {
	return app.task[cmdId]
}

// Return all tasks currently managed by the Commander
func (app *App) getTasks(all bool) (taskList []command) {

	for _, task := range app.task {
		if all {
			taskList = append(taskList, task)
		} else if !task.Completed() {
			taskList = append(taskList, task)
		}
	}
	return
}

// Remove task after its completed its life cycle.
func (app *App) delTask(cmdId string) (taskList command) {
	taskList = app.task[cmdId]
	delete(app.task, cmdId)
	return
}
