package minesweeper

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func SampleRenderedGame() Minesweeper {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()
	return minesweeper
}

func TestGameActualBombLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	game := minesweeper.(*game)
	properties := minesweeper.(RuntimeGameProperties)

	bombPlacements := make([]Position, int(float32(game.Height*game.Width)*game.difficultyMultiplier))

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				bombPlacements = append(bombPlacements, Position{x, y})
			}
		}
	}

	assert.NotEmpty(t, bombPlacements)

	actualBombLocations := properties.BombLocations()
	for i, bomb := range bombPlacements {
		assert.Equal(t, bomb, actualBombLocations[i])
	}
}

func TestGameActualHintLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	game := minesweeper.(*game)
	properties := minesweeper.(RuntimeGameProperties)

	hintPlacements := make([]Position, 0)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				hintPlacements = append(hintPlacements, Position{x, y})
			}
		}
	}

	assert.NotEmpty(t, hintPlacements)

	actualHintPlacements := properties.HintLocations()
	for i, bomb := range hintPlacements {
		assert.Equal(t, bomb, actualHintPlacements[i])
	}
}

func TestBothBombsAndHintsDoNotShareSameLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	properties := minesweeper.(RuntimeGameProperties)

	hintPlacements := properties.HintLocations()
	bombPlacements := properties.BombLocations()

	assert.NotEmpty(t, hintPlacements)
	assert.NotEmpty(t, bombPlacements)
	for _, hint := range hintPlacements {
		for _, bomb := range bombPlacements {
			if hint.X == bomb.X && hint.Y == bomb.Y {
				assert.Fail(t, fmt.Sprintf("A hint at %v:%v shares the same location with a bomb at %v:%v", hint.X, hint.Y, bomb.X, bomb.Y))
			}
		}
	}
}