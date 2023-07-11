package client

type Status string

const (
	Init       Status = "init"
	Connecting Status = "connecting"
	Connected  Status = "connected"
	ConnError  Status = "error"
)

type StateUpdate struct {
	ReportInterval string `json:"report_interval"`
}
