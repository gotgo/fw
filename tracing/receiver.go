package tracing

import (
	"encoding/json"

	"github.com/gotgo/fw/logging"
)

type Receiver struct {
	ReceiverName string
	Writer       TraceMessageWriter
	Log          logging.Logger `inject:""`
}

func (ft *Receiver) Name() string {
	return ft.ReceiverName
}

func (ft *Receiver) Receive(m *TraceMessage) {
	if bytes, err := json.MarshalIndent(m, "", "\t"); err != nil {
		ft.Log.MarshalFail("trace receive", err)
	} else if _, err := ft.Writer.Write(bytes); err != nil {
		//if that was a partial write, could corrupt the log
		ft.Log.Error("failed to write to receiver", err)
	}
}

func (ft *Receiver) Close() {
	if err := ft.Writer.Close(); err != nil {
		ft.Log.Error("failed to close trace writer", err)
	}
}
