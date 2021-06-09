package main

import (
	"fmt"
	"os"
	"time"

	"github.com/maverickames/commander"
)

type App struct {
	cmdr *commander.Commander
	task map[string]command
}

func main() {
	fmt.Println("Commander Example")

	// Setup
	app := App{}
	app.task = make(map[string]command)
	app.cmdr = commander.NewCommander()

	// Commander responder channel. Handler errors
	// process analytics.
	go app.handleResponder()

	// Commander to run in background routine to process task.
	go app.cmdr.Run()

	// Create and added task to dispatcher.
	var task1 = command{
		WSConn:    os.Stdout,
		CmdId:     "command1",
		SchedTime: time.Now().Add(20 * time.Second),
		Cmd:       "do stuff1",
	}
	cancel := app.addTask(task1)

	var task2 = command{
		WSConn:    os.Stdout,
		CmdId:     "command2",
		SchedTime: time.Now().Add(10 * time.Second),
		Cmd:       "do stuff2",
	}
	app.addTask(task2)

	var task3 = command{
		WSConn:    os.Stdout,
		CmdId:     "command3",
		SchedTime: time.Now(),
		Cmd:       "do stuff3",
	}
	app.addTask(task3)

	// Cancel task1
	cancel()

	// Check what commands have been added.
	fmt.Println(app.getTask("command1"))

	// Delete from local state after
	// processing.
	del := app.getTask("command1")
	fmt.Printf("Deleted: %s\n", del.CmdId)

	// These are just here for example sake.
	// We would have a blocking call here irl.
	time.Sleep(30 * time.Second)
	app.cmdr.Halt()
	time.Sleep(5 * time.Second)

	// Display task with results
	tasks := app.getTasks(true)
	for _, t := range tasks {
		fmt.Printf("TaskName:%s\n  --  ScheduleTime:%s\n  --  Dispatched:%s\n  --  ExecuteTime:%s\n", t.CmdId, t.SchedTime, t.DispatchedTime, t.ExcTime)
	}
}
