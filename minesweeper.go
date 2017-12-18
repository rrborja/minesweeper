package minesweeper

import (
	"crypto/rand"
	"encoding/binary"
)

type Node uint8
type Blocks [][]Block
type Grid struct{ width, height int }

type Difficulty uint8

const (
	UNKNOWN Node = 1 << iota >> 1
	BOMB
	NUMBER
)

const (
	NOTSET Difficulty = iota
	EASY
	MEDIUM
	HARD
)

const CONSECUTIVE_RANDOM_LIMIT = 3
const EASY_MULTIPLIER = 0.1
const MEDIUM_MULTIPLIER = 0.2
const HARD_MULTIPLIER = 0.5

type Block struct {
	Node
	value            int
	visited, flagged bool
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

	Flag(int, int)

	Visit(int, int) error
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

func (game *game) Flag(x, y int) {
	game.Blocks[x][y].flagged = true
}

func (game *game) Visit(x, y int) error {
	game.Blocks[x][y].visited = true
	switch game.Blocks[x][y].Node {
	case BOMB:
		return new(Exploded)
	case UNKNOWN:
		game.Blocks[x][y].visited = false //to avoid infinite recursion, first is to set the base case
		autoRevealUnmarkedBlock(game, x, y)
	}
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
	tallyHints(game)
	return nil
}

func (block *Block) SetBlock(node Node) {
	block.Node = node
}

// Shifts to the right
func shiftPosition(grid *Grid, x, y int) (_x, _y int) {
	width := grid.width
	height := grid.height
	if x+1 >= width {
		if y+1 >= height {
			_x, _y = 0, 0
		} else {
			_x, _y = 0, y+1
		}
	} else {
		_x, _y = x+1, y
	}
	return
}

func createBombs(game *game) {
	area := int(game.width * game.height)
	for i := 0; i < int(float32(area)*game.difficultyMultiplier); i++ {
		randomPos := randomNumber(area)

		x, y := randomPos%game.width, randomPos/game.width

		countLimit := 0
		for game.Board.Blocks[x][y].Node != UNKNOWN {
			x, y = shiftPosition(game.Grid, x, y)
			countLimit++
		}
		if countLimit >= CONSECUTIVE_RANDOM_LIMIT {
			i--
			continue
		}

		game.Blocks[x][y].Node = BOMB
	}
}

func tallyHints(game *game) {
	width := game.width
	height := game.height

	tally := func(blocks Blocks, x, y int) {
		if x >= 0 && y >= 0 &&
			x < width && y < height &&
			blocks[x][y].Node != BOMB {
			blocks[x][y].Node = NUMBER
			blocks[x][y].value++
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				tally(game.Blocks, x-1, y-1)
				tally(game.Blocks, x-1, y)
				tally(game.Blocks, x-1, y+1)
				tally(game.Blocks, x, y-1)
				tally(game.Blocks, x, y+1)
				tally(game.Blocks, x+1, y-1)
				tally(game.Blocks, x+1, y)
				tally(game.Blocks, x+1, y+1)
			}
		}
	}
}

func createBoard(game *game) {
	game.Blocks = make([][]Block, game.width)
	for x := range game.Blocks {
		game.Blocks[x] = make([]Block, game.height)
	}
}

func autoRevealUnmarkedBlock(game *game, x, y int) {
	blocks := game.Blocks
	width := game.width
	height := game.height

	if x >= 0 && y >= 0 && x < width && y < height {
		if blocks[x][y].visited {
			return
		}
		if blocks[x][y].Node == UNKNOWN {
			blocks[x][y].visited = true

			autoRevealUnmarkedBlock(game, x-1, y-1)
			autoRevealUnmarkedBlock(game, x-1, y)
			autoRevealUnmarkedBlock(game, x-1, y+1)
			autoRevealUnmarkedBlock(game, x, y-1)
			autoRevealUnmarkedBlock(game, x, y+1)
			autoRevealUnmarkedBlock(game, x+1, y-1)
			autoRevealUnmarkedBlock(game, x+1, y)
			autoRevealUnmarkedBlock(game, x+1, y+1)
		} else if blocks[x][y].Node == NUMBER {
			blocks[x][y].visited = true
		}
	}
}

func randomNumber(max int) int {
	var number uint16
	binary.Read(rand.Reader, binary.LittleEndian, &number)
	return int(number) % max
}
