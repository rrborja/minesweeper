package minesweeper

type Node uint8
type Blocks [][]Block
type Grid struct {width, height uint}

type Difficulty uint8

const (
	UNKNOWN Node = iota
	BOMB
	NUMBER
	FLAGGED
)

const (
	NOTSET Difficulty = iota
	EASY
	MEDIUM
	HARD
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
	Difficulty
}

type Minesweeper interface {
	SetGrid(uint, uint) error

	SetDifficulty(Difficulty)

	Play() error

	Flag(uint, uint) error
}

func NewGame(grid ...Grid) Minesweeper {
	game := new(game)
	if len(grid) > 0 {
		game.SetGrid(grid[0].width, grid[0].height)
	}
	return game
}

func (game *game) SetGrid(width, height uint) error {
	if game.Grid != nil {
		return new(GameAlreadyStarted)
	}
	game.Grid = &Grid{width, height}
	createBoard(game)
	return nil
}

func (game *game) Flag(x, y uint) error {
	game.Blocks[x][y].Node = FLAGGED
	return nil
}

func (game *game) SetDifficulty(difficulty Difficulty) {
	game.Difficulty = difficulty
}

func (game *game) Play() error {
	return nil
}

func (block *Block) SetBlock(node Node) {
	block.Node = node
}

func createBoard(game *game) {
	game.Blocks = make([][]Block, game.height)
	for y := range game.Blocks {
		game.Blocks[y] = make([]Block, game.width)
	}
}