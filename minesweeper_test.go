package minesweeper

import (
"testing"
"github.com/stretchr/testify/assert"
)

const (
	SAMPLE_GRID_WIDTH = 50
	SAMPLE_GRID_HEIGHT = 100
)

func TestGridMustNotBeSquaredForTheSakeOfTesting(t *testing.T) {
	assert.True(t, SAMPLE_GRID_WIDTH != SAMPLE_GRID_HEIGHT)
}

func TestBlock_SetBlock(t *testing.T) {
	block := new(Block)

	block.SetBlock(UNKNOWN)
	assert.Equal(t, block.Node, UNKNOWN)
}

func TestGame_SetGrid(t *testing.T) {
	minesweeper := NewGame()
	minesweeper.SetGrid(SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT)
	assert.Equal(t, minesweeper.(*game).Board.Grid, Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}