package main

import (
	"math"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type phase int

const (
	// Generation: No active tetromino.
	//             Wait 0.2 seconds before spawning a new one.
	phase_Generation phase = iota
	// Falling: Active tetromino is falling until it hits something.
	phase_Falling

	// Lock: Tetrimino has hit something.
	phase_Lock

	// Completion: Active tetromino has landed.
	//             Check for full rows before moving to generation
	phase_Completion

	// Paused: Game is paused. Can only be entered from Falling.
	phase_Paused

	// GameOver: Game is over. Goes to end after any input.
	phase_GameOver

	// End: Exit immediately.
	phase_End
)

type cell struct {
	IsFilled bool
	IsGhost  bool
	Color    rl.Color
}

type board [boardCellsY][boardCellsX]cell

type gameState struct {
	sync.RWMutex
	ActiveTetromino  *tetromino
	HoldingTetromino *tetromino
	TetrominoQueue   [tetrominoQueueSize]tetromino
	Board            board
	Phase            phase

	linesCleared int
	Score        int

	IsDone bool
}

func newGameState() *gameState {
	gs := &gameState{}
	gs.Phase = phase_Generation
	gs.IsDone = false

	// Initialize the board
	for i := int32(0); i < boardCellsY; i++ {
		for j := int32(0); j < boardCellsX; j++ {
			gs.Board[i][j] = cell{
				IsFilled: false,
				Color:    rl.Black,
			}
		}
	}

	// Initialize the tetromino queue
	for i := 0; i < len(gs.TetrominoQueue); i++ {
		gs.TetrominoQueue[i] = *NewRandomTetromino(
			tetrominoQueueX,
			tetrominoQueueY+int32(i*4),
		)
	}

	return gs
}

func (gs *gameState) Level() int {
	level := int(gs.linesCleared/linesClearedPerLevel) + 1
	return level
}

func (gs *gameState) DropInterval(multiplier float64) time.Duration {
	// Formula taken from Tetris Guide 2009, added multiplier
	level := float64(gs.Level() - 1)
	interval := math.Pow((0.8-(level*0.007)), level) * multiplier
	return time.Duration(interval * float64(time.Second))
}

type withLockFunc func() bool

// WithLock executes the given function with a lock on the game state.
func (gs *gameState) WithLock(fn withLockFunc) bool {
	gs.Lock()
	ret := fn()
	gs.Unlock()

	return ret
}

//// Tetromino Actions

func (gs *gameState) ActiveTetrominoDown() (didCollide bool) {
	return gs.WithLock(func() bool {
		gs.ActiveTetromino.OriginY -= 1

		didCollide := gs.ActiveTetromino.CheckCollision(&gs.Board)
		if didCollide {
			gs.ActiveTetromino.OriginY += 1
		}

		return didCollide
	})
}

func (gs *gameState) ActiveTetrominoLeft() (didCollide bool) {
	return gs.WithLock(func() bool {
		gs.ActiveTetromino.OriginX -= 1

		didCollide := gs.ActiveTetromino.CheckCollision(&gs.Board)
		if didCollide {
			gs.ActiveTetromino.OriginX += 1
		}

		return didCollide
	})
}

func (gs *gameState) ActiveTetrominoRight() (didCollide bool) {
	return gs.WithLock(func() bool {
		gs.ActiveTetromino.OriginX += 1

		didCollide := gs.ActiveTetromino.CheckCollision(&gs.Board)
		if didCollide {
			gs.ActiveTetromino.OriginX -= 1
		}

		return didCollide
	})
}

func (gs *gameState) ActiveTetrominoRotateClockwise() (didCollide bool) {
	return gs.WithLock(func() bool {
		gs.ActiveTetromino.RotateClockwise()

		didCollide := gs.ActiveTetromino.CheckCollision(&gs.Board)
		if didCollide {
			gs.ActiveTetromino.RotateCounterClockwise()
		}

		return didCollide
	})
}

func (gs *gameState) ActiveTetrominoRotateCounterClockwise() (didCollide bool) {
	return gs.WithLock(func() bool {
		gs.ActiveTetromino.RotateCounterClockwise()

		didCollide := gs.ActiveTetromino.CheckCollision(&gs.Board)
		if didCollide {
			gs.ActiveTetromino.RotateClockwise()
		}

		return didCollide
	})
}

func (gs *gameState) ActiveTetrominoHold() (shouldGenerate bool) {
	return gs.WithLock(func() bool {
		gs.HoldingTetromino, gs.ActiveTetromino = gs.ActiveTetromino, gs.HoldingTetromino
		gs.HoldingTetromino.OriginX = tetrominoHoldingX
		gs.HoldingTetromino.OriginY = tetrominoHoldingY

		if gs.ActiveTetromino != nil {
			gs.ActiveTetromino.OriginX = tetrominoGenerateX
			gs.ActiveTetromino.OriginY = tetrominoGenerateY
			return false
		}

		return true
	})
}

func (gs *gameState) ActiveTetrominoHardDown() {
	gs.WithLock(func() bool {
		for {
			gs.ActiveTetromino.OriginY -= 1

			if gs.ActiveTetromino.CheckCollision(&gs.Board) {
				gs.ActiveTetromino.OriginY += 1
				break
			}

			gs.ActiveTetromino.CommitTrailToBoard(&gs.Board)
		}

		return false
	})

	time.AfterFunc(trailClearDelay, func() {
		gs.WithLock(func() bool {
			for i := int32(0); i < boardCellsY; i++ {
				for j := int32(0); j < boardCellsX; j++ {
					gs.Board[i][j].IsGhost = false
				}
			}
			return false
		})
	})
}

//// Phases

func (gs *gameState) GenerationPhase() {
	// Spawn a new tetromino
	time.Sleep(generationDelay)
	gameOver := gs.WithLock(func() bool {
		topmino := gs.TetrominoQueue[tetrominoQueueSize-1]
		gs.ActiveTetromino = &topmino
		gs.ActiveTetromino.OriginX = tetrominoGenerateX
		gs.ActiveTetromino.OriginY = tetrominoGenerateY

		// Drop the tetromino once to check for collisions
		return gs.ActiveTetromino.CheckCollision(&gs.Board)
	})

	// Move all tetrominos in the queue up
	for i := tetrominoQueueSize - 1; i > 0; i-- {
		gs.TetrominoQueue[i] = gs.TetrominoQueue[i-1]
		gs.TetrominoQueue[i].OriginY += 4
	}

	// Generate a new tetromino
	gs.TetrominoQueue[0] = *NewRandomTetromino(
		tetrominoQueueX,
		tetrominoQueueY,
	)

	if gameOver {
		gs.Phase = phase_GameOver
		return
	}
	gs.Phase = phase_Falling
}

func (gs *gameState) FallingPhase(inputEvents chan InputEvent) {
	dropTicker := time.NewTicker(gs.DropInterval(1.0))
	done := false

	stop := func(nextPhase phase) {
		dropTicker.Stop()
		done = true
		gs.Phase = nextPhase
	}

	adjustTicker := func(multiplier float64) {
		newInterval := gs.DropInterval(multiplier)
		dropTicker.Reset(newInterval)
	}

	for !done {
		select {
		case <-dropTicker.C:
			if didCollide := gs.ActiveTetrominoDown(); didCollide {
				stop(phase_Lock)
			}
		case event := <-inputEvents:
			switch event.Input {
			case Input_Pause:
				if event.Action == Action_Up {
					stop(phase_Paused)
				}
			case Input_Hold:
				if event.Action == Action_Down {
					if shouldGenerate := gs.ActiveTetrominoHold(); shouldGenerate {
						stop(phase_Generation)
					}
				}
			case Input_MoveLeft:
				if event.Action == Action_Down {
					gs.ActiveTetrominoLeft()
				}
			case Input_MoveRight:
				if event.Action == Action_Down {
					gs.ActiveTetrominoRight()
				}
			case Input_RotateClockwise:
				if event.Action == Action_Down {
					gs.ActiveTetrominoRotateClockwise()
				}
			case Input_RotateCounterClockwise:
				if event.Action == Action_Down {
					gs.ActiveTetrominoRotateCounterClockwise()
				}
			case Input_SoftDrop:
				switch event.Action {
				case Action_Down:
					adjustTicker(softDropMultiplier)
				case Action_Up:
					adjustTicker(1)
				}
			case Input_HardDrop:
				if event.Action == Action_Down {
					gs.ActiveTetrominoHardDown()
					stop(phase_Lock)
				}
			}
		}
	}
}

func (gs *gameState) LockPhase() {
	// Done falling, commit active tetromino to board
	gs.WithLock(func() bool {
		gs.ActiveTetromino.CommitToBoard(&gs.Board)
		gs.ActiveTetromino = nil

		return true
	})
	gs.Phase = phase_Completion
}

func (gs *gameState) CompletionPhase() {
	rowsToDelete := []int32{}

	// Mark rows for deletion
	shouldDeleteRows := gs.WithLock(func() bool {
		for i := boardCellsY - 1; i >= 0; i-- {
			row := gs.Board[i]
			isRowComplete := true
			for _, cell := range row {
				if !cell.IsFilled {
					isRowComplete = false
					break
				}
			}

			if isRowComplete {
				rowsToDelete = append(rowsToDelete, i)
				// Mark cells visually as deleted
				for j := int32(0); j < boardCellsX; j++ {
					gs.Board[i][j].IsFilled = false
					gs.Board[i][j].IsGhost = true
				}
			}
		}

		return len(rowsToDelete) > 0
	})

	if !shouldDeleteRows {
		gs.Phase = phase_Generation
		return
	}

	linesCleared := len(rowsToDelete)
	gs.linesCleared += linesCleared
	switch linesCleared {
	case 1:
		gs.Score += 100 * gs.Level()
	case 2:
		gs.Score += 300 * gs.Level()
	case 3:
		gs.Score += 500 * gs.Level()
	case 4:
		gs.Score += 800 * gs.Level()
	default:
		gs.Score += 1200 * gs.Level()
	}

	// Small delay so the user can see the rows being deleted
	time.Sleep(rowClearDelay)

	// Delete marked rows
	gs.WithLock(func() bool {
		// Start with the topmost row to delete
		for i := 0; i < len(rowsToDelete); i++ {
			// starting from current row, move all rows above down
			for j := int32(rowsToDelete[i]); j < boardCellsY-1; j++ {
				gs.Board[j] = gs.Board[j+1]
			}

			// clear the top row
			gs.Board[boardCellsY-1] = [boardCellsX]cell{}
		}

		return true
	})

	gs.Phase = phase_Generation
}

func (gs *gameState) PausedPhase(inputEvents chan InputEvent) {
	for {
		event := <-inputEvents
		if event.Input == Input_Pause && event.Action == Action_Up {
			gs.Phase = phase_Falling
			return
		}
	}
}

func (gs *gameState) GameOverPhase(inputEvents chan InputEvent) {
	<-inputEvents
	gs.Phase = phase_End
}

//// Main loop

func (gs *gameState) Run(inputEvents chan InputEvent) {
	go func() {
		for {
			// Each phase will run until it's ready to move to another phase
			switch gs.Phase {
			case phase_Generation:
				gs.GenerationPhase()
			case phase_Falling:
				gs.FallingPhase(inputEvents)
			case phase_Lock:
				gs.LockPhase()
			case phase_Completion:
				gs.CompletionPhase()
			case phase_Paused:
				gs.PausedPhase(inputEvents)
			case phase_GameOver:
				gs.GameOverPhase(inputEvents)
			case phase_End:
				gs.IsDone = true
				return
			}
		}
	}()
}
