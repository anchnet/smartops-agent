package forwarder

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anchnet/smartops-agent/cmd/common"
	"github.com/anchnet/smartops-agent/pkg/executor"
	"github.com/anchnet/smartops-agent/pkg/packet"
	"github.com/cihub/seelog"
)

func (f *defaultForwarder) StartLocalHttp() {
	http.HandleFunc(common.LocalMetricHandle, f.indexHandler)
	http.HandleFunc(common.PhysicalDeviceHandle, f.physical_device)
	http.ListenAndServe(common.LocalMetricListen, nil)
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
