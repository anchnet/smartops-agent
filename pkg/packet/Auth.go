package packet

type Authorize struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
