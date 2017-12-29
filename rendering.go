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
	"strings"

	"github.com/rrborja/minesweeper/rendering"
	"github.com/rrborja/minesweeper/visited"
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

func (game *game) Print() {
	bombs := game.BombLocations()
	hints := game.HintLocations()

	star := '*'

	var board = make([][]*rune, game.Width)
	for i := range board {
		board[i] = make([]*rune, game.Height)
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

	var boardLayout = make([]string, game.Width)
	for i, row := range board {
		cellLayout := make([]rune, (game.Height * 2))
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
