package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maverickames/commander"
)

type command struct {
	// WSConn         *websocket.Conn `json:"-"`
	WSConn         *os.File  `json:"-"`
	CmdId          string    // Command GUID
	SchedTime      time.Time // Planned Runtime
	ExcTime        time.Time // Actuall Execution Time
	DispatchedTime time.Time // Dispatched and queued
	Cmd            string    // Json data
	response       string    // Response from Sever
}

func main() {
	fmt.Println("Commander")

	cmdr := commander.NewCommander()
	go cmdr.Run()

	var job commander.Command = command{
		WSConn:    os.Stdout,
		CmdId:     "command1",
		SchedTime: time.Now().Add(30 * time.Second),
		Cmd:       "do stuff",
	}
	cmdr.Add(job)

	time.Sleep(1 * time.Minute)
	cmdr.Halt()
	time.Sleep(10 * time.Second)
}

func (cmd command) TaskId() string {
	return cmd.CmdId
}

func (cmd command) Start() {

	jsonData, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println(err)
		// should handle errors over an erros chanell to provide a path up the app to the gui
		//return app.em.ErrorHandler(err, nil, r.RequestURI)
	}
	cmd.WSConn.Write(jsonData)
}

func (cmd command) ScheduledTime() time.Time {

	return cmd.SchedTime
}

// Completed returns bool true if not nil assignment and exctime.
func (cmd command) Completed() (time.Time, bool) {

	return cmd.ExcTime, !cmd.ExcTime.IsZero()
}