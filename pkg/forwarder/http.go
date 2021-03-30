package forwarder

import (
	"io/ioutil"
	"net/http"

	"github.com/anchnet/smartops-agent/pkg/executor"
	"github.com/cihub/seelog"
)

const (
	LocalMetricListen string = "127.0.0.1:8899"
	LocalMetricHandle string = "/localmetric"
)

func (f *defaultForwarder) StartLocalHttp() {
	http.HandleFunc(LocalMetricHandle, f.indexHandler)
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
