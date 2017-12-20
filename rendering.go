package minesweeper

type Position struct {x, y int}
type Message string
type History []Record

type Record struct {
	Position
	Message
}

type RuntimeGameProperties interface {

	BombLocations() Blocks

	HintLocations() Blocks

	History() History

	LastAction() Record
}

func (game *game) BombLocations() Blocks {
	return nil
}

func (game *game) HintLocations() Blocks {
	return nil
}

func (game *game) History() History {
	return nil
}

func (game *game) LastAction() Record {
	return Record{}
}