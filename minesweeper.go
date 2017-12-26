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

const (
	// Unknown is the type of the value contained in a minesweeper
	// cell. It has no value which means no mines are neighbored in
	// the cell.
	Unknown Node = 1 << iota

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

// Minesweeper is the main point of consumption and manipulation of the Minesweeper's
// game state. Methods provided by this interface are common use cases to solve
// a minesweeper game.
//
// Any instance derived by this interface is compatible for type casting to the
// rendering.Tracker and visited.StoryTeller interfaces.
type Minesweeper interface {
	SetGrid(int, int) error

	SetDifficulty(Difficulty) error

	Play() error

	Flag(int, int)

	Visit(int, int) ([]Block, error)
}

// NewGame creates a separate minesweeper instance. Unlike minesweeper.New,
// this function creates a non-singleton instance. Functions of this package
// such as Visit will become the methods of this instance.
func NewGame(grid ...Grid) (Minesweeper, Event) {
	game := new(game)

	if len(grid) > 0 {
		game.SetGrid(grid[0].Width, grid[0].Height)
	}

	game.Event = make(chan eventType, 1)

	return game, game.Event
}

// New creates a new minesweeper environment. Note that this only creates the minesweeper
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
func New(grid ...Grid) Event {
	minesweeper, mainEvent := NewGame(grid...)
	singleton = minesweeper
	return mainEvent
}

// SetGrid sets the board's size. You can't change the board's size once
// the Play() method has been called, otherwise, a GameAlreadyStartedError
// will return.
//
// Nothing will return if the setting the board's size is successful.
//
// Whenever the game ends, this instance will be garbage collected. Calling
// any exported methods when the game ended may result in nil pointer
// dereference. Creating a new setup of the game is ideal.
func SetGrid(width int, height int) error {
	return singleton.SetGrid(width, height)
}

// SetDifficulty sets the difficulty of the game as a basis of the number of mines. An
// error will return if this method is being called when a game is already
// being played or better yet, the Play() method has already been called.
func SetDifficulty(difficulty Difficulty) error {
	return singleton.SetDifficulty(difficulty)
}

// Play allows the game to setup all the mines in place randomly. The
// placement is non-deterministic since the implementation uses the
// "crypto/rand" package.
//
// An error will return when this method is called twice or more.
//
// More importantly, Grid size and Difficulty must be specified, otherwise,
// you will encounter an UnspecifiedGridError and UnspecifiedDifficultyError,
// respectively.
func Play() error {
	return singleton.Play()
}

// Flag marks the cell, according to the coordinates supplied in the
// method argument, as flagged. When a particular cell is flagged, the
// cell in question will prevent from being handled by the game when
// the Visit(int, int) method with the same coordinate of the cell in
// question is called.
func Flag(x int, y int) {
	singleton.Flag(x, y)
}

// Visit visits a particular cell according to the xy-coordinates of the argument
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
func Visit(x int, y int) ([]Block, error) {
	return singleton.Visit(x, y)
}
