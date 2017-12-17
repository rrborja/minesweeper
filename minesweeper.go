package minesweeper

type Node uint8
type Blocks []Block
type Grid struct {width, height int}

const (
	UNKNOWN Node = iota
	BOMB
	NUMBER
	FLAGGED
)

type Block struct {
	Node
}

type Board struct {
	*Grid
	Blocks
}

type game struct {
	Board
}

type Minesweeper interface {
	SetGrid(int, int) error

}

func NewGame(grid ...Grid) Minesweeper {
	game := new(game)
	if len(grid) > 0 {
		game.Grid = &grid[0]
	}
	return game
}

func (game *game) SetGrid(width, height int) error {
	if game.Grid != nil {
		return new(GameAlreadyStarted)
	}
	game.Grid = &Grid{width, height}
	return nil
}

func (block *Block) SetBlock(node Node) {
	block.Node = node
}



