package minesweeper

type Position struct {x, y int}
type Message string
type History []Record

type Record struct {
	Position
	Message
}

type RuntimeGameProperties interface {

	BombLocations() []Position

	HintLocations() []Position

	History() History

	LastAction() Record
}

func (game *game) BombLocations() []Position {
	return nil
}

func (game *game) HintLocations() []Position {
	return nil
}

func (game *game) History() History {
	return nil
}

func (game *game) LastAction() Record {
	return Record{}
}