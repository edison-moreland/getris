package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Input int

const (
	Input_Pause Input = iota
	Input_Hold
	Input_RotateCounterClockwise
	Input_RotateClockwise
	Input_HardDrop
	Input_SoftDrop
	Input_MoveLeft
	Input_MoveRight
)

type Action int

const (
	Action_Down Action = iota
	Action_Hold
	Action_Up
)

type InputEvent struct {
	Input   Input
	Action  Action
	KeyCode int32
}

// KeyMap maps raylib's key codes to game's input codes
var KeyMap = map[int32]Input{
	// Pause
	rl.KeyEscape: Input_Pause,
	rl.KeyF1:     Input_Pause,

	// Hold
	rl.KeyLeftShift:  Input_Hold,
	rl.KeyRightShift: Input_Hold,
	rl.KeyC:          Input_Hold,
	rl.KeyKp0:        Input_Hold,

	// Rotate Counter-Clockwise
	rl.KeyLeftControl:  Input_RotateCounterClockwise,
	rl.KeyRightControl: Input_RotateCounterClockwise,
	rl.KeyZ:            Input_RotateCounterClockwise,
	rl.KeyKp3:          Input_RotateCounterClockwise,
	rl.KeyKp7:          Input_RotateCounterClockwise,

	// Rotate Clockwise
	rl.KeyX:   Input_RotateClockwise,
	rl.KeyUp:  Input_RotateClockwise,
	rl.KeyKp1: Input_RotateClockwise,
	rl.KeyKp5: Input_RotateClockwise,
	rl.KeyKp9: Input_RotateClockwise,

	// Hard Drop
	rl.KeySpace: Input_HardDrop,
	rl.KeyKp8:   Input_HardDrop,

	// Soft Drop
	rl.KeyDown: Input_SoftDrop,
	rl.KeyKp2:  Input_SoftDrop,

	// Move Left
	rl.KeyLeft: Input_MoveLeft,
	rl.KeyKp4:  Input_MoveLeft,

	// Move Right
	rl.KeyRight: Input_MoveRight,
	rl.KeyKp6:   Input_MoveRight,
}

// InverseKeyMap maps game's input codes to raylib's key codes
var InverseKeyMap map[Input][]int32

func init() {
	InverseKeyMap = make(map[Input][]int32)
	for key, input := range KeyMap {
		InverseKeyMap[input] = append(InverseKeyMap[input], key)
	}
}

func InputForwarder(eventChannel chan InputEvent) {
	keyPressed := rl.GetKeyPressed()
	for keyPressed != 0 {
		if input, ok := KeyMap[keyPressed]; ok {
			eventChannel <- InputEvent{
				Input:   input,
				Action:  Action_Down,
				KeyCode: keyPressed,
			}
			waitForLongPress(eventChannel, keyPressed)
		}

		keyPressed = rl.GetKeyPressed()
	}

	for input, keyCodes := range InverseKeyMap {
		for _, keyCode := range keyCodes {
			if rl.IsKeyReleased(keyCode) {
				eventChannel <- InputEvent{
					Input:   input,
					Action:  Action_Up,
					KeyCode: keyCode,
				}
			}
		}
	}
}

func waitForLongPress(eventChannel chan<- InputEvent, keyCode int32) {
	time.AfterFunc(inputLongPressInterval, func() {
		if rl.IsKeyDown(keyCode) {
			eventChannel <- InputEvent{
				Input:   Input(KeyMap[keyCode]),
				Action:  Action_Hold,
				KeyCode: keyCode,
			}
		}
	})
}

func DebugInputEvent(eventChannel chan InputEvent) {
	go func() {
		for e := range eventChannel {
			var actionText string
			switch e.Action {
			case Action_Down:
				actionText = "down"
			case Action_Hold:
				actionText = "hold"
			case Action_Up:
				actionText = "up"
			}

			var inputText string
			switch e.Input {
			case Input_Pause:
				inputText = "pause"
			case Input_Hold:
				inputText = "hold"
			case Input_RotateCounterClockwise:
				inputText = "rotate counter-clockwise"
			case Input_RotateClockwise:
				inputText = "rotate clockwise"
			case Input_HardDrop:
				inputText = "hard drop"
			case Input_SoftDrop:
				inputText = "soft drop"
			case Input_MoveLeft:
				inputText = "move left"
			case Input_MoveRight:
				inputText = "move right"
			}

			fmt.Printf("Input: %s, Action: %s, KeyCode: %d\n", inputText, actionText, e.KeyCode)
		}
	}()
}
