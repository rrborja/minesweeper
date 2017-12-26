package minesweeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionNew(t *testing.T) {
	assert.Empty(t, singleton)
	event := New()
	assert.Equal(t, event, singleton.(*game).Event)
}

func TestFunctionSetGrid(t *testing.T) {
	New()
	SetGrid(3, 4)
	assert.EqualValues(t, &Grid{3, 4}, singleton.(*game).Grid)
}

func TestFunctionSetDifficulty(t *testing.T) {
	New()
	SetDifficulty(Hard)
	assert.EqualValues(t, Hard, singleton.(*game).Difficulty)
}

func TestFunctionPlay(t *testing.T) {
	New()
	SetGrid(4, 5)
	SetDifficulty(Easy)
	Play()
	assert.NotEmpty(t, singleton.(*game).Mutex)
}

func TestFunctionFlag(t *testing.T) {
	New()
	SetGrid(4, 5)
	SetDifficulty(Medium)
	Play()
	Flag(0, 0)
	assert.True(t, singleton.(*game).blocks[0][0].flagged)
}

func TestFunctionVisit(t *testing.T) {
	New()
	SetGrid(3, 8)
	SetDifficulty(Easy)
	Play()
	Visit(2, 2)
	assert.True(t, singleton.(*game).blocks[2][2].visited)
}
