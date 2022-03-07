package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawCenteredText(text string) {
	font := rl.GetFontDefault()
	size := rl.MeasureTextEx(
		font,
		text,
		titleTextSize,
		titleTextSpacing,
	)

	rl.DrawTextEx(
		font,
		text,
		rl.Vector2{
			X: -(size.X / 2),
			Y: -(size.Y / 2),
		},
		titleTextSize,
		titleTextSpacing,
		textColor,
	)
}

func drawBorderedRectangle(x, y, width, height int32, backgroundColor, outlineColor rl.Color) {
	rl.DrawRectangle(x, y, width, height, backgroundColor)
	rl.DrawRectangleLines(x, y, width, height, outlineColor)
}

// Draw is intended to be called from the render loop
type drawBoardCellCallback func(gridX, gridY, screenX, screenY int32) (color rl.Color, cellFilled bool)

func drawBoard(bottomLeftX, bottomLeftY, cellsX, cellsY int32, fn drawBoardCellCallback) {
	boardSizeX := cellsX * cellSizeX
	boardSizeY := cellsY * cellSizeY

	drawBorderedRectangle(
		bottomLeftX, (bottomLeftY - boardSizeY),
		boardSizeX, boardSizeY,
		boardColor,
		boardOutlineColor,
	)

	for gridY := int32(0); gridY < cellsY; gridY++ {
		for gridX := int32(0); gridX < cellsX; gridX++ {
			screenX := bottomLeftX + (cellSizeX * gridX)
			screenY := bottomLeftY - (cellSizeY * (gridY + 1))

			if color, isFilled := fn(gridX, gridY, screenX, screenY); isFilled {
				drawBorderedRectangle(
					screenX, screenY,
					cellSizeX, cellSizeY,
					color,
					boardColor,
				)
			}
		}
	}
}

func (gs *gameState) Draw() {
	gs.RLock()
	defer gs.RUnlock()

	gs.DrawMainBoard()

	gs.DrawHoldingBoard()

	gs.DrawQueueBoard()

	gs.DrawScore()

	// Draw paused message
	switch gs.Phase {
	case phase_Paused:
		drawCenteredText(pausedText)
	case phase_GameOver:
		drawCenteredText(gameOverText)
	}
}

func (gs *gameState) DrawMainBoard() {
	drawBoard(
		boardBottomLeftX, boardBottomLeftY,
		boardCellsX, boardCellsY_Visible,
		func(gridX, gridY, screenX, screenY int32) (color rl.Color, cellFilled bool) {
			if gs.ActiveTetromino != nil {
				if color, isFilled := gs.ActiveTetromino.IsCell(gridX, gridY); isFilled {
					return color, isFilled
				}
			}

			cell := gs.Board[gridY][gridX]
			if cell.IsFilled {
				return cell.Color, true
			}
			if cell.IsGhost {
				return rl.ColorAlpha(cell.Color, ghostCellAlpha), true
			}

			return rl.Color{}, false
		},
	)
}

func (gs *gameState) DrawHoldingBoard() {
	drawBoard(
		holdingBoardBottomLeftX, holdingBoardBottomLeftY,
		holdingBoardCellsX, holdingBoardCellsY,
		func(gridX, gridY, screenX, screenY int32) (color rl.Color, cellFilled bool) {
			if gs.HoldingTetromino != nil {
				return gs.HoldingTetromino.IsCell(gridX, gridY)
			}

			return rl.Color{}, false
		},
	)
}

func (gs *gameState) DrawQueueBoard() {
	drawBoard(
		queueBoardBottomX, queueBoardBottomY,
		queueBoardCellsX, queueBoardCellsY,
		func(gridX, gridY, screenX, screenY int32) (color rl.Color, cellFilled bool) {
			for _, tetromino := range gs.TetrominoQueue {
				if color, isFilled := tetromino.IsCell(gridX, gridY); isFilled {
					return color, isFilled
				}
			}
			return rl.Color{}, false
		},
	)
}

func (gs *gameState) DrawScore() {
	rl.DrawText(
		scoreText,
		scoreTextX, scoreTextY,
		scoreTextSize,
		textColor,
	)

	rl.DrawText(
		fmt.Sprint(gs.Score),
		scoreNumberX, scoreNumberY,
		scoreTextSize,
		textColor,
	)

	rl.DrawText(
		levelText,
		levelTextX, levelTextY,
		scoreTextSize,
		textColor,
	)

	rl.DrawText(
		fmt.Sprint(gs.Level()),
		levelNumberX, levelNumberY,
		scoreTextSize,
		textColor,
	)
}
