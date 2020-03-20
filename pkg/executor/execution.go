package executor

import (
	"bufio"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"log"
	"os/exec"
	"strings"
)

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
	for scan.Scan() {
		s := scan.Text()
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Code:   unknownError,
			Output: s,
		}))
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			// send message
			sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
				TaskId: task.Id,
				Output: scanner.Text(),
			}))
			fmt.Println(scanner.Text())
		}
	}()
	err := cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
