package listener

import (
	"context"
	"fmt"

	input "github.com/nathan-fiscaletti/dev-input"
	"github.com/nathan-fiscaletti/kbm-overlay/internal/config"
	"github.com/nathan-fiscaletti/kbm-overlay/internal/events"
	"github.com/samber/lo"
)

func Listen(ctx context.Context, cfg config.Config, outChan chan any) error {
	keyboards, err := input.ListKeyboards()
	if err != nil {
		return err
	}

	if len(keyboards) == 0 {
		return fmt.Errorf("No keyboards found")
	}

	for _, keyboard := range keyboards {
		fmt.Printf("Keyboard found: %s\n", keyboard.Name)
		go handleKeyboardEvents(ctx, cfg, outChan, keyboard)
	}

	mice, err := input.ListPointerDevices()
	if err != nil {
		return err
	}

	if len(mice) == 0 {
		return fmt.Errorf("No mice found")
	}

	for _, mouse := range mice {
		fmt.Printf("Mouse found: %s %s\n", mouse.Name, mouse.Path)
		go handleMouseEvents(ctx, cfg, outChan, mouse)
	}
	fmt.Println()

	return nil
}

func handleKeyboardEvents(ctx context.Context, cfg config.Config, outChan chan any, keyboard *input.Device) {
	err := keyboard.Listen(ctx, func(event input.Event) {
		if lo.Contains(cfg.Monitor.Keys, event.Code) && event.Value != 2 {
			outChan <- events.KeyEvent{
				Type:  events.Key,
				Code:  event.Code,
				State: event.Value == 1,
			}
		}
	})
	if err != nil {
		fmt.Println("Error listening for keyboard events:", err)
	}
}

func handleMouseEvents(ctx context.Context, cfg config.Config, outChan chan any, mouse *input.Device) {
	var lastKnownXVal, lastKnownYVal, lastKnownWheelVal int32

	err := mouse.Listen(ctx, func(event input.Event) {
		if event.Type == input.EV_TYPE_SYN {
			return
		}

		var axis events.Axis
		var directions []events.Direction = []events.Direction{}

		if event.Type == input.EV_TYPE_REL {
			switch event.Code {
			case events.ABS_X:
				axis = events.Location
				if event.Value > 0 {
					directions = append(directions, events.DirectionRight)
				} else if event.Value < 0 {
					directions = append(directions, events.DirectionLeft)
				}
			case events.REL_Y:
				axis = events.Location
				if event.Value > 0 {
					directions = append(directions, events.DirectionDown)
				} else if event.Value < 0 {
					directions = append(directions, events.DirectionUp)
				}
			case events.REL_WHEEL:
				axis = events.Wheel
				if event.Value > 0 {
					directions = append(directions, events.DirectionUp)
				} else if event.Value < 0 {
					directions = append(directions, events.DirectionDown)
				}
			}
		} else if event.Type == input.EV_TYPE_ABS {
			switch event.Code {
			case events.ABS_X:
				axis = events.Location
				if event.Value > lastKnownXVal {
					directions = append(directions, events.DirectionRight)
				} else if event.Value < lastKnownXVal {
					directions = append(directions, events.DirectionLeft)
				}
				lastKnownXVal = event.Value
			case events.ABS_Y:
				axis = events.Location
				if event.Value > lastKnownYVal {
					directions = append(directions, events.DirectionDown)
				} else if event.Value < lastKnownYVal {
					directions = append(directions, events.DirectionUp)
				}
				lastKnownYVal = event.Value
			case events.ABS_WHEEL:
				axis = events.Wheel
				if event.Value > lastKnownWheelVal {
					directions = append(directions, events.DirectionUp)
				}
				lastKnownWheelVal = event.Value
			default:
				fmt.Printf("Unknown ABS code: %d\n", event.Code)
			}
		} else if event.Type == input.EV_TYPE_KEY {
			if lo.Contains(cfg.Monitor.MouseButtons, event.Code) {
				outChan <- events.KeyEvent{
					Type:  events.Key,
					Code:  event.Code,
					State: event.Value == 1,
				}
			}
			return
		}

		outChan <- events.MouseEvent{
			Type:       events.MouseMovement,
			Axis:       axis,
			Directions: directions,
		}
	})
	if err != nil {
		fmt.Println("Error listening for mouse events:", err)
	}
}
