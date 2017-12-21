package minesweeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExploded_Error(t *testing.T) {
	err := Exploded{struct{ x, y int }{x: 3, y: 9}}
	assert.Equal(t, "Game over at X=3 Y=9", err.Error())
}

func TestGameAlreadyStarted_Error(t *testing.T) {
	err := GameAlreadyStarted{}
	assert.Equal(t, "Game already started. Try setting a new board.", err.Error())
}
