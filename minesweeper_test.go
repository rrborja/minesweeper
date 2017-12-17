package minesweeper

import (
"testing"
"github.com/stretchr/testify/assert"
	"fmt"
	"math/rand"
	"time"
)

const (
	SAMPLE_GRID_WIDTH = 50
	SAMPLE_GRID_HEIGHT = 100
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

func TestBlock_SetBlock(t *testing.T) {
	block := new(Block)

	block.SetBlock(UNKNOWN)
	assert.Equal(t, block.Node, UNKNOWN)
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
	assert.NotNil(t, err, "Must report an error upon setting a new grid from an already started game")
	assert.IsType(t, new(GameAlreadyStarted), err, "The error must be GameAlreadyStarted error type")
}

func TestFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.Flag(3, 6)
	assert.Equal(t, minesweeper.(*game).Blocks[3][6].Node, FLAGGED)
}

func TestGame_SetDifficulty(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	assert.Equal(t, minesweeper.(*game).Difficulty, EASY)
}

func TestBombsInPlace(t *testing.T) {

	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	rand.Seed(time.Now().Unix())

	numOfBombs := int(float32(game.width * game.height) * EASY_MULTIPLIER)
	countedBombs := 0
	for _, row := range game.Blocks {
		fmt.Println()
		for _, block := range row {
			fmt.Print(block.Node)
			if block.Node == BOMB {
				countedBombs ++
			}
		}
	}
	assert.Equal(t, numOfBombs, countedBombs)
}