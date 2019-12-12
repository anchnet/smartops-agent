package executor

import (
	"bytes"
	"encoding/gob"
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
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(task.Content)
	if err != nil {
		_ = seelog.Errorf("encode task content error, %v", err)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: err.Error(),
			Code:   contentEncodeError,
		}))
		return
	}
	file, err := ioutil.TempFile("/tmp/smartops-agent/var/run/", task.Id+"_*.sh")
	if err != nil {
		_ = seelog.Errorf("create script file error, %v", err)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: err.Error(),
			Code:   createFileError,
		}))
		return
	}
	defer func() {
		err = os.Remove(file.Name())
		if err != nil {
			_ = seelog.Errorf("remove file %s error, %v", file.Name(), err)
		}
	}()
	err = ioutil.WriteFile(file.Name(), buf.Bytes(), 0644)
	if err != nil {
		_ = seelog.Errorf("save script content to file %s error, %v", file, err)
		sendMessage(packet.NewTaskResultPacket(packet.TaskResult{
			TaskId: task.Id,
			Output: err.Error(),
			Code:   saveContentError,
		}))
		return
	}
	out, err := exec.Command(file.Name()).Output()
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
}
