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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExploded_Error(t *testing.T) {
	err := Exploded{struct{ x, y int }{x: 3, y: 9}}
	assert.EqualError(t, err, "Game over at X=3 Y=9")
}

func TestGameAlreadyStarted_Error(t *testing.T) {
	err := GameAlreadyStarted{}
	assert.EqualError(t, err, "Game already started. Try setting a new board.")
}

func TestUnspecifiedDifficulty_Error(t *testing.T) {
	err := UnspecifiedDifficulty{}
	assert.EqualError(t, err, "Difficulty was not specified. Use Difficulty(Difficulty) method before calling Play().")
}

func TestUnspecifiedGrid_Error(t *testing.T) {
	err := UnspecifiedGrid{}
	assert.EqualError(t, err, "Grid was not specified. Pass a Grid object with the corresponding coordinates before calling Play().")
}
