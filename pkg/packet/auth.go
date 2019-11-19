package packet

type Authorize struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
