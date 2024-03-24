package sse

import (
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Event interface {
	String() string
}

func WriteKeepAlive(w io.Writer) error {
	_, err := fmt.Fprintln(w, ":keepalive")
	if err != nil {
		return err
	}
	return nil
}

func WriteEvent(w io.Writer, e Event) error {
	_, err := fmt.Fprintln(w, e)
	if err != nil {
		return err
	}
	return nil
}

type NamedEvent struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Data  string `json:"data"`
	Retry int    `json:"retry"`
}

func NewEvent(eventName string, data string) *NamedEvent {
	return &NamedEvent{
		ID:    uuid.NewString(),
		Event: eventName,
		Data:  data,
	}
}

func (e *NamedEvent) String() string {
	return fmt.Sprintf(`event: %s
data: %s
id: %s
retry: %d
`, e.Event, e.Data, e.ID, e.Retry)
}

type UnamedEvent struct {
	ID    string `json:"id"`
	Data  string `json:"data"`
	Retry int    `json:"retry"`
}

func NewUnamedEvent(data string) *UnamedEvent {
	return &UnamedEvent{
		ID:   uuid.NewString(),
		Data: data,
	}
}

func (e *UnamedEvent) String() string {
	return fmt.Sprintf(`data: %s
id: %s
retry: %d
`, e.Data, e.ID, e.Retry)
}
