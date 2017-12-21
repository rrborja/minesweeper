// Copyright 2017 Ritchie Borja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package minesweeper

import (
	"fmt"
	"testing"

	"github.com/rrborja/minesweeper-go/rendering"
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
	properties := minesweeper.(rendering.RuntimeGameProperties)

	bombPlacements := make([]rendering.Position, int(float32(game.Height*game.Width)*game.difficultyMultiplier))

	var counter int
	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				bombPlacements[counter] = rendering.Position{x, y}
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
	properties := minesweeper.(rendering.RuntimeGameProperties)

	hintPlacements := make([]rendering.Position, 0)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				hintPlacements = append(hintPlacements, rendering.Position{x, y})
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
	properties := minesweeper.(rendering.RuntimeGameProperties)

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
