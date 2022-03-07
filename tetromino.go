package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type tetromino struct {
	// Origin is in gameboard space, not screen space
	OriginX, OriginY int32
	// Cells are offsets from the origin
	cells [4][2]int32
	color rl.Color
}

type cellIteratorFunction func(x, y int32) bool

// cellIterator call the iterator function for each cell in the tetromino
// The iterator function is given the absolute coordinates of the cell
// If the iter function returns true, iteration is stopped and this function returns true
// If no functions return true, this function returns false
func (t *tetromino) cellIterator(f cellIteratorFunction) bool {
	for _, cell := range t.cells {
		if ok := f(cell[0]+t.OriginX, cell[1]+t.OriginY); ok {
			return true
		}
	}

	return false
}

// IsCell given coordinates in gameboard space, returns true if one of the tetrominos blocks is in those coordinates
func (t *tetromino) IsCell(x, y int32) (rl.Color, bool) {
	return t.color, t.cellIterator(func(cx, cy int32) bool {
		return cx == x && cy == y
	})
}

// CheckCollision returns true if the tetromino is colliding with the gameboard
func (t *tetromino) CheckCollision(b *board) bool {
	return t.cellIterator(func(x, y int32) bool {
		// Check if the cell is outside the gameboard
		if x < 0 || x >= boardCellsX || y < 0 || y >= boardCellsY {
			return true
		}

		// Check if the cell is already filled
		return b[y][x].IsFilled
	})
}

// CommitToBoard copies each cell of tetromino to the gameboard
func (t *tetromino) CommitToBoard(b *board) {
	t.cellIterator(func(x, y int32) bool {
		b[y][x].IsFilled = true
		b[y][x].Color = t.color
		return false
	})
}

// CommitTrailToBoard copies each cell of the tetromino as a trail to the gameboard
func (t *tetromino) CommitTrailToBoard(b *board) {
	t.cellIterator(func(x, y int32) bool {
		b[y][x].IsGhost = true
		b[y][x].Color = t.color
		return false
	})
}

func (t *tetromino) RotateClockwise() {
	// To rotate counter clockwise,
	// first swap the x and y components,
	// then invert the y component
	for i := 0; i < 4; i++ {
		t.cells[i][0], t.cells[i][1] = t.cells[i][1], -t.cells[i][0]
	}
}

func (t *tetromino) RotateCounterClockwise() {
	// To rotate counter clockwise,
	// do the opposite of the clockwise rotation
	for i := 0; i < 4; i++ {
		t.cells[i][0], t.cells[i][1] = -t.cells[i][1], t.cells[i][0]
	}
}

func NewRandomTetromino(originX, originY int32) *tetromino {
	// Note cells are defined in clockwise order
	// Todo: Come up with a better pallete, instead of using builtin colors
	switch rl.GetRandomValue(0, 6) {
	case 0:
		// O-Tetromino
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   oTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{0, 1},
				{1, 1},
				{1, 0},
			},
		}
	case 1:
		// I-Tetromino
		// Note: In the real tetris, the I-tetromino's origin is
		//       in the center of 4 cells, not the cengter-left block.
		//       It's an edge case we won't care about for now.
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   iTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{-1, 0},
				{1, 0},
				{2, 0},
			},
		}
	case 2:
		// T-Tetromino
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   tTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{-1, 0},
				{0, 1},
				{1, 0},
			},
		}
	case 3:
		// L-Tetromino
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   lTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{-1, 0},
				{1, 1},
				{1, 0},
			},
		}
	case 4:
		// J-Tetromino
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   jTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{-1, 0},
				{-1, -1},
				{1, 0},
			},
		}
	case 5:
		// S-Tetromino
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   sTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{-1, 0},
				{0, 1},
				{1, 1},
			},
		}
	case 6:
		// Z-Tetromino
		return &tetromino{
			OriginX: originX,
			OriginY: originY,
			color:   zTetriminoColor,
			cells: [4][2]int32{
				{0, 0},
				{-1, 1},
				{0, 1},
				{1, 0},
			},
		}
	}

	log.Fatal("NewRandomTetromino: Invalid random value")
	return nil // Should never happen
}
