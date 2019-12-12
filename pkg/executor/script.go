package executor

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func RunScript(task packet.Task, sendMessage func(p packet.Packet)) {
	if task.Content == nil || strings.Trim(task.Content.(string), "") == "" {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: "task content is empty",
			Code:   contentEmptyError,
		}))
		return
	}
	file := fmt.Sprintf("/opt/smartops-agent/var/cache/%s.sh", task.Id)
	err := ioutil.WriteFile(file, []byte(task.Content.(string)), 0744)
	if err != nil {
		_ = seelog.Errorf("save script content to file %s error, %v", file, err)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: err.Error(),
			Code:   saveContentError,
		}))
		return
	}
	defer func() {
		//err = os.Remove(file.Name())
		//if err != nil {
		//	_ = seelog.Errorf("remove file %s error, %v", file.Name(), err)
		//}
	}()
	err = exec.Command("chmod", "+x", file).Run()
	if err != nil {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: "update file permission error",
			Code:   unknownError,
		}))
		return
	}
	out, err := exec.Command("/bin/bash", file).Output()
	if err != nil {
		result := packet.TaskResult{
			TaskId: task.Id,
		}
		switch e := err.(type) {
		case *exec.ExitError:
			result.Code = e.ExitCode()
			result.Output = string(e.Stderr)
			break
		case *exec.Error:
			result.Code = unknownError
			result.Output = e.Error()
		case *os.PathError:
			fmt.Println(e.Err)
		}
		sendMessage(packet.NewTaskResultPacket(result))
		_ = seelog.Errorf("run cmd error,%v", err)
		return
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: line,
		}))
	}
	sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
		TaskId:    task.Id,
		Output:    "SUCCESS",
		Completed: true,
	}))
	seelog.Infof("Task %s completed.", task.Id)
}
