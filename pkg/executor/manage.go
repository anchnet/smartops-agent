package executor

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
)

// sender := func(packet packet.Packet) {
// 	fmt.Println(packet)
// }
// routineManage := executor.NewRoutineManage("", sender)
// routineManage.Init()
// go routineManage.Go("./a.sh", packet.Task{
// 	Id: "ddddd",
// })
// go func() {
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("into stop")
// 	routineManage.Stop("ddddd")
// }()

type ManageTask struct {
	Task      packet.Task
	CtxCancel context.Context
	PID       int
	Params    string
	Cmd       *exec.Cmd
}

type RoutineManage struct {
	TaskMap      map[string]ManageTask //key is task id
	SnapFilepath string
	SnapTicker   *time.Ticker
	Sender       func(packet packet.Packet)
	Mu           sync.RWMutex
}

func (r *RoutineManage) SetTaskMap(id string, mngTask ManageTask) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.TaskMap[id] = mngTask
}

func NewRoutineManage(snapPath string, sender func(packet packet.Packet)) *RoutineManage {
	if snapPath == "" {
		snapPath = "./snapPath.json"
	}
	return &RoutineManage{
		TaskMap:      make(map[string]ManageTask),
		SnapFilepath: snapPath,
		SnapTicker:   time.NewTicker(time.Second * 10),
		Sender:       sender,
	}
}

//Init init manage data
func (r *RoutineManage) Init() error {
	err := r.Reload()
	if err != nil {
		seelog.Error("Reload json file err : ", err)
		return err
	}

	// start snap file.
	go func() {
		for {
			seelog.Debug("Start snap file.")
			<-r.SnapTicker.C
			if err := r.Snap(); err != nil {
				seelog.Error("Snap err: ", err)
			}
		}
	}()

	//Start stopped task.
	for _, taskMng := range r.TaskMap {
		go r.Go(taskMng.Params, taskMng.Task)
	}

	return nil
}

//Go start task .
func (r *RoutineManage) Go(params string, task packet.Task) {
	seelog.Infof("start run task id: %s, params: %s", task.Id, params)
	mngTask := ManageTask{}
	cmd = exec.Command(commandName, params)

	err := cmd.Start()
	if err != nil {
		seelog.Error(err)
		return
	}

	{ // init manage task
		mngTask.PID = cmd.Process.Pid
		mngTask.Params = params
		mngTask.Cmd = cmd
		mngTask.Task = task
	}

	r.SetTaskMap(task.Id, mngTask)

	err = cmd.Wait()
	if err != nil {
		seelog.Error("manage routine cmd wait err: ", err)
	}

	//remote task
	seelog.Infof("Custom task end %s: ", task.Id)
	delete(r.TaskMap, task.Id)
}

//Stop task with task id .
func (r *RoutineManage) Stop(id string) {
	mngTask, ok := r.TaskMap[id]
	if !ok {
		seelog.Info("Task not exist!")
		return
	}

	cmd = mngTask.Cmd
	if err := cmd.Process.Kill(); err != nil {
		seelog.Error("Stop task error! ", err)
		return
	}
}

//Snap store data into local file
func (r *RoutineManage) Snap() error {
	byts, err := json.Marshal(r.TaskMap)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.SnapFilepath, byts, 0666)
}

//Reload reload data into local file
func (r *RoutineManage) Reload() error {
	jsonFile, err := os.Open(r.SnapFilepath)
	if os.IsNotExist(err) {
		seelog.Info("Snap json file not found . ")
		return nil
	} else if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	json.Unmarshal(byteValue, &r.TaskMap)
	return nil
}
