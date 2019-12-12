package packet

import "fmt"

type WsResponse struct {
	Type string      `json:"type"`
	Code int32       `json:"code"`
	Body interface{} `json:"body"`
}

func (wr *WsResponse) String() string {
	return fmt.Sprintf("type: %s, code: %s, content: %s", wr.Type, fmt.Sprint(wr.Code), wr.Body)
}

type Task struct {
	Id      string      `json:"id"`
	NodeId  string      `json:"nodeId"`
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func (t *Task) String() string {
	return fmt.Sprintf("id: %s, nodeId: %s, type: %s, cnt: %s", t.Id, t.NodeId, t.Type, t.Content)
}
