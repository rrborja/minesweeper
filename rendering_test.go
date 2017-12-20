package minesweeper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGame_BombLocations(t *testing.T) {
	minesweeper := newSampleGame()
	game := minesweeper.(*game)
	properties := minesweeper.(RuntimeGameProperties)

	bombPlacements := make([]Position, int(float32(game.Height * game.Width) * game.difficultyMultiplier))

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				bombPlacements = append(bombPlacements, Position{x, y})
			}
		}
	}

	assert.Equal(t, bombPlacements, properties.BombLocations())
}