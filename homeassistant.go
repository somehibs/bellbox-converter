package converter

import (
	"encoding/json"
	"fmt"

	"git.circuitco.de/self/bellbox"
)

type HomeAssistant struct {
}

type HassMessage struct {
}

func (p HomeAssistant) Handle(input []byte) bellbox.Message {
	fmt.Printf("HASS MSG: %s\n", input)
	msg := bellbox.Message{}
	j := HassMessage{}
	err := json.Unmarshal(input, &j)
	fmt.Printf("e: %s j: %+v\n", err, j)
	if err != nil {
		return msg
	}
	msg.Message = ""
	return msg
}
