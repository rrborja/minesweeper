package minesweeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExploded_Error(t *testing.T) {
	err := Exploded{struct{ x, y int }{x: 3, y: 9}}
	assert.EqualError(t, err, "Game over at X=3 Y=9")
}

func TestGameAlreadyStarted_Error(t *testing.T) {
	err := GameAlreadyStarted{}
	assert.EqualError(t, err, "Game already started. Try setting a new board.")
}

func TestUnspecifiedDifficulty_Error(t *testing.T) {
	err := UnspecifiedDifficulty{}
	assert.EqualError(t, err, "Difficulty was not specified. Use Difficulty(Difficulty) method before calling Play()")
}
