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

import "fmt"

// ExplodedError is the error type used to handle a situation when a mine is visited
type ExplodedError struct {
	location struct{ x, y int }
}

func (Exploded ExplodedError) Error() string {
	return fmt.Sprintf("Game over at X=%v Y=%v", Exploded.location.x, Exploded.location.y)
}

// GameAlreadyStartedError is the error type used to handle errors when attempting
// to restart the game, changing the Grid size when the game is started, or changing
// the difficulty while the game is started.
type GameAlreadyStartedError struct{}

func (GameAlreadyStarted GameAlreadyStartedError) Error() string {
	return "Game already started. Try setting a new board."
}

// UnspecifiedDifficultyError is the error type used to handle errors when the Play()
// method is called but the Difficulty is not set in the game.
type UnspecifiedDifficultyError struct{}

func (UnspecifiedDifficulty UnspecifiedDifficultyError) Error() string {
	return "Difficulty was not specified. Use Difficulty(Difficulty) method before calling Play()."
}

// UnspecifiedGridError is the error type used to handle errors when the Play() method
// is called but the Grid size is not set in the game.
type UnspecifiedGridError struct{}

func (UnspecifiedGrid UnspecifiedGridError) Error() string {
	return "Grid was not specified. Pass a Grid object with the corresponding coordinates before calling Play()."
}
