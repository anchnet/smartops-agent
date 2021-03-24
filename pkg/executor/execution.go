package executor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	err := cmd.Start()
	if err != nil {
		sendSuccess(task, FormatOutput(task.ResourceName, err.Error()), sendMessage)
		return
	}
	err = stdReadOnce(stdoutIn, stderrIn, task, sendMessage)
	if err != nil {
		sendSuccess(task, FormatOutput(task.ResourceName, err.Error()), sendMessage)
		return
	}
	err = cmd.Wait()
	if err != nil {
		seelog.Infof("err: %s", err)
	}
	if errStdout != nil {
		seelog.Infof("errStdout : ", errStdout)
	}
	if errStderr != nil {
		seelog.Infof("errStderr: ", errStderr)
	}
}

//stdReadOnce  stdReadOnce read pipe stdout and stderr, gather the  message and send.
func stdReadOnce(stdout io.Reader, stderr io.Reader, task packet.Task, sender func(packet packet.Packet)) error {
	stdourBuffer := make([]byte, 4*1024)
	stdErrBuffer := make([]byte, 4*1024)
	var stdourErr, stdErrErr error
	var sendStr string
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		buffer := bufio.NewReader(stdout)
		cmd.Process.Wait()
		var count int
		count, stdourErr = buffer.Read(stdourBuffer)
		seelog.Infof("Stdout read finish: count：%d", count)
	}()
	go func() {
		defer wg.Done()
		buffer := bufio.NewReader(stderr)
		cmd.Process.Wait()
		var count int
		count, stdErrErr = buffer.Read(stdErrBuffer)
		seelog.Infof("Stderr read finish: count：%d", count)
	}()
	wg.Wait()
	//
	if stdourErr == io.EOF && stdErrErr == io.EOF {
		// FormatOutput(task.ResourceName, "SUCCESS")
		sendSuccess(task, "", sender)
		return nil
	}

	if stdourErr != nil && stdourErr != io.EOF {
		return errors.New(fmt.Sprintf("Read stdout pipe error : %s", stdourErr))
	} else {
		if stdourErr != io.EOF {
			sendStr = sendStr + formatOutputGather(task.ResourceName, stdourBuffer)
		}
	}

	if stdErrErr != nil && stdErrErr != io.EOF {
		return errors.New(fmt.Sprintf("Read stderr pipe error : %s", stdErrErr))
	} else {
		if stdErrErr != io.EOF {
			sendStr = sendStr + formatOutputGather(task.ResourceName, stdErrBuffer)
		}
	}
	sendSuccess(task, sendStr, sender)
	return nil
}

func formatOutputGather(resourceName string, elems []byte) string {
	arrElems := bytes.Split(elems, []byte("\n"))
	strs := make([]string, 0)
	for _, val := range arrElems {
		strs = append(strs, FormatOutput(resourceName, string(val)))
	}
	str := strings.Join(strs, "<br>")
	return str
}

func sendSuccess(task packet.Task, str string, sender func(packet packet.Packet)) {
	str += fmt.Sprintf("<br>%s", FormatOutput(task.ResourceName, "SUCCESS"))
	seelog.Infof("Send to server,  task: %v ， str： %s", task, str)
	sender(packet.NewTaskResultPacket(packet.TaskResult{
		TaskId:    task.Id,
		Output:    str,
		Completed: true,
	}))
}

func stdRead(reader io.Reader, code int, task packet.Task, sender func(packet packet.Packet)) error {
	count := 0
	fmt.Printf("code is : %d \n", code)
	//buf := make([]byte, 1024, 1024)
	buffer := bufio.NewReader(reader)
	for {
		//n, err := reader.Read(buf[:])
		buffers, _, err := buffer.ReadLine()
		if runtime.GOOS == "windows" {
			buffers, _ = GbkToUtf8(buffers)
		}
		if err != nil {
			//if return is.EOF and count =0, then maybe we are running a eary and quick command
			//then in this time this command has completed, we should return STD_SUCCESS
			if err == io.EOF && count == 0 {
				sendCommandLineMessage(STD_SUCCESS, task, nil, sender)
			}

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
		if string(buffers) == "" {
			count++
			continue
		}
		sendCommandLineMessage(code, task, buffers, sender)
		if code == STD_READ && len(buffers) > 0 {
			count++
		}
	}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func sendCommandLineMessage(code int, task packet.Task, buffers []byte, sender func(packet packet.Packet)) {
	seelog.Infof("Send to server,  code: %d , return result: %s.", code, string(buffers))
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
