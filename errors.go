package minesweeper

import "fmt"

type Exploded struct{
	location struct {x, y int}
}

func (Exploded *Exploded) Error() string {
	return fmt.Sprintf("Game over at X=%v Y=%v", Exploded.location.x, Exploded.location.y)
}

type GameAlreadyStarted struct{}

func (GameAlreadyStarted *GameAlreadyStarted) String() string {
	return "Game already started. Try setting a new board."
}

func (GameAlreadyStarted *GameAlreadyStarted) Error() string {
	return GameAlreadyStarted.String()
}
