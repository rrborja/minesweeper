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
	bombPlacements := make([]Position, int(float32(game.Height * game.Width) * game.difficultyMultiplier))

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				bombPlacements = append(bombPlacements, Position{x, y})
			}
		}
	}

	return bombPlacements
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