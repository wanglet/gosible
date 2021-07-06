package gosible

import "fmt"

type EventHandler struct {
	UUID    string
	channel chan map[string]string
	silent  bool
}

func (e EventHandler) Write(p []byte) (n int, err error) {
	event := make(map[string]string)
	if e.channel != nil {
		event["uuid"] = e.UUID
		event["out"] = string(p)
		e.channel <- event
	}

	if !e.silent {
		fmt.Printf("%v %v", e.UUID, string(p))
	}
	return len(p), nil
}
