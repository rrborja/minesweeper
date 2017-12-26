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
	"time"

	"github.com/rrborja/minesweeper-go/visited"
	"github.com/stretchr/testify/assert"
)

const (
	sampleGridWidth  = 10
	sampleGridHeight = 40
)

func newBlankGame() Minesweeper {
	game, _ := NewGame()
	return game
}

func newSampleGame() Minesweeper {
	game, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	return game
}

func TestGridMustNotBeSquaredForTheSakeOfTesting(t *testing.T) {
	assert.True(t, sampleGridWidth != sampleGridHeight)
}

func TestGame_SetGrid(t *testing.T) {
	minesweeper := newBlankGame()
	minesweeper.SetGrid(sampleGridWidth, sampleGridHeight)
	assert.Equal(t, minesweeper.(*game).board.Grid, &Grid{sampleGridWidth, sampleGridHeight})
}

func TestGameWithGridArgument(t *testing.T) {
	minesweeper := newSampleGame()
	assert.Equal(t, minesweeper.(*game).board.Grid, &Grid{sampleGridWidth, sampleGridHeight})
}

func TestNewGridWhenStartedGame(t *testing.T) {
	minesweeper := newSampleGame()
	err := minesweeper.SetGrid(10, 20)
	assert.Error(t, err)
	assert.NotNil(t, err, "Must report an error upon setting a new grid from an already started game")
	assert.IsType(t, new(GameAlreadyStartedError), err, "The error must be GameAlreadyStarted error type")
}

func TestFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.Flag(3, 6)
	assert.Equal(t, minesweeper.(*game).blocks[3][6].flagged, true)
}

func TestGame_SetDifficulty(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	assert.Equal(t, minesweeper.(*game).Difficulty, Easy)
}

func TestShiftFromMaxPosition(t *testing.T) {
	grid := Grid{5, 5}
	x, y := shiftPosition(&grid, 4, 4)
	assert.Equal(t, struct {
		x int
		y int
	}{0, 0}, struct {
		x int
		y int
	}{x, y})
}

func TestBombsInPlace(t *testing.T) {

	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	numOfBombs := int(float32(game.Width*game.Height) * easyMultiplier)
	countedBombs := 0
	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node == Bomb {
				countedBombs++
			}
		}
	}
	assert.Equal(t, numOfBombs, countedBombs)
}

func TestTalliedBomb(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)
	width := game.Width
	height := game.Height

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x-1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x-1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x-1, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x+1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x+1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.blocks, width, height, x+1, y+1))
			}
		}
	}

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Number {
				var counted int
				counted = count(game.blocks, width, height, x-1, y-1) +
					count(game.blocks, width, height, x-1, y) +
					count(game.blocks, width, height, x-1, y+1) +
					count(game.blocks, width, height, x, y-1) +
					count(game.blocks, width, height, x, y+1) +
					count(game.blocks, width, height, x+1, y-1) +
					count(game.blocks, width, height, x+1, y) +
					count(game.blocks, width, height, x+1, y+1)
				assert.Equal(t, counted, block.Value)
			}
		}
	}
}

func TestVisitedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Number {
				game.Visit(x, y)
				assert.True(t, game.blocks[x][y].visited)
			}
		}
	}

}

func TestVisitedBombToGameOver(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)
	var x, y int
	var err error

	for i, row := range game.blocks {
		for j, block := range row {
			if block.Node == Bomb {
				x, y = i, j
				_, err = game.Visit(x, y)
				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.IsType(t, new(ExplodedError), err)
			}
		}
	}

}

func TestVisitedBombToGameOverWithCorrectLocationReason(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)
	var x, y int
	var err error

	for i, row := range game.blocks {
		for j, block := range row {
			if block.Node == Bomb {
				x, y = i, j
				_, err = game.Visit(x, y)
				assert.Error(t, err)
				assert.EqualError(t, err,
					fmt.Sprintf("Game over at X=%v Y=%v",
						x, y))
				assert.IsType(t, new(ExplodedError), err)
			}
		}
	}

}

func TestVisitedUnmarkedBlockDistributeVisit(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Unknown && !block.visited {
				minesweeper.Visit(x, y)
			}
		}
	}

	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node == Unknown {
				assert.True(t, block.visited)
			}
		}
	}
}

func TestVisitAFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				minesweeper.Flag(x, y)
				_, err := minesweeper.Visit(x, y)
				assert.NoError(t, err)
				if err != nil {
					assert.IsType(t, new(ExplodedError), err)
				}
			}
		}
	}
}

func TestVisitedBlocksReturnOneBlockWhenAHintBlockIsVisited(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Number {
				visitedBlocks, err := minesweeper.Visit(x, y)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(visitedBlocks))
				assert.Equal(t, block.Value, visitedBlocks[0].Value)
				assert.Equal(t, visitedBlocks[0], game.blocks[x][y])
			}
		}
	}
}

func TestVisitedBlocksWhenBlockIsABomb(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				_, err := minesweeper.Visit(x, y)
				assert.Error(t, err)
				assert.EqualError(t, err, (&ExplodedError{x: x, y: y}).Error())
			}
		}
	}
}

func TestVisitedBlockWhenBlockIsUnknownAndSpreadVisits(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	var x, y int
	var actualVisitedBlocks []Block
first:
	for i, row := range game.blocks {
		for j, block := range row {
			if block.Node == Unknown && !block.visited {
				x, y = i, j
				var err error
				actualVisitedBlocks, err = minesweeper.Visit(x, y)
				assert.NoError(t, err)
				break first
			}
		}
	}

	var visitedBlocks []Block
	for _, row := range game.blocks {
		for _, block := range row {
			if block.visited {
				visitedBlocks = append(visitedBlocks, block)
			}
		}
	}

	assert.NotEmpty(t, actualVisitedBlocks)

	for _, block1 := range visitedBlocks {
		found := false
		for _, block2 := range actualVisitedBlocks {
			if block1 == block2 {
				found = true
				break
			}
		}
		assert.Truef(t, found, "%v not found in list %v", block1, actualVisitedBlocks)
	}
}

func TestBlockLocationAfterNewGame(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				assert.Equal(t, struct{ x, y int }{x: x, y: y}, block.location)
			}
		}
	}

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Number {
				assert.Equal(t, struct{ x, y int }{x: x, y: y}, block.location)
			}
		}
	}

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Unknown {
				assert.Equal(t, struct{ x, y int }{x: x, y: y}, block.location)
			}
		}
	}
}

func TestCheckEventOfGameWhenWinning(t *testing.T) {
	minesweeper, event := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node != Bomb && !block.visited {
				minesweeper.Visit(x, y)
			}
		}
	}

	go func() {
		time.Sleep(5 * time.Second)
		assert.Fail(t, "Was expecting any event in less than 5 seconds of runtime")
		close(event)
	}()

	if won, ok := <-event; ok {
		assert.Equal(t, Win, won, "Expecting a winning event")
	} else {
		assert.Fail(t, "Channel event closed. Broken code.")
	}
}

func TestCheckEventOfGameWhenLosing(t *testing.T) {
	minesweeper, event := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

mainLoop:
	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				minesweeper.Visit(x, y)
				break mainLoop
			}
		}
	}

	go func() {
		time.Sleep(5 * time.Second)
		assert.Fail(t, "Was expecting any event in less than 5 seconds of runtime")
		close(event)
	}()

	if won, ok := <-event; ok {
		assert.Equal(t, Lose, won, "Expecting a losing event")
	} else {
		assert.Fail(t, "Channel event closed. Broken code.")
	}
}

func TestGameEasyDifficultyIsSet(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	assert.Equal(t, Easy, game.Difficulty)
	assert.Equal(t, int(sampleGridWidth*sampleGridHeight*easyMultiplier), game.totalBombs())
}

func TestGameMediumDifficultyIsSet(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Medium)
	minesweeper.Play()

	game := minesweeper.(*game)

	assert.Equal(t, Medium, game.Difficulty)
	assert.Equal(t, int(sampleGridWidth*sampleGridHeight*mediumMultiplier), game.totalBombs())
}

func TestGameHardDifficultyIsSet(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Hard)
	minesweeper.Play()

	game := minesweeper.(*game)

	assert.Equal(t, Hard, game.Difficulty)
	assert.Equal(t, int(sampleGridWidth*sampleGridHeight*hardMultiplier), game.totalBombs())
}

func TestGameOverReturnAllBombLocations(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	var bombLocations []Block

mainLoop:
	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				bombLocations, _ = minesweeper.Visit(x, y)
				break mainLoop
			}
		}
	}

	assert.Equalf(t, game.totalBombs(), len(bombLocations), "Number of bombs must be %v", game.totalBombs())

	for _, bombLocation := range bombLocations {
		x := bombLocation.location.x
		y := bombLocation.location.y
		assert.Equalf(t, game.blocks[x][y].Node, Bomb, "Block at %v:%v is not a bomb.", x, y)
	}
}

func TestPlayGameWithoutSettingDifficulty(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	err := minesweeper.Play()

	assert.Error(t, err)
	assert.Equal(t, new(UnspecifiedDifficultyError), err)
}

func TestPlayGameWithoutSettingGrid(t *testing.T) {
	minesweeper, _ := NewGame()
	minesweeper.SetDifficulty(Hard)
	err := minesweeper.Play()

	assert.Nil(t, minesweeper.(*game).Grid,
		"For the sake of testing this, we expect Grid is not specified. Therefore this test must fail.")
	assert.Error(t, err)
	assert.Equal(t, new(UnspecifiedGridError), err)
}

func TestVisitedUnknownIsTheFirstInTheListOfDistributedVisits(t *testing.T) {
	minesweeper, _ := NewGame()
	minesweeper.SetDifficulty(Hard)
	minesweeper.Play()

	game := minesweeper.(*game)

	for i, row := range game.blocks {
		for j, block := range row {
			if block.Node == Unknown && !block.visited {
				visitedBlocks, _ := minesweeper.Visit(i, j)
				assert.Equal(t, block, visitedBlocks[0])
			}
		}
	}
}

func TestGameDoesRecordPlayersAction(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Medium)
	minesweeper.Play()

	var expectedHistory []visited.Record
	var story visited.StoryTeller = minesweeper.(*game)

	maxMoves := 10
mainLoop:
	for i := 0; i < maxMoves; i++ {
		randomX := randomNumber(sampleGridWidth)
		randomY := randomNumber(sampleGridHeight)
		blocks, _ := minesweeper.Visit(randomX, randomY)

		if blocks == nil { // Either already visited block or flagged block
			continue
		}

		switch len(blocks) {
		case 0: // Either already visited block or flagged block
			// panic("Unexpected")
			// TODO: Action for auto visit
		case 1: // Number
			switch blocks[0].Node {
			case Bomb:
				expectedHistory = append(expectedHistory, visited.Record{Position: blocks[0], Action: visited.Bomb})
				break mainLoop
			case Number:
				expectedHistory = append(expectedHistory, visited.Record{Position: blocks[0], Action: visited.Number})
			default:
				panic("Unexpected")
			}
		default: // Unknown
			if blocks[0].Node == Bomb {
				expectedHistory = append(expectedHistory, visited.Record{Position: blocks[0], Action: visited.Bomb})
			} else if blocks[0].Node == Unknown {
				expectedHistory = append(expectedHistory, visited.Record{Position: blocks[0], Action: visited.Unknown})
			} else {
				panic("Unexpected")
			}
		}
	}

	assert.NotNil(t, story.History(), "Initial phase of comparing list must pass")

	for cursor, i := story.History(), len(expectedHistory)-1; cursor != nil && i >= 0; cursor, i = cursor.History, i-1 {
		assert.Equal(t, expectedHistory[i], cursor.Record)
	}

}

// Test cancelled by virtue of auto filling a visited block. Cannot be oblivious as it will do something
//
// func TestRevisitedBlockDoCompletelyOblivious(t *testing.T) {
// 	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
// 	minesweeper.SetDifficulty(Easy)
// 	minesweeper.Play()

// 	game := minesweeper.(*game)

// 	for x, row := range game.blocks {
// 		for y, block := range row {
// 			if block.Node != Bomb {
// 				minesweeper.Visit(x, y)
// 				result, _ := minesweeper.Visit(x, y) // Visit again. Point of this test.
// 				assert.Nil(t, result, "Game must be oblivious of a visited block.")
// 			}
// 		}
// 	}
// }

func TestBlock_String(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for _, row := range game.blocks {
		for _, block := range row {
			var expectedType string
			switch block.Node {
			case Unknown:
				expectedType = "blank"
			case Number:
				expectedType = "number"
			case Bomb:
				expectedType = "bomb"
			}

			var value string
			if block.Value > 0 {
				value = string(block.Value)
			}

			assert.Equal(t,
				fmt.Sprintf("\n\nBlock: \n\tValue\t :\t%v\n\tLocation :\tx:%v y:%v\n\tType\t :\t%v\n\tVisited? :\t%v\n\tFlagged? :\t%v\n\n",
					value, block.location.x, block.location.y, expectedType, block.visited, block.flagged),
				block.String())
		}
	}
}

func TestAttemptVisitWithoutSettingUpGameEnvironmentOfGrid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, UnspecifiedGridError{}, r.(error).Error())
		} else {
			assert.Fail(t, "We are expecting an error when grid is not set.")
		}
	}()

	minesweeper, _ := NewGame()
	minesweeper.SetDifficulty(Hard)
	minesweeper.Play()

	minesweeper.Visit(0, 0)
}

func TestAttemptVisitWithoutSettingUpGameEnvironmentOfDifficulty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, UnspecifiedDifficultyError{}, r.(error).Error())
		} else {
			assert.Fail(t, "We are expecting an error when difficulty is not set.")
		}
	}()

	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.Play()

	minesweeper.Visit(0, 0)
}

func TestRepeatThePlayMethodThenReturnError(t *testing.T) {
	go func() {
		time.Sleep(5 * time.Second)
		t.Fail()
		panic("Testing timeout. Recursion loop in tallyHints()")
	}()

	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Medium)

	minesweeper.Play()

	assert.EqualError(t, minesweeper.Play(), GameAlreadyStartedError{}.Error())
}

func TestCannotChangeDifficultyOnceGameIsStarted(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Medium)

	minesweeper.Play()

	assert.EqualError(t, minesweeper.SetDifficulty(Hard), GameAlreadyStartedError{}.Error())
}

func TestBlock_Visited(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node != Bomb && !block.visited {
				minesweeper.Visit(x, y)
				assert.True(t, game.blocks[x][y].Visited())
			}
		}
	}
}

func TestBlock_Flagged(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb {
				minesweeper.Flag(x, y)
				assert.True(t, game.blocks[x][y].Flagged())
			}
		}
	}
}

func TestVisitedFlaggedBlock(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node != Bomb && !block.visited {
				minesweeper.Visit(x, y)
				minesweeper.Flag(x, y)
				assert.False(t, game.blocks[x][y].flagged)
			}
		}
	}
}

func TestUnflagBlock(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Bomb && !block.flagged {
				minesweeper.Flag(x, y)
				assert.True(t, game.blocks[x][y].flagged)
				minesweeper.Flag(x, y)
				assert.False(t, game.blocks[x][y].flagged)
			}
		}
	}
}

func TestAutoVisitNeighboringUnprobedNumbersAfterMineFlagged(t *testing.T) {
	minesweeper, event := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	go func() {
		if <-event == Lose {
			assert.FailNow(t, "Broken test")
		}
	}()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Number && !block.visited && !block.flagged {
				minesweeper.Visit(x, y)

				bombsNearby := make([]*Block, 0, 8)
				nonBombs := make([]*Block, 0, 8)

				game.traverseAdjacentCells(x, y, func(cell *Block) {
					bombsNearby, nonBombs = appendBomb(game, cell.X(), cell.Y(), bombsNearby, nonBombs)
				})

				for _, bomb := range bombsNearby {
					minesweeper.Flag(bomb.X(), bomb.Y())
				}

				// Visit a visited cell again after flagging all neighboring mine to auto visit neighboring cells
				minesweeper.Visit(x, y)

				for _, nonBomb := range nonBombs {

					assert.True(t, nonBomb.visited)

				}
				return
			}
		}
	}
}

func TestAutoVisitWithWronglyFlaggedBomb(t *testing.T) {
	minesweeper, _ := NewGame(Grid{10, 10})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {

			if block.Node == Number && !block.visited && !block.flagged {
				minesweeper.Visit(x, y)

				var numberCellToFlag *Block
				var bombCellToNotFlag *Block

				game.traverseAdjacentCells(x, y, func(cell *Block) {
					if numberCellToFlag == nil && cell.Node == Number {
						minesweeper.Flag(cell.X(), cell.Y())
						numberCellToFlag = cell
					} else if bombCellToNotFlag == nil && cell.Node == Bomb {
						bombCellToNotFlag = cell
					} else if cell.Node == Bomb {
						minesweeper.Flag(cell.X(), cell.Y())
					}
				})

				_, err := minesweeper.Visit(x, y)
				assert.Error(t, err)
				if err != nil {
					assert.IsType(t, &ExplodedError{}, err)
					return // Since the game is over, we also need to stop the test
				}
			}
		}
	}
}

func TestDoNotAutoVisitIfAllNeighboringBombsAreNotFlagged(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.blocks {
		for y, block := range row {
			if block.Node == Number && !block.visited && !block.flagged {
				minesweeper.Visit(x, y)

				// Visit again
				blocks, _ := minesweeper.Visit(x, y)

				assert.Empty(t, blocks)
			}
		}
	}
}

func appendBomb(game *game, x, y int, bombs []*Block, nonBombs []*Block) ([]*Block, []*Block) {
	if x >= 0 && y >= 0 &&
		x < game.Width && y < game.Height {
		if game.blocks[x][y].Node == Bomb {
			return append(bombs, &game.blocks[x][y]), nonBombs
		}
		return bombs, append(nonBombs, &game.blocks[x][y])
	}
	return bombs, nonBombs
}

func count(blocks blocks, width, height, x, y int) (has int) {
	if x >= 0 && y >= 0 &&
		x < width && y < height &&
		blocks[x][y].Node&Bomb == 1 {
		return 1
	}
	return
}

func hasSurroundingTally(blocks blocks, width, height, x, y int) int {
	if x >= 0 && y >= 0 &&
		x < width && y < height {
		switch blocks[x][y].Node {
		case Number:
			return 1
		case Bomb:
			return -1
		default:
			return 0
		}
	}
	return -1
}

func print(game *game) {
	for _, row := range game.blocks {
		fmt.Println()
		for _, block := range row {
			if block.Node == Bomb {
				fmt.Print("* ")
			} else if block.Node == Unknown {
				fmt.Print("  ")
			} else {
				fmt.Printf("%v ", block.Value)
			}
		}
	}
}
