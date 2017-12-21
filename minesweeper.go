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
	"container/list"
	"crypto/rand"
	"encoding/binary"
	"sync"
)

type Node uint8
type Blocks [][]Block
type Grid struct{ Width, Height int }

type Difficulty uint8

type EventType uint8
type Event chan EventType

var eventLock sync.Mutex

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
	Value            int
	Location         Position
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
	if !game.Blocks[x][y].flagged {
		game.Blocks[x][y].visited = true
		defer func() {
			go game.validateSolution()
		}()
		switch game.Blocks[x][y].Node {
		case NUMBER:
			return []Block{game.Blocks[x][y]}, nil
		case BOMB:
			return []Block{game.Blocks[x][y]}, &Exploded{struct{ x, y int }{x: x, y: y}}
		case UNKNOWN:
			game.Blocks[x][y].visited = false //to avoid infinite recursion, first is to set the base case

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
	createBombs(game)
	tallyHints(game)
	return nil
}

func (block *Block) SetBlock(node Node) {
	block.Node = node
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
		randomPos := randomNumber(area)

		x, y := randomPos%game.Width, randomPos/game.Width

		countLimit := 0
		for game.Board.Blocks[x][y].Node != UNKNOWN {
			x, y = shiftPosition(game.Grid, x, y)
			countLimit++
		}
		if countLimit >= CONSECUTIVE_RANDOM_LIMIT {
			i--
			continue
		}

		game.Blocks[x][y].Node = BOMB
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
			game.Blocks[x][y].Location = Position{x, y}
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
	eventLock.Lock()
	defer eventLock.Unlock()

	var visitTally int
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node != BOMB && block.visited {
				visitTally++
			}
		}
	}
	if visitTally == game.totalNonBombs() {
		game.Event <- WIN
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

func randomNumber(max int) int {
	var number uint16
	binary.Read(rand.Reader, binary.LittleEndian, &number)
	return int(number) % max
}
