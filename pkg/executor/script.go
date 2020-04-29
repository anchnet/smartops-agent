package executor

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

var file string

func RunScript(task packet.Task, sendMessage func(p packet.Packet)) {
	if task.Content == nil || strings.Trim(task.Content.(string), "") == "" {
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: FormatOutput(task.ResourceName, "task content is empty"),
			Code:   contentEmptyError,
		}))
		return
	}
	sysType := runtime.GOOS
	if sysType == "windows" {
		file = fmt.Sprintf("C://cloudops-agent/var/cache/%s.ps1", task.Id)
	} else {
		file = fmt.Sprintf("/opt/cloudops-agent/var/cache/%s.sh", task.Id)
	}
	//file := fmt.Sprintf("/Users/james/scripts/smartops-agent/var/cache/%s.sh", task.Id)
	err := ioutil.WriteFile(file, []byte(task.Content.(string)), 0744)
	if err != nil {
		_ = seelog.Errorf("save script content to file %s error, %v", file, err)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: FormatOutput(task.ResourceName, err.Error()),
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
	if runtime.GOOS != "windows" {
		err = exec.Command("chmod", "+x", file).Run()
		if err != nil {
			sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
				TaskId: task.Id,
				Output: FormatOutput(task.ResourceName, "update file permission error"),
				Code:   unknownError,
			}))
			return
		}
	}

	execCommand(file, task, "script", sendMessage)
	//if err != nil {
	//	result := packet.TaskResult{
	//		TaskId: task.Id,
	//	}
	//	switch e := err.(type) {
	//	case *exec.ExitError:
	//		result.Code = e.ExitCode()
	//		result.Output = FormatOutput(task.ResourceName, string(e.Stderr))
	//		break
	//	case *exec.Error:
	//		result.Code = unknownError
	//		result.Output = FormatOutput(task.ResourceName, e.Error())
	//	case *os.PathError:
	//		fmt.Println(e.Err)
	//	default:
	//		result.Code = unknownError
	//		result.Output = FormatOutput(task.ResourceName, e.Error())
	//	}
	//	sendMessage(packet.NewTaskResultPacket(result))
	//	_ = seelog.Errorf("run cmd error,%v", err)
	//	return
	//}
	seelog.Infof("Task %s completed.", task.Id)
}
