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

import "fmt"

type Exploded struct {
	location struct{ x, y int }
}

func (Exploded *Exploded) Error() string {
	return fmt.Sprintf("Game over at X=%v Y=%v", Exploded.location.x, Exploded.location.y)
}

type GameAlreadyStarted struct{}

func (GameAlreadyStarted *GameAlreadyStarted) Error() string {
	return "Game already started. Try setting a new board."
}
