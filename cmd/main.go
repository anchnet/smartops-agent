package main

import (
	"github.com/anchnet/smartops-agent/cmd/app"
	"os"
)

func main() {
	if err := app.Command.Execute(); err != nil {
		os.Exit(-1)
	}

	//wsAddr := "wss://devtransfer.smartops.anchnet.com/ws"
	//headers := make(http.Header)
	//headers.Add("a","b")
	//wsConn, _, err := websocket.DefaultDialer.Dial(wsAddr, headers)
	//websocket.DefaultDialer.Dial(wsAddr, headers)
	//if err != nil {
	//	_ = fmt.Errorf("connect error", err)
	//	return
	//}
	//_ = wsConn.WriteMessage(websocket.TextMessage, []byte("hello"))
	//fmt.Print(wsConn.Subprotocol())

}
