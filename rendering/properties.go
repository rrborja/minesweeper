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

package rendering

// Position is used to interface the cell's xy-coordinates used for this package
type Position interface {
	// X returns the x-coordinate of the cell in the grid
	X() int

	// Y returns the y-coordinate of the cell in the grid
	Y() int
}

// Tracker is used to interface the instance of the Minesweeper game to retrieve
// certain information such as the location of all mines and the location of all
// non-zero warning values
type Tracker interface {
	// BombLocations returns all the location of all mines in the grid
	BombLocations() []Position

	// HintLocations returns all the warning numbers in the grid
	HintLocations() []Position
}
