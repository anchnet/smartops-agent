package forwarder

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anchnet/smartops-agent/pkg/executor"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
)

const (
	LocalMetricListen    string = "127.0.0.1:48001"
	LocalMetricHandle    string = "/localmetric"
	PhysicalDeviceHandle string = "/physical_device"
)

func (f *defaultForwarder) StartLocalHttp() {
	http.HandleFunc(LocalMetricHandle, f.indexHandler)
	http.HandleFunc(PhysicalDeviceHandle, f.physical_device)
	http.ListenAndServe(LocalMetricListen, nil)
}

func (f *defaultForwarder) indexHandler(w http.ResponseWriter, r *http.Request) {
	byts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		seelog.Error("Http server get body error: ", err)
		return
	}
	defer r.Body.Close()
	rm := executor.GetRoutineManage()
	rm.Send(byts)
}

func (f *defaultForwarder) physical_device(w http.ResponseWriter, r *http.Request) {
	byts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		seelog.Error("Http server get body error: ", err)
		return
	}
	defer r.Body.Close()
	pkg := packet.Packet{}
	pkg.Data = string(byts)
	pkg.Time = time.Now()
	pkg.Type = "physical_device"

	err = f.sendMessage(pkg)
	if err != nil {
		seelog.Error("Send to transfer error %v", err)
	}
}
