package minesweeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateInterface(t *testing.T) {
	var state State = new(Block)

	state.SetBlock(BOMB)

	assert.Equal(t, state.(*Block).Node, BOMB)
}
