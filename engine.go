package minesweeper

import (
	"container/list"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/rrborja/minesweeper-go/visited"
)

type eventType uint8
type blocks [][]Block

const consecutiveRandomLimit = 3

const easyMultiplier = 0.1
const mediumMultiplier = 0.2
const hardMultiplier = 0.5

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

var singleton Minesweeper
var mainEvent Event

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
	game.iterateBlocks(func(block *Block) bool {
		if block.Node == Bomb {
			game.traverseAdjacentCells(block.X(), block.Y(), func(cell *Block) {
				if cell.Node != Bomb {
					cell.Node = Number
					cell.Value++
				}
			})
		}
		return true
	})
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
	iterateNotInterrupted := game.iterateBlocks(func(block *Block) bool {
		if block.Node != Bomb && block.visited {
			visitTally++
		} else if block.Node == Bomb && block.visited {
			game.Event <- Lose
			return false
		}
		return true
	})
	if iterateNotInterrupted && visitTally == game.totalNonBombs() {
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

func (game *game) iterateBlocks(do func(*Block) bool) bool {
	success := true
mainLoop:
	for _, row := range game.blocks {
		for _, block := range row {
			if success = do(&game.blocks[block.X()][block.Y()]); !success {
				break mainLoop
			}
		}
	}
	return success
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
