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
	"github.com/rrborja/minesweeper-go/rendering"
	"github.com/rrborja/minesweeper-go/visited"
)

type recordedActions struct {
	*visited.History
}

func (game *game) BombLocations() []rendering.Position {
	bombPlacements := make([]rendering.Position, int(float32(game.Height*game.Width)*game.difficultyMultiplier))

	var counter int
	game.iterateBlocks(func(block *Block) bool {
		if block.Node == Bomb {
			bombPlacements[counter] = *block
			counter++
		}
		return true
	})

	return bombPlacements
}

// Not recommended to call this function until a new update to improve the performance of this method
func (game *game) HintLocations() []rendering.Position {
	hintPlacements := make([]rendering.Position, 0) // TODO: Improve this performance

	game.iterateBlocks(func(block *Block) bool {
		if block.Node == Number {
			hintPlacements = append(hintPlacements, *block)
		}
		return true
	})

	return hintPlacements
}

func (game *game) History() *visited.History {
	return game.recordedActions.History
}

func (game *game) LastAction() visited.Record {
	return game.recordedActions.History.Record
}

func (game *recordedActions) add(record visited.Record) {
	if game.History == nil {
		game.History = new(visited.History)
	} else {
		temp := game.History
		game.History = new(visited.History)
		game.History.History = temp
	}
	game.Record = record
}
