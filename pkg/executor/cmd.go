package executor

import (
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"os/exec"
	"strings"
)

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
	lines := strings.Split(cnt, "\n")
	out, err := exec.Command(lines[0]).Output()
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
	lines = strings.Split(string(out), "\n")
	for _, line := range lines {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: line,
		}))
	}
	sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
		TaskId:    task.Id,
		Output:    "success",
		Completed: true,
	}))
}
