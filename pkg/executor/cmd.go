package executor

import (
	"bufio"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"io"
	"os"
	"os/exec"
	"strings"
)

var commandName = "/bin/bash"

func ExecCommand(task packet.Task, sendMessage func(packet packet.Packet)) {
	if task.Content == nil || strings.Trim(task.Content.(string), " ") == "" {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: "task content is empty",
			Code:   contentEmptyError,
		}))
		return
	}
	cnt := task.Content.(string)
	cmdLine := strings.Split(cnt, "\n")[0]
	err := execCommand(cmdLine, task, sendMessage)
	if err != nil {
		result := packet.TaskResult{
			TaskId: task.Id,
		}
		switch e := err.(type) {
		case *exec.ExitError:
			result.Code = e.ExitCode()
			result.Output = string(e.Stderr)
			break
		default:
			result.Code = unknownError
			result.Output = e.Error()
		}
		sendMessage(packet.NewTaskResultPacket(result))
		_ = seelog.Errorf("run cmd error,%v", err)
		return
	}
	sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
		TaskId:    task.Id,
		Completed: true,
	}))
	seelog.Infof("Task %s completed.", task.Id)
}

func execCommand(params string, task packet.Task, sendMessage func(packet packet.Packet)) error {
	cmd := exec.Command(commandName, "-c", params)
	//显示运行的命令
	fmt.Printf("执行命令: %s\n", strings.Join(cmd.Args[1:], " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error=>", err.Error())
		return err
	}
	cmd.Start() // Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		// send message
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: line,
		}))
		fmt.Println(line)
	}

	cmd.Wait()
	return err
}
