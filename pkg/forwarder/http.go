package forwarder

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cihub/seelog"
)

const (
	LocalMetricListen string = "127.0.0.1:8899"
	LocalMetricHandle string = "/localmetric"
)

func StartLocalHttp() {
	http.HandleFunc(LocalMetricHandle, indexHandler)
	http.ListenAndServe(LocalMetricListen, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	byts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		seelog.Error("Http server get body error: ", err)
		return
	}
	defer r.Body.Close()
	fmt.Println("---------->", string(byts))
}
