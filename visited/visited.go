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

package visited

// Action is the type of the cell used to record the player's move
type Action uint8

const (
	// Unknown is the value of the cell with no value
	Unknown Action = 1 << iota

	// Number is the value of the cell with a warning value
	Number

	// Bomb is the value of the cell with the mine
	Bomb
)

// History contains the information of all player's movements in a linked list
type History struct {
	Record
	*History
}

// Record contains the information of the player's movement such as the position
// of the cell visited and visited cell's type
type Record struct {
	Position
	Action
}

// Position is used to interface the cell's xy-coordinates used for this package
type Position interface {
	// X returns the x-coordinate of the cell in the grid
	X() int

	// Y returns the y-coordinate of the cell in the grid
	Y() int
}

// StoryTeller is used to interface the instance of the Minesweeper game to retrieve
// certain information such as the history of all player's move and the player's
// recent action
type StoryTeller interface {
	// History returns the list of player's move. The iteration of the returned value
	// is not the same as iterating an array because the returned value is in the
	// implementation of a linked-list.
	History() *History

	// LastAction returns the player's recent move
	LastAction() Record
}
