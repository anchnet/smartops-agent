package executor

import (
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"strings"
)

func ExecCommand(task packet.Task, sendMessage func(packet packet.Packet)) {
	if task.Content == nil || strings.Trim(task.Content.(string), " ") == "" {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: FormatOutput(task.ResourceName, "task content is empty"),
			Code:   contentEmptyError,
		}))
		return
	}
	cnt := task.Content.(string)
	cmdLine := strings.Split(cnt, "\n")[0]

	execCommand(cmdLine, task, "cmd", sendMessage)
	//if err != nil {
	//	result := packet.TaskResult{
	//		TaskId: task.Id,
	//	}
	//	switch e := err.(type) {
	//	case *exec.ExitError:
	//		result.Code = e.ExitCode()
	//		result.Output = FormatOutput(task.ResourceName, string(e.Stderr))
	//		break
	//	default:
	//		result.Code = unknownError
	//		result.Output = FormatOutput(task.ResourceName, e.Error())
	//	}
	//	sendMessage(packet.NewTaskResultPacket(result))
	//	_ = seelog.Errorf("run cmd error,%v", err)
	//	return
	//}
	//sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
	//	TaskId:    task.Id,
	//	Completed: true,
	//}))
	seelog.Infof("Task %s completed.", task.Id)
}
