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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/rrborja/minesweeper/rendering"
	"github.com/rrborja/minesweeper/visited"
	"github.com/stretchr/testify/assert"
)

func SampleRenderedGame() Minesweeper {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()
	return minesweeper
}

func TestGameActualBombLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	game := minesweeper.(*game)
	properties := minesweeper.(rendering.Tracker)

	bombPlacements := make([]rendering.Position, int(float32(game.Height*game.Width)*game.difficultyMultiplier))

	var counter int
	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node == Bomb {
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
	properties := minesweeper.(rendering.Tracker)

	hintPlacements := make([]rendering.Position, 0)

	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node == Number {
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
	properties := minesweeper.(rendering.Tracker)

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
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Medium)
	minesweeper.Play()

	var story visited.StoryTeller = minesweeper.(*game)

	var recentMove visited.Record

	maxMoves := 10
	for i := 0; i < maxMoves; i++ {
		randomX := randomNumber(sampleGridWidth)
		randomY := randomNumber(sampleGridHeight)
		blocks, err := minesweeper.Visit(randomX, randomY)

		if len(blocks) == 0 { // Either already visited block or flagged block
			continue
		}

		var expectedAction visited.Action
		switch blocks[0].Node {
		case Unknown:
			expectedAction = visited.Unknown
		case Number:
			expectedAction = visited.Number
		case Bomb:
			expectedAction = visited.Bomb
		}

		recentMove = visited.Record{Position: blocks[0], Action: expectedAction}

		if err != nil {
			break
		}
	}

	assert.Equal(t, recentMove, story.LastAction())
}

func TestGamePrintBoard(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	properties := minesweeper.(rendering.Tracker)
	bombs := properties.BombLocations()
	hints := properties.HintLocations()

	star := '*'

	var board = make([][]*rune, sampleGridWidth)
	for i := range board {
		board[i] = make([]*rune, sampleGridHeight)
	}

	for _, bomb := range bombs {
		x := bomb.X()
		y := bomb.Y()
		board[x][y] = &star
	}

	for _, hint := range hints {
		x := hint.X()
		y := hint.Y()
		value := rune(hint.(Block).Value + 48)
		board[x][y] = &value
	}

	var boardLayout = make([]string, sampleGridWidth)
	for i, row := range board {
		cellLayout := make([]rune, (sampleGridHeight * 2))
		for j, cell := range row {
			switch cell {
			case nil:
				cellLayout[j*2] = '.'
			default:
				cellLayout[j*2] = *cell
			}
			cellLayout[j*2+1] = ' '
		}
		boardLayout[i] = string(cellLayout)
	}

	fmt.Println(strings.Join(boardLayout, "\n"))

	w.Close()
	expected, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	rescueStdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	properties.(rendering.Printer).Print()

	w.Close()
	actual, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(t, string(expected), string(actual))

}
