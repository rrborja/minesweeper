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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExploded_Error(t *testing.T) {
	err := ExplodedError{x: 3, y: 9}
	assert.EqualError(t, err, "Game over at X=3 Y=9")
}

func TestGameAlreadyStarted_Error(t *testing.T) {
	err := GameAlreadyStartedError{}
	assert.EqualError(t, err, "Game already started. Try setting a new board.")
}

func TestUnspecifiedDifficulty_Error(t *testing.T) {
	err := UnspecifiedDifficultyError{}
	assert.EqualError(t, err, "Difficulty was not specified. Use Difficulty(Difficulty) method before calling Play().")
}

func TestUnspecifiedGrid_Error(t *testing.T) {
	err := UnspecifiedGridError{}
	assert.EqualError(t, err, "Grid was not specified. Pass a Grid object with the corresponding coordinates before calling Play().")
}
