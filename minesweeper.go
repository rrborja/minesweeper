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
	"sync"

	"github.com/rrborja/minesweeper-go/visited"
)

// Node is the type of the cell's value.
// Values of this type are minesweeper.Unknown, minesweeper.Bomb and
// minesweeper.Number
type Node uint8

// Grid is the game's board size defined by its Grid.Width and Grid.Height
type Grid struct{ Width, Height int }

// Difficulty is the state of the game's difficulty.
// Values of this type are minesweeper.Easy, minesweeper.Medium
// and minesweeper.Hard
type Difficulty uint8

// Event is a channel type used to receive the game's realtime event.
// Particular values accepted to this channel are minesweeper.Win and
// minesweeper.Lose
type Event chan eventType

type eventType uint8
type blocks [][]Block

const (
	// Unknown is the type of the value contained in a minesweeper
	// cell. It has no value which means no mines are neighbored in
	// the cell.
	Unknown Node = 1 << iota >> 1

	// Bomb is the type of the value contained in a minesweeper cell.
	// This is the equivalent of a mine in the game.
	Bomb

	// Number is the type of the value contained in a minesweeper
	// cell. The number indicated the number of mines neighbored in
	// the cell.
	Number
)

const (
	notSet Difficulty = iota

	// Easy is the difficulty of the game. The amount of mines present
	// in the game with this difficulty would result to 10% of the
	// total area of the board's size.
	Easy

	// Medium is the difficulty of the game. The amount of mines
	// present in the game with this difficulty would result to 20%
	// of the total area of the board's size.
	Medium

	// Hard is the difficulty of the game. The amount of mines present
	// in the game with this difficulty would result to 50% of the
	// total area of the board's size.
	Hard
)

const (
	ongoing eventType = iota

	// Win is the game's event. This will trigger whenever all non-mine
	// cells are visited.
	Win

	// Lose is the game's event. This will trigger whenever a cell
	// containing the mine is visited.
	Lose
)

const consecutiveRandomLimit = 3

const easyMultiplier = 0.1
const mediumMultiplier = 0.2
const hardMultiplier = 0.5

// Block stores the information of a particular minesweeper cell. It
// has the value itself, type of the value in the cell, and the location in
// the grid. It also reports when a cell is visited or flagged.
type Block struct {
	Node
	Value    int
	location struct {
		x int
		y int
	}
	visited, flagged bool
}

type board struct {
	*Grid
	blocks
	difficultyMultiplier float32
}

type game struct {
	Event
	board
	Difficulty
	recordedActions
	*sync.Mutex
}

// Minesweeper is the main point of consumption and manipulation of the Minesweeper's
// game state. Methods provided by this interface are common use cases to solve
// a minesweeper game.
//
// Any instance derived by this interface is compatible for type casting to the
// rendering.Tracker and visited.StoryTeller interfaces.
type Minesweeper interface {
	// Sets or changes the board's size. You can't change the board's size once
	// the Play() method has been called, otherwise, a GameAlreadyStartedError
	// will return.
	//
	// Nothing will return if the setting the board's size is successful.
	//
	// Whenever the game ends, this instance will be garbage collected. Calling
	// any exported methods when the game ended may result in nil pointer
	// dereference. Creating a new setup of the game is ideal.
	SetGrid(int, int) error

	// Sets the difficulty of the game as a basis of the number of mines. An
	// error will return if this method is being called when a game is already
	// being played or better yet, the Play() method has already been called.
	SetDifficulty(Difficulty) error

	// When called, the game sets up all the mines in place randomly. The
	// placement is non-deterministic since the implementation uses the
	// "crypto/rand" package.
	//
	// An error will return when this method is called twice or more.
	//
	// More importantly, Grid size and Difficulty must be specified, otherwise,
	// you will encounter an UnspecifiedGridError and UnspecifiedDifficultyError,
	// respectively.
	Play() error

	// When called, the cell, according to the coordinates supplied in the
	// method argument, will be marked as flagged. When a particular cell is
	// flagged, the cell in question will prevent from being handled by the
	// game when the Visit(int, int) method with the same coordinate of the
	// cell in question is called.
	Flag(int, int)

	// Visits a particular cell according to the xy-coordinates of the argument
	// supplied by this method being called. There are three scenarios that
	// depend to the generated configuration of the game:
	//
	// A warning number:
	// We may expect a number when visiting a particular cell. The returned
	// []Block will return an array of revealed cells. Accordingly, this
	// type of scenario will always return an array of one Block element.
	//
	// No mine but no warning number:
	// When encountered, a visited cell will recursively visit all neighboring
	// or adjacent cells because if a visited cells contains no warning number,
	// then there are no neighboring mines per se. Thus, visiting all the
	// neighboring cells are safe to be visited and continues to visit any
	// probed blank cell recursively until no blank cell is left to be probed.
	// The returned []Block will return an array of all probed cells with the
	// first element as the original visited cell.
	//
	// A mine:
	// When encountered, the game ends revealing all mines by the returned
	// []Block array. It will also return an ExplodedError as a second return
	// value. Furthermore, when event handling is supported, a Lose event will
	// enqueue to the game's even channel buffer. You can remember this buffer
	// when you originally called the NewGame(...Grid) function and store the second
	// returned value as your event listener in a separate goroutine.
	//
	// There is also one trivial scenario. Visiting a block that is flagged
	// (i.e. originally called the Flag(int, int) method as a result) will have
	// no effect on the state of the game. If a flagged cell with a suspected mine
	// is visited, it will prevent it from being visited and the game would likely
	// treat the called method as if nothing was called at all.
	//
	// The last Visit() method call with the last non-mine cell will trigger
	// the Win event. The game ends eventually.
	Visit(int, int) ([]Block, error)
}

// NewGame Creates a new minesweeper instance. Note that this only creates the minesweeper
// instance without the necessary settings such as the game's difficulty and the
// game's board size and calling this method will not start the game.
//
// This method returns the minesweeper interface and the event handler. The event
// handler allows you to setup event listeners in such situation when a game
// seamlessly triggers a win or lose event. The event handler is a buffered
// channel that, when used, allows you to setup a particular goroutine to
// independently listen for these events.
//
// As for this method's argument, this method appears to accept an arbitrary
// number of trailing arguments of type Grid. It can only, however, handle only
// one Grid instance and the rest of the arguments will be ignored. Although
// supplying this Grid is also optional, you may encounter an UnspecifiedGridError
// panic when calling the Play() method if the Grid is not supplied. You may
// explicitly supply it by calling the SetGrid(int, int) method.
func NewGame(grid ...Grid) (Minesweeper, Event) {
	game := new(game)

	if len(grid) > 0 {
		game.SetGrid(grid[0].Width, grid[0].Height)
	}

	game.Event = make(chan eventType, 1)

	return game, game.Event
}

func (game *game) SetGrid(width, height int) error {
	if game.Grid != nil {
		return new(GameAlreadyStartedError)
	}
	game.Grid = &Grid{width, height}
	createBoard(game)
	return nil
}

func (game *game) Flag(x, y int) {
	blockPtr := &game.blocks[x][y]
	if !blockPtr.visited {
		blockPtr.flagged = !blockPtr.flagged
	}
}

func (game *game) Visit(x, y int) ([]Block, error) {
	game.validateGameEnvironment()

	game.Lock()
	defer game.Unlock()

	block := &game.blocks[x][y]
	if block.Node == Number && block.visited {
		countedFlaggedBlock := 0
		resultedBlocks := make([]Block, 0)
		blocksToBeVisited := make([]*Block, 0)

		game.traverseAdjacentCells(x, y, func(cell *Block) {
			if cell.flagged {
				countedFlaggedBlock++
			} else {
				blocksToBeVisited = append(blocksToBeVisited, cell)
			}
		})

		if countedFlaggedBlock == block.Value {
			for _, block := range blocksToBeVisited {
				blocks, err := game.visit(block.X(), block.Y())
				if err != nil {
					return blocks, err
				}
				resultedBlocks = append(resultedBlocks, blocks...)
			}
		}
		return resultedBlocks, nil
	}
	return game.visit(x, y)
}

func (game *game) visit(x, y int) ([]Block, error) {
	block := &game.blocks[x][y]

	if !block.flagged && !block.visited {
		block.visited = true
		defer func() {
			go game.validateSolution()
		}()
		switch block.Node {
		case Number:
			defer game.add(visited.Record{
				Position: *block, Action: visited.Number})
			return []Block{*block}, nil
		case Bomb:
			defer game.add(visited.Record{
				Position: *block, Action: visited.Bomb})

			bombLocations := make([]Block, 0, game.totalBombs()-1)

			for _, bombLocation := range game.BombLocations() {
				if bombLocation != *block {
					bombLocations = append(bombLocations, bombLocation.(Block))
				}
			}

			bombLocations = append([]Block{*block}, bombLocations...)

			return bombLocations, &ExplodedError{x: x, y: y}
		case Unknown:
			defer game.add(visited.Record{
				Position: *block, Action: visited.Unknown})
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

func (game *game) SetDifficulty(difficulty Difficulty) error {
	if game.Mutex != nil {
		return new(GameAlreadyStartedError)
	}

	game.Difficulty = difficulty
	switch difficulty {
	case Easy:
		game.difficultyMultiplier = easyMultiplier
	case Medium:
		game.difficultyMultiplier = mediumMultiplier
	case Hard:
		game.difficultyMultiplier = hardMultiplier
	}

	return nil
}

func (game *game) Play() error {
	if game.Mutex != nil {
		return new(GameAlreadyStartedError)
	}
	game.Mutex = new(sync.Mutex)

	if game.Difficulty == notSet {
		return new(UnspecifiedDifficultyError)
	}
	if game.Grid == nil {
		return new(UnspecifiedGridError)
	}
	createBombs(game)
	tallyHints(game)
	return nil
}

// X returns the X coordinate of the block in the minesweeper grid
func (block Block) X() int {
	return block.location.x
}

// Y returns the Y coordinate of the block in the minesweeper grid
func (block Block) Y() int {
	return block.location.y
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
			for game.board.blocks[x][y].Node != Unknown {
				x, y = shiftPosition(game.Grid, x, y)
				countLimit++
			}

			if countLimit <= consecutiveRandomLimit {
				game.blocks[x][y].Node = Bomb
				break
			}
		}
	}
}

func tallyHints(game *game) {
	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				game.traverseAdjacentCells(x, y, func(cell *Block) {
					if cell.Node != Bomb {
						cell.Node = Number
						cell.Value++
					}
				})
			}
		}
	}
}

func createBoard(game *game) {
	game.blocks = make([][]Block, game.Width)
	for x := range game.blocks {
		game.blocks[x] = make([]Block, game.Height)
	}
	for x, row := range game.blocks {
		for y := range row {
			game.blocks[x][y].location = struct{ x, y int }{x: x, y: y}
		}
	}
}

func autoRevealUnmarkedBlock(game *game, visitedBlocks *list.List, x, y int) {
	blocks := game.blocks
	width := game.Width
	height := game.Height

	if x >= 0 && y >= 0 && x < width && y < height {
		if blocks[x][y].visited {
			return
		}
		if blocks[x][y].Node == Unknown {
			blocks[x][y].visited = true

			visitedBlocks.PushBack(blocks[x][y])

			game.traverseAdjacentCells(x, y, func(cell *Block) {
				autoRevealUnmarkedBlock(game, visitedBlocks, cell.X(), cell.Y())
			})

		} else if blocks[x][y].Node == Number {
			blocks[x][y].visited = true

			visitedBlocks.PushBack(blocks[x][y])
		}
	}
}

func (game *game) validateSolution() {
	var visitTally int
	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node != Bomb && block.visited {
				visitTally++
			} else if block.Node == Bomb && block.visited {
				game.Event <- Lose
				return
			}
		}
	}
	if visitTally == game.totalNonBombs() {
		game.Event <- Win
	}
}

func (game *game) validateGameEnvironment() {
	if game.Grid == nil {
		panic(UnspecifiedGridError{})
	}
	if game.Difficulty == notSet {
		panic(UnspecifiedDifficultyError{})
	}
}

func (game *game) traverseAdjacentCells(x, y int, do func(*Block)) {
	game.recursivelyTraverseAdjacentCells(x-1, y-1, do)
	game.recursivelyTraverseAdjacentCells(x, y-1, do)
	game.recursivelyTraverseAdjacentCells(x+1, y-1, do)
	game.recursivelyTraverseAdjacentCells(x-1, y, do)
	game.recursivelyTraverseAdjacentCells(x+1, y, do)
	game.recursivelyTraverseAdjacentCells(x-1, y+1, do)
	game.recursivelyTraverseAdjacentCells(x, y+1, do)
	game.recursivelyTraverseAdjacentCells(x+1, y+1, do)
}

func (game *game) recursivelyTraverseAdjacentCells(x, y int, do func(*Block)) {
	width := game.Width
	height := game.Height

	if x >= 0 && y >= 0 && x < width && y < height {
		do(&game.blocks[x][y])
	}
}

func (game *game) area() int {
	return len(game.blocks) * len(game.blocks[0])
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

// Visited responds if a cell is visited or not
func (block *Block) Visited() bool {
	return block.visited
}

// Flagged responds if a cell is visited or not
func (block *Block) Flagged() bool {
	return block.flagged
}

func (block Block) String() string {
	var nodeType string
	switch block.Node {
	case Unknown:
		nodeType = "blank"
	case Number:
		nodeType = "number"
	case Bomb:
		nodeType = "bomb"
	}

	var value string
	if block.Value > 0 {
		value = string(block.Value)
	}

	return fmt.Sprintf("\n\nBlock: \n\tValue\t :\t%v\n\tLocation :\tx:%v y:%v\n\tType\t :\t%v\n\tVisited? :\t%v\n\tFlagged? :\t%v\n\n",
		value, block.location.x, block.location.y, nodeType, block.visited, block.flagged)
}

func randomNumber(max int) int {
	var number uint16
	binary.Read(rand.Reader, binary.LittleEndian, &number)
	return int(number) % max
}
