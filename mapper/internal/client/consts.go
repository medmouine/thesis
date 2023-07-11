package client

import "fmt"

type Topic string

const (
	UpdateStateTopic Topic = "device/%s/state/update"
)

func (t Topic) Fmt(id string, args ...interface{}) string {
	switch t {
	case UpdateStateTopic:
		return fmt.Sprintf(string(t), id)
	default:
		return fmt.Sprintf(string(t), id, args)
	}
}
