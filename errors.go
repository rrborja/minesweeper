package minesweeper

type Exploded struct {
	error string
}

func (exploded *Exploded) Error() string {
	return exploded.error
}