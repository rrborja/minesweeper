package minesweeper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStateInterface(t *testing.T) {
	var state State = new(Block)

	state.SetBlock(BOMB)

	assert.Equal(t, state.(*Block).Node, BOMB)
}