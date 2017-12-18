package minesweeper

type Exploded struct{}

func (Exploded *Exploded) Error() string {
	return "Game over"
}

type GameAlreadyStarted struct{}

func (GameAlreadyStarted *GameAlreadyStarted) String() string {
	return "Game already started. Try setting a new board."
}

func (GameAlreadyStarted *GameAlreadyStarted) Error() string {
	return GameAlreadyStarted.String()
}
