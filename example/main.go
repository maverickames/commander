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
	WSConn         *os.File  `json:"-"` // This is jsut temporary until its tie into the websockets.
	CmdId          string    // Command GUID
	SchedTime      time.Time // Planned Runtime
	ExcTime        time.Time // Actuall Execution Time
	DispatchedTime time.Time // Dispatched and queued
	Cmd            string    // Json data
	//response       string    // Response from Sever //Still paying with this idea. Need more time to work it out.
}

func main() {
	fmt.Println("Commander")

	cmdr := commander.NewCommander()
	go cmdr.Run()

	var job commander.Command = command{
		WSConn:    os.Stdout,
		CmdId:     "command1",
		SchedTime: time.Now().Add(20 * time.Second),
		Cmd:       "do stuff",
	}
	cmdr.Add(job)

	jobs := cmdr.GetJobs(true)
	for _, j := range jobs {
		fmt.Println("Current job list:" + j.(command).CmdId)
	}

	fmt.Println(cmdr.GetJob("command1").(command))

	// These are just here for example sake.
	// We would have a blocking call here irl.
	time.Sleep(30 * time.Second)
	cmdr.Halt()
	time.Sleep(5 * time.Second)
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

func (cmd command) TaskId() string {
	return cmd.CmdId
}

func (cmd command) ScheduledTime() time.Time {
	return cmd.SchedTime
}

// Completed returns bool true if not nil assignment and exctime.
func (cmd command) Completed() bool {
	return !cmd.ExcTime.IsZero()
}
