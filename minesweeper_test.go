package minesweeper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	SAMPLE_GRID_WIDTH  = 10
	SAMPLE_GRID_HEIGHT = 40
)

func newBlankGame() Minesweeper {
	return NewGame()
}

func newSampleGame() Minesweeper {
	return NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestGridMustNotBeSquaredForTheSakeOfTesting(t *testing.T) {
	assert.True(t, SAMPLE_GRID_WIDTH != SAMPLE_GRID_HEIGHT)
}

func TestGame_SetGrid(t *testing.T) {
	minesweeper := newBlankGame()
	minesweeper.SetGrid(SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT)
	assert.Equal(t, minesweeper.(*game).Board.Grid, &Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestGameWithGridArgument(t *testing.T) {
	minesweeper := newSampleGame()
	assert.Equal(t, minesweeper.(*game).Board.Grid, &Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestNewGridWhenStartedGame(t *testing.T) {
	minesweeper := newSampleGame()
	err := minesweeper.SetGrid(10, 20)
	assert.Error(t, err)
	assert.NotNil(t, err, "Must report an error upon setting a new grid from an already started game")
	assert.IsType(t, new(GameAlreadyStarted), err, "The error must be GameAlreadyStarted error type")
}

func TestFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.Flag(3, 6)
	assert.Equal(t, minesweeper.(*game).Blocks[3][6].flagged, true)
}

func TestGame_SetDifficulty(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	assert.Equal(t, minesweeper.(*game).Difficulty, EASY)
}

func TestShiftFromMaxPosition(t *testing.T) {
	grid := Grid{5, 5}
	x, y := shiftPosition(&grid, 4, 4)
	assert.Equal(t, struct {
		x int
		y int
	}{0, 0}, struct {
		x int
		y int
	}{x, y})
}

func TestBombsInPlace(t *testing.T) {

	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	numOfBombs := int(float32(game.Width*game.Height) * EASY_MULTIPLIER)
	countedBombs := 0
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == BOMB {
				countedBombs++
			}
		}
	}
	assert.Equal(t, numOfBombs, countedBombs)
}

func TestTalliedBomb(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	width := game.Width
	height := game.Height

	count := func(blocks Blocks, x, y int) (has int) {
		if x >= 0 && y >= 0 &&
			x < width && y < height &&
			blocks[x][y].Node&BOMB == 1 {
			return 1
		}
		return
	}

	hasSurroundingTally := func(blocks Blocks, x, y int) int {
		if x >= 0 && y >= 0 &&
			x < width && y < height {
			switch blocks[x][y].Node {
			case NUMBER:
				return 1
			case BOMB:
				return -1
			default:
				return 0
			}
		}
		return -1
	}
	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y+1))
			}
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				var counted int
				counted = count(game.Blocks, x-1, y-1) +
					count(game.Blocks, x-1, y) +
					count(game.Blocks, x-1, y+1) +
					count(game.Blocks, x, y-1) +
					count(game.Blocks, x, y+1) +
					count(game.Blocks, x+1, y-1) +
					count(game.Blocks, x+1, y) +
					count(game.Blocks, x+1, y+1)
				assert.Equal(t, counted, block.value)
			}
		}
	}
}

func TestVisitedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				game.Visit(x, y)
				assert.True(t, game.Blocks[x][y].visited)
			}
		}
	}

}

func TestVisitedBombToGameOver(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	var x, y int
	var err error

	for i, row := range game.Blocks {
		for j, block := range row {
			if block.Node == BOMB {
				x, y = i, j
				_, err = game.Visit(x, y)
				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.IsType(t, new(Exploded), err)
			}
		}
	}

}

func TestVisitedBombToGameOverWithCorrectLocationReason(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	var x, y int
	var err error

	for i, row := range game.Blocks {
		for j, block := range row {
			if block.Node == BOMB {
				x, y = i, j
				_, err = game.Visit(x, y)
				assert.Error(t, err)
				assert.EqualError(t, err,
					fmt.Sprintf("Game over at X=%v Y=%v",
						x, y))
				assert.IsType(t, new(Exploded), err)
			}
		}
	}

}

func TestVisitedUnmarkedBlockDistributeVisit(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == UNKNOWN && !block.visited {
				minesweeper.Visit(x, y)
			}
		}
	}

	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == UNKNOWN {
				assert.True(t, block.visited)
			}
		}
	}
}

func TestVisitAFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				minesweeper.Flag(x, y)
				_, err := minesweeper.Visit(x, y)
				assert.NoError(t, err)
				if err != nil {
					assert.IsType(t, new(Exploded), err)
				}
			}
		}
	}
}

func print(game *game) {
	for _, row := range game.Blocks {
		fmt.Println()
		for _, block := range row {
			if block.Node == BOMB {
				fmt.Print("* ")
			} else if block.Node == UNKNOWN {
				fmt.Print("  ")
			} else {
				fmt.Printf("%v ", block.value)
			}
		}
	}
}
