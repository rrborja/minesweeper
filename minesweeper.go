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
	Grid
	Blocks
}

type game struct {
	Board
}

type Minesweeper interface {
	SetGrid(int, int) *game

}

func NewGame() Minesweeper {
	return new(game)
}

func (game *game) SetGrid(width, height int) *game {
	game.Grid = Grid{width, height}
	return game
}

func (block *Block) SetBlock(node Node) {
	block.Node = node
}

