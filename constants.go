package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	internalScreenX int32 = 800
	internalScreenY int32 = 450

	cellSizeX int32 = 20
	cellSizeY int32 = 20

	// Main board
	boardCellsX         int32 = 10
	boardCellsY         int32 = 40
	boardCellsY_Visible int32 = 20 // Only the bottom 20 lines are visible
	boardSizeX          int32 = cellSizeX * boardCellsX
	boardSizeY          int32 = cellSizeY * boardCellsY_Visible
	boardBottomLeftX    int32 = -(boardSizeX / 2)
	boardBottomLeftY    int32 = (boardSizeY / 2)

	// Pos in main board
	tetrominoGenerateX int32 = 5
	tetrominoGenerateY int32 = 20

	// Holding board
	holdingBoardCellsX      int32 = 5
	holdingBoardCellsY      int32 = 5
	holdingBoardSizeX       int32 = cellSizeX * holdingBoardCellsX
	holdingBoardSizeY       int32 = cellSizeY * holdingBoardCellsY
	holdingBoardMargin      int32 = cellSizeX
	holdingBoardBottomLeftX int32 = boardBottomLeftX - holdingBoardSizeX - holdingBoardMargin
	holdingBoardBottomLeftY int32 = (boardBottomLeftY - (boardSizeY - holdingBoardSizeY))

	// Pos in holding board
	tetrominoHoldingX int32 = 2
	tetrominoHoldingY int32 = 2

	// Queue board
	tetrominoQueueSize int32 = 5
	queueBoardCellsX   int32 = 5
	queueBoardCellsY   int32 = 4 * tetrominoQueueSize
	queueBoardMargin   int32 = cellSizeX
	queueBoardSizeX    int32 = cellSizeX * queueBoardCellsX
	queueBoardSizeY    int32 = cellSizeY * queueBoardCellsY
	queueBoardBottomX  int32 = boardBottomLeftX + boardSizeX + queueBoardMargin
	queueBoardBottomY  int32 = boardBottomLeftY

	// Pos in queue board
	tetrominoQueueX int32 = 2
	tetrominoQueueY int32 = 1 // Multiplied by pos in queue

	softDropMultiplier float64 = 0.05 // 20 times faster

	generationDelay        time.Duration = time.Millisecond * 200 // 0.2 seconds
	rowClearDelay          time.Duration = time.Millisecond * 75  // 0.075 seconds
	trailClearDelay        time.Duration = time.Millisecond * 50  // 0.05 seconds
	inputLongPressInterval time.Duration = time.Millisecond * 300 // 0.3 seconds
	inputPollingInterval   time.Duration = time.Millisecond * 10  // 0.01 seconds

	// During auto repeat a tetromino should be able to move to the edge in 0.5 seconds
	//autoRepeatInterval time.Duration = time.Millisecond * time.Duration(500/boardCellsX)

	linesClearedPerLevel int = 10

	titleTextSize    float32 = 50.0
	titleTextSpacing float32 = 1.0
	pausedText       string  = "PAUSED"
	gameOverText     string  = "GAME OVER"

	scoreTextSize int32  = 30.0
	levelText     string = "LEVEL"
	scoreText     string = "SCORE"

	scoreLineSpacing int32 = 10
	scoreLineSizeY   int32 = scoreTextSize + scoreLineSpacing
	scoreSizeY       int32 = scoreLineSizeY * 3 // 4 lines

	scoreBottomLeftX int32 = holdingBoardBottomLeftX
	scoreBottomLeftY int32 = holdingBoardBottomLeftY + holdingBoardMargin + scoreSizeY

	levelTextX   int32 = scoreBottomLeftX
	levelTextY   int32 = scoreBottomLeftY - (scoreLineSizeY * 3)
	levelNumberX int32 = scoreBottomLeftX
	levelNumberY int32 = scoreBottomLeftY - (scoreLineSizeY * 2)

	scoreTextX   int32 = scoreBottomLeftX
	scoreTextY   int32 = scoreBottomLeftY - (scoreLineSizeY)
	scoreNumberX int32 = scoreBottomLeftX
	scoreNumberY int32 = scoreBottomLeftY
)

// Pallette has to be var because the rl.Color type can't be a constant
var (
	// Pallette
	ghostCellAlpha    float32  = 0.25
	backgroundColor   rl.Color = rl.GetColor(0x3E363FFF)
	boardColor        rl.Color = rl.GetColor(0x504850FF)
	boardOutlineColor rl.Color = rl.Black
	textColor         rl.Color = rl.RayWhite
	// https://coolors.co/ffa122-fcfc32-00c400-ac17ac-f50000-5193e8-310ca9
	oTetriminoColor rl.Color = rl.GetColor(0xFCFC32FF) // Yellow
	iTetriminoColor rl.Color = rl.GetColor(0x5193E8FF) // Light Blue
	tTetriminoColor rl.Color = rl.GetColor(0xAC17ACFF) // Purple
	lTetriminoColor rl.Color = rl.GetColor(0xFFA122FF) // Orange
	jTetriminoColor rl.Color = rl.GetColor(0x310CA9FF) // Dark Blue
	sTetriminoColor rl.Color = rl.GetColor(0x00C400FF) // Green
	zTetriminoColor rl.Color = rl.GetColor(0xF50000FF) // Red
)
