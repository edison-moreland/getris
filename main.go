package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(
		internalScreenX,
		internalScreenY,
		"Getris",
	)

	// Tetris uses esc to pause, so rebind window close
	rl.SetExitKey(rl.KeyQ)

	// Camera puts (0, 0) at the center of the screen
	camera := rl.NewCamera2D(
		rl.NewVector2(float32(internalScreenX/2), float32(internalScreenY/2)),
		rl.NewVector2(0.0, 0.0),
		0.0, 1.0,
	)

	inputEventChannel := make(chan InputEvent)
	defer close(inputEventChannel)

	game := newGameState()
	game.Run(inputEventChannel)

	rl.SetTargetFPS(60)
	for (!rl.WindowShouldClose()) && (!game.IsDone) {
		InputForwarder(inputEventChannel)

		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)
		rl.BeginMode2D(camera)

		game.Draw()

		rl.EndMode2D()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func drawCoordDebug(x, y, size int32) {
	rl.DrawLine(x, y, x+size, y, rl.Red)
	rl.DrawLine(x, y, x, y+size, rl.Green)
}
