package events

const (
	// Event codes for relative movement (e.g., mouse)
	REL_X     = 0x00
	REL_Y     = 0x01
	REL_WHEEL = 0x08

	// Event codes for absolute positioning (e.g., touch screen)
	ABS_X     = 0x00
	ABS_Y     = 0x01
	ABS_WHEEL = 0x08
)

type EventType string

const (
	MouseMovement EventType = "mouse_movement"
	MouseButton   EventType = "mouse_button"
	Key           EventType = "key"
)

type Axis string

const (
	Location Axis = "location"
	Wheel    Axis = "wheel"
)

type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

type MouseEvent struct {
	Type       EventType   `json:"type"`
	Axis       Axis        `json:"axis"`
	Directions []Direction `json:"directions"`
}

type KeyEvent struct {
	Type  EventType `json:"type"`
	Code  uint16    `json:"code"`
	State bool      `json:"state"`
}
