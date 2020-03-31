package executor

import (
	"bufio"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"io"
	"os/exec"
	"runtime"
)

const (
	STD_READ    = 0
	STD_ERR     = 1
	STD_SUCCESS = 2
	//STD_FAIL    = 3
)

var commandName = "/bin/bash"
var powershell = "powershell"
var cmd *exec.Cmd

const WINDOWS = "windows"

func FormatOutput(resName, output string) string {
	return fmt.Sprintf("%s: %s", resName, output)
}
func execCommand(params string, task packet.Task, action string, sendMessage func(packet packet.Packet)) {
	if runtime.GOOS != WINDOWS {

		if action == "cmd" {
			cmd = exec.Command(commandName, "-c", params)
		}
		if action == "script" {
			cmd = exec.Command(commandName, params)
		}
	} else {
		cmd = exec.Command(powershell, params)
	}

	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()
	go func() {
		errStdout = stdRead(stdoutIn, STD_READ, task, sendMessage)
	}()
	go func() {
		errStderr = stdRead(stderrIn, STD_ERR, task, sendMessage)
	}()
	err := cmd.Wait()
	if err != nil {
		//log.Fatalf("cmd.Run() failed with %s\n", err)
		fmt.Println("err")
	}
	if errStdout != nil || errStderr != nil {
		//log.Fatalf("failed to capture stdout or stderr\n")
		fmt.Println("read and write error")
	}
}

func stdRead(reader io.Reader, code int, task packet.Task, sender func(packet packet.Packet)) error {
	count := 0
	fmt.Printf("code is : %d \n", code)
	//buf := make([]byte, 1024, 1024)
	buffer := bufio.NewReader(reader)
	for {
		//n, err := reader.Read(buf[:])
		buffers, _, err := buffer.ReadLine()
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			if code == STD_READ && count > 0 {
				sendCommandLineMessage(STD_SUCCESS, task, nil, sender)
			}
			return err
		}
		//fmt.Printf("this code is : %d and val : %s \n", code, string(buffers))
		sendCommandLineMessage(code, task, buffers, sender)
		if code == STD_READ && len(buffers) > 0 {
			count++
		}
	}
}

func sendCommandLineMessage(code int, task packet.Task, buffers []byte, sender func(packet packet.Packet)) {
	switch code {
	case STD_READ:
		sender(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: FormatOutput(task.ResourceName, string(buffers)),
		}))
		break
	case STD_ERR:
		// default  unknown error
		sender(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: FormatOutput(task.ResourceName, string(buffers)),
			Code:   unknownError,
		}))
		break
	case STD_SUCCESS:
		sender(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId:    task.Id,
			Output:    FormatOutput(task.ResourceName, "SUCCESS"),
			Completed: true,
		}))
		break
	}
}
