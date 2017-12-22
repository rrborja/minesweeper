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
	"container/list"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/rrborja/minesweeper-go/visited"
	"sync"
)

type Node uint8
type Blocks [][]Block
type Grid struct{ Width, Height int }

type Difficulty uint8

type EventType uint8
type Event chan EventType

const (
	UNKNOWN Node = 1 << iota >> 1
	BOMB
	NUMBER
)

const (
	NOTSET Difficulty = iota
	EASY
	MEDIUM
	HARD
)

const (
	ONGOING EventType = 1 << iota
	WIN
	LOSE
)

const CONSECUTIVE_RANDOM_LIMIT = 3

const EASY_MULTIPLIER = 0.1
const MEDIUM_MULTIPLIER = 0.2
const HARD_MULTIPLIER = 0.5

type Block struct {
	Node
	Value    int
	Location struct {
		X int
		Y int
	}
	visited, flagged bool
}

type Board struct {
	*Grid
	Blocks
	difficultyMultiplier float32
}

type game struct {
	Event
	Board
	Difficulty
	RecordedActions
	*sync.Mutex
}

type Minesweeper interface {
	SetGrid(int, int) error

	SetDifficulty(Difficulty)

	Play() error

	Flag(int, int)

	Visit(int, int) ([]Block, error)
}

func NewGame(grid ...Grid) (Minesweeper, Event) {
	game := new(game)

	if len(grid) > 0 {
		game.SetGrid(grid[0].Width, grid[0].Height)
	}

	game.Mutex = new(sync.Mutex)
	game.Event = make(chan EventType, 1)

	return game, game.Event
}

func (game *game) SetGrid(width, height int) error {
	if game.Grid != nil {
		return new(GameAlreadyStarted)
	}
	game.Grid = &Grid{width, height}
	createBoard(game)
	return nil
}

func (game *game) Flag(x, y int) {
	game.Blocks[x][y].flagged = true
}

func (game *game) Visit(x, y int) ([]Block, error) {
	game.validateGameEnvironment()

	game.Lock()
	defer game.Unlock()

	block := &game.Blocks[x][y]

	if !block.flagged && !block.visited {
		block.visited = true
		defer func() {
			go game.validateSolution()
		}()
		switch block.Node {
		case NUMBER:
			defer game.Add(visited.Record{*block, visited.Number})
			return []Block{*block}, nil
		case BOMB:
			defer game.Add(visited.Record{*block, visited.Bomb})

			bombLocations := make([]Block, 0, game.totalBombs()-1)

			for _, bombLocation := range game.BombLocations() {
				if bombLocation != *block {
					bombLocations = append(bombLocations, bombLocation.(Block))
				}
			}

			bombLocations = append([]Block{*block}, bombLocations...)

			return bombLocations, &Exploded{struct{ x, y int }{x: x, y: y}}
		case UNKNOWN:
			defer game.Add(visited.Record{*block, visited.Unknown})
			block.visited = false //to avoid infinite recursion, first is to set the base case

			visitedList := list.New()
			autoRevealUnmarkedBlock(game, visitedList, x, y)

			visitedBlocks := make([]Block, visitedList.Len())

			var counter int
			for e := visitedList.Front(); e != nil; e = e.Next() {
				visitedBlocks[counter] = e.Value.(Block)
				counter++
			}

			return visitedBlocks, nil
		}
	}
	return nil, nil
}

func (game *game) SetDifficulty(difficulty Difficulty) {
	game.Difficulty = difficulty
	switch difficulty {
	case EASY:
		game.difficultyMultiplier = EASY_MULTIPLIER
	case MEDIUM:
		game.difficultyMultiplier = MEDIUM_MULTIPLIER
	case HARD:
		game.difficultyMultiplier = HARD_MULTIPLIER
	}
}

func (game *game) Play() error {
	if game.Difficulty == NOTSET {
		return new(UnspecifiedDifficulty)
	}
	if game.Grid == nil {
		return new(UnspecifiedGrid)
	}
	createBombs(game)
	tallyHints(game)
	return nil
}

func (block Block) X() int {
	return block.Location.X
}

func (block Block) Y() int {
	return block.Location.Y
}

// Shifts to the right
func shiftPosition(grid *Grid, x, y int) (_x, _y int) {
	width := grid.Width
	height := grid.Height
	if x+1 >= width {
		if y+1 >= height {
			_x, _y = 0, 0
		} else {
			_x, _y = 0, y+1
		}
	} else {
		_x, _y = x+1, y
	}
	return
}

func createBombs(game *game) {
	area := int(game.Width * game.Height)
	for i := 0; i < int(float32(area)*game.difficultyMultiplier); i++ {
		for {
			randomPos := randomNumber(area)

			x, y := randomPos%game.Width, randomPos/game.Width

			countLimit := 0
			for game.Board.Blocks[x][y].Node != UNKNOWN {
				x, y = shiftPosition(game.Grid, x, y)
				countLimit++
			}

			if countLimit <= CONSECUTIVE_RANDOM_LIMIT {
				game.Blocks[x][y].Node = BOMB
				break
			}
		}
	}
}

func tallyHints(game *game) {
	width := game.Width
	height := game.Height

	tally := func(blocks Blocks, x, y int) {
		if x >= 0 && y >= 0 &&
			x < width && y < height &&
			blocks[x][y].Node != BOMB {
			blocks[x][y].Node = NUMBER
			blocks[x][y].Value++
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				tally(game.Blocks, x-1, y-1)
				tally(game.Blocks, x-1, y)
				tally(game.Blocks, x-1, y+1)
				tally(game.Blocks, x, y-1)
				tally(game.Blocks, x, y+1)
				tally(game.Blocks, x+1, y-1)
				tally(game.Blocks, x+1, y)
				tally(game.Blocks, x+1, y+1)
			}
		}
	}
}

func createBoard(game *game) {
	game.Blocks = make([][]Block, game.Width)
	for x := range game.Blocks {
		game.Blocks[x] = make([]Block, game.Height)
	}
	for x, row := range game.Blocks {
		for y := range row {
			game.Blocks[x][y].Location = struct{ X, Y int }{X: x, Y: y}
		}
	}
}

func autoRevealUnmarkedBlock(game *game, visitedBlocks *list.List, x, y int) {
	blocks := game.Blocks
	width := game.Width
	height := game.Height

	if x >= 0 && y >= 0 && x < width && y < height {
		if blocks[x][y].visited {
			return
		}
		if blocks[x][y].Node == UNKNOWN {
			blocks[x][y].visited = true

			visitedBlocks.PushBack(blocks[x][y])

			autoRevealUnmarkedBlock(game, visitedBlocks, x-1, y-1)
			autoRevealUnmarkedBlock(game, visitedBlocks, x-1, y)
			autoRevealUnmarkedBlock(game, visitedBlocks, x-1, y+1)
			autoRevealUnmarkedBlock(game, visitedBlocks, x, y-1)
			autoRevealUnmarkedBlock(game, visitedBlocks, x, y+1)
			autoRevealUnmarkedBlock(game, visitedBlocks, x+1, y-1)
			autoRevealUnmarkedBlock(game, visitedBlocks, x+1, y)
			autoRevealUnmarkedBlock(game, visitedBlocks, x+1, y+1)
		} else if blocks[x][y].Node == NUMBER {
			blocks[x][y].visited = true

			visitedBlocks.PushBack(blocks[x][y])
		}
	}
}

func (game *game) validateSolution() {
	var visitTally int
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node != BOMB && block.visited {
				visitTally++
			} else if block.Node == BOMB && block.visited {
				game.Event <- LOSE
				return
			}
		}
	}
	if visitTally == game.totalNonBombs() {
		game.Event <- WIN
	}
}

func (game *game) validateGameEnvironment() {
	if game.Grid == nil {
		panic(UnspecifiedGrid{})
	}
	if game.Difficulty == NOTSET {
		panic(UnspecifiedDifficulty{})
	}
}

func (game *game) area() int {
	return len(game.Blocks) * len(game.Blocks[0])
}

func (game *game) areaInFloat() float32 {
	return float32(game.area())
}

func (game *game) totalBombs() int {
	return int(game.areaInFloat() * game.difficultyMultiplier)
}

func (game *game) totalNonBombs() int {
	return game.area() - game.totalBombs()
}

func (block Block) String() string {
	var nodeType string
	switch block.Node {
	case UNKNOWN:
		nodeType = "blank"
	case NUMBER:
		nodeType = "number"
	case BOMB:
		nodeType = "bomb"
	}

	var value string
	if block.Value > 0 {
		value = string(block.Value)
	}

	return fmt.Sprintf("\n\nBlock: \n\tValue\t :\t%v\n\tLocation :\tx:%v y:%v\n\tType\t :\t%v\n\tVisited? :\t%v\n\tFlagged? :\t%v\n\n",
		value, block.Location.X, block.Location.Y, nodeType, block.visited, block.flagged)
}

func randomNumber(max int) int {
	var number uint16
	binary.Read(rand.Reader, binary.LittleEndian, &number)
	return int(number) % max
}
