package staff

import hsm "github.com/hhkbp2/go-hsm"

const (
	EventBase				hsm.EventType = hsm.EventUser + 1 + iota
	EventStartOrder
	EventFinishOrder
	EventYes
	EventCancel
	EventMenu

	EventUpdate
)

var EventsToStr = map[hsm.EventType]string {
	EventStartOrder:	"StartOrder",
	EventFinishOrder:	"FinishOrder",
	EventCancel:		"Cancel",
}

func PrintEvent(eventType hsm.EventType) string {
	return EventsToStr[eventType]
}

type StaffEvent interface {
	hsm.Event
}

type GeneralEvent struct {
	*hsm.StdEvent
}

func NewIntEvent(event int) *GeneralEvent {
	return NewEvent(hsm.EventType(event))
}

func NewEvent(eventType hsm.EventType) *GeneralEvent {
	return &GeneralEvent{
		hsm.NewStdEvent(eventType),
	}
}

