package packet

type Task struct {
	Id      string      `json:"id"`
	NodeId  string      `json:"node_id"`
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
