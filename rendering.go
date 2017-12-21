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

import "github.com/rrborja/minesweeper-go/rendering"

func (game *game) BombLocations() []rendering.Position {
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

	return bombPlacements
}

// Not recommended to call this function until a new update to improve the performance of this method
func (game *game) HintLocations() []rendering.Position {
	hintPlacements := make([]rendering.Position, 0) // TODO: Improve this performance

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				hintPlacements = append(hintPlacements, rendering.Position{x, y})
			}
		}
	}

	return hintPlacements
}

func (game *game) History() rendering.History {
	return nil
}

func (game *game) LastAction() rendering.Record {
	return rendering.Record{}
}
