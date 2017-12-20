package minesweeper

type Position struct {X, Y int}
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

// Not recommended to call this function until a new update to improve the performance of this method
func (game *game) HintLocations() []Position {
	hintPlacements := make([]Position, 0) // TODO: Improve this performance

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				hintPlacements = append(hintPlacements, Position{x, y})
			}
		}
	}

	return hintPlacements
}

func (game *game) History() History {
	return nil
}

func (game *game) LastAction() Record {
	return Record{}
}