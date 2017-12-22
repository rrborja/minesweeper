/*
 * Minesweeper API
 * Copyright (C) 2017  Ritchie Borja
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package minesweeper

import (
	"fmt"
	"testing"

	"github.com/rrborja/minesweeper-go/rendering"
	"github.com/rrborja/minesweeper-go/visited"
	"github.com/stretchr/testify/assert"
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
	properties := minesweeper.(rendering.Locations)

	bombPlacements := make([]rendering.Position, int(float32(game.Height*game.Width)*game.difficultyMultiplier))

	var counter int
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == BOMB {
				bombPlacements[counter] = block
				counter++
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
	properties := minesweeper.(rendering.Locations)

	hintPlacements := make([]rendering.Position, 0)

	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == NUMBER {
				hintPlacements = append(hintPlacements, block)
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
	properties := minesweeper.(rendering.Locations)

	hintPlacements := properties.HintLocations()
	bombPlacements := properties.BombLocations()

	assert.NotEmpty(t, hintPlacements)
	assert.NotEmpty(t, bombPlacements)
	for _, hint := range hintPlacements {
		for _, bomb := range bombPlacements {
			if hint.X() == bomb.X() && hint.Y() == bomb.Y() {
				assert.Fail(t, fmt.Sprintf("A hint at %v:%v shares the same location with a bomb at %v:%v", hint.X(), hint.Y(), bomb.X(), bomb.Y()))
			}
		}
	}
}

func TestRecentPlayersMove(t *testing.T) {
	minesweeper, _ := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	minesweeper.SetDifficulty(MEDIUM)
	minesweeper.Play()

	var story visited.Story = minesweeper.(*game)

	var recentMove visited.Record

	maxMoves := 10
	for i := 0; i < maxMoves; i++ {
		randomX := randomNumber(SAMPLE_GRID_WIDTH)
		randomY := randomNumber(SAMPLE_GRID_HEIGHT)
		blocks, err := minesweeper.Visit(randomX, randomY)

		if len(blocks) == 0 { // Either already visited block or flagged block
			continue
		}

		var expectedAction visited.Action
		switch blocks[0].Node {
		case UNKNOWN:
			expectedAction = visited.Unknown
		case NUMBER:
			expectedAction = visited.Number
		case BOMB:
			expectedAction = visited.Bomb
		}

		recentMove = visited.Record{Position: blocks[0], Action: expectedAction}

		if err != nil {
			break
		}
	}

	assert.Equal(t, recentMove, story.LastAction())
}
