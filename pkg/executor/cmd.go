package executor

import (
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"os/exec"
	"strings"
)

const SUCCESS int = 0
const DefaultError int = -1

func ExecCommand(task packet.Task, sendMessage func(packet packet.Packet)) {
	if task.Content == nil || strings.Trim(task.Content.(string), "") == "" {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: "task content is empty",
			Code:   DefaultError,
		}))
		return
	}
	cnt := task.Content.(string)
	lines := strings.Split(cnt, "\n")
	out, err := exec.Command(lines[0]).Output()
	if err != nil {
		e := err.(*exec.ExitError)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: string(e.Stderr),
			Code:   e.ExitCode(),
		}))
		_ = seelog.Errorf("run cmd error,%v", e)
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
		TaskId: task.Id,
		Output: "success",
		Code:   SUCCESS,
	}))
}
