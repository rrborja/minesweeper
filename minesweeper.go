package minesweeper

import (
	"crypto/rand"
	"encoding/binary"
)

type Node uint8
type Blocks [][]Block
type Grid struct {width, height int}

type Difficulty uint8

const (
	UNKNOWN Node = iota
	BOMB
	NUMBER
	FLAGGED
)

const (
	NOTSET Difficulty = iota
	EASY
	MEDIUM
	HARD
)

const CONSECUTIVE_RANDOM_LIMIT = 3
const EASY_MULTIPLIER = 0.2
const MEDIUM_MULTIPLIER = 0.4
const HARD_MULTIPLIER = 0.6

type Block struct {
	Node
}

type Board struct {
	*Grid
	Blocks
	difficultyMultiplier float32
}

type game struct {
	Board
	Difficulty
}

type Minesweeper interface {
	SetGrid(int, int) error

	SetDifficulty(Difficulty)

	Play() error

	Flag(int, int) error
}

func NewGame(grid ...Grid) Minesweeper {
	game := new(game)
	if len(grid) > 0 {
		game.SetGrid(grid[0].width, grid[0].height)
	}
	return game
}

func (game *game) SetGrid(width, height int) error {
	if game.Grid != nil {
		return new(GameAlreadyStarted)
	}
	game.Grid = &Grid{width, height}
	createBoard(game)
	return nil
}

func (game *game) Flag(x, y int) error {
	game.Blocks[x][y].Node = FLAGGED
	return nil
}

func (game *game) SetDifficulty(difficulty Difficulty) {
	game.Difficulty = difficulty
	switch difficulty {
	case EASY:
		game.difficultyMultiplier = EASY_MULTIPLIER
	case MEDIUM:
		game.difficultyMultiplier = MEDIUM_MULTIPLIER
	case HARD:
		game.difficultyMultiplier = HARD_MULTIPLIER
	}
}

func (game *game) Play() error {
	createBombs(game)
	return nil
}

func (block *Block) SetBlock(node Node) {
	block.Node = node
}

// Shifts to the right
func shiftPosition(grid *Grid, x, y int) (_x, _y int){
	width := grid.width
	height := grid.height
	if x + 1 >= width {
		if y + 1 >= height {
			_x, _y = 0, 0
		} else {
			_x, _y = 0, y + 1
		}
	} else {
		_x, _y = x + 1, y
	}
	return
}

func createBombs(game *game) {
	area := int(game.width * game.height)
	for i := 0; i < int(float32(area) * game.difficultyMultiplier); i++ {
		randomPos := randomNumber(area)

		x, y := randomPos%game.width, randomPos/game.width

		countLimit := 0
		for game.Board.Blocks[x][y].Node != UNKNOWN {
			x, y = shiftPosition(game.Grid, x, y)
			countLimit ++
		}
		if countLimit >= CONSECUTIVE_RANDOM_LIMIT {
			i --
			continue
		}

		game.Blocks[x][y].Node = BOMB
	}
}

func createBoard(game *game) {
	game.Blocks = make([][]Block, game.width)
	for x := range game.Blocks {
		game.Blocks[x] = make([]Block, game.height)
	}
}

func randomNumber(max int) int {
	var n uint16
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return int(n) % max
}