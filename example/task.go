package main

import (
	"context"
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

	Cancel context.CancelFunc
}

func (app *App) addTask(task command) {

	task.id = "id"
	ctx, cancel := context.WithCancel(context.Background())
	ctxVal := context.WithValue(ctx, task.id, task.CmdId)
	task.Cancel = cancel

	// Add to commander
	var newTask commander.Command = task
	app.cmdr.Add(ctxVal, newTask)

	// Add to local task list.
	app.task[task.CmdId] = task
}

func (app *App) handleResponder() {
	for taskResp := range app.cmdr.Responder {
		// Handle errors
		if taskResp.Err != nil {
			fmt.Printf("Updating Task Details for: %s\n  -- Error: %s\n", taskResp.CmdId, taskResp.Err.Error())
		} else {
			fmt.Printf("Updating Task Details for: %s\n  -- Error: %s\n", taskResp.CmdId, "no errors")
		}
		app.updateTask(taskResp)
	}
}

// Update the task with the reponse dataß
func (app *App) updateTask(resp commander.Response) {
	cmd := app.getTask(resp.CmdId)
	cmd.ExcTime = resp.ExcTime
	cmd.DispatchedTime = resp.SchTime
	app.task[resp.CmdId] = cmd
}

func (cmd command) cancel() {
	cmd.Cancel()
}

func (cmd command) Start() (time.Time, error) {
	exctime := time.Now()
	// jsonData, err := json.Marshal(cmd)
	// if err != nil {
	// 	return exctime, err
	// }
	_, err := cmd.WSConn.Write([]byte("Running command" + cmd.Cmd + " \n"))
	return exctime, err
}

func (cmd command) ScheduledTime() time.Time {
	return cmd.SchedTime
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
