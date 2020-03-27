package executor

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"log"
	"os/exec"
	"strings"
	"sync"
)

var ok = false
var wg sync.WaitGroup
var newErr error

func FormatOutput(resName, output string) string {
	return fmt.Sprintf("%s: %s", resName, output)
}
func execCommand(params string, task packet.Task, action string, sendMessage func(packet packet.Packet)) error {

	cmd := exec.Command(commandName, "-c", params)
	if action == "script" {
		cmd = exec.Command(commandName, params)
	}

	//显示运行的命令
	fmt.Printf("执行命令: %s\n", strings.Join(cmd.Args[1:], " "))
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Println("exec the cmd failed")
		return err
	}

	scan := bufio.NewScanner(stderr)
	//go func() {
	for scan.Scan() {
		ok = false
		s := scan.Text()
		if !ok {
			newErr = errors.New(s)
		}
		//sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
		//	TaskId: task.Id,
		//	Code:   unknownError,
		//	Output: s,
		//}))
	}
	//}()

	scanner := bufio.NewScanner(stdout)
	//go func() {
	for scanner.Scan() {
		ok = true
		newErr = nil
		// send message
		s := scanner.Text()
		fmt.Println("this is success: " + s)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: FormatOutput(task.ResourceName, s),
		}))
	}
	//if ok {
	//	sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
	//		TaskId:    task.Id,
	//		Output:    FormatOutput(task.ResourceName, "SUCCESS"),
	//		Completed: true,
	//	}))
	//}

	//}()

	cmd.Wait()
	return newErr
}
