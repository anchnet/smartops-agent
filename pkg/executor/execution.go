package executor

import (
	"bufio"
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"io"
	"os"
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
