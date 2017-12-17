package minesweeper

type Node uint8

const (
	UNKNOWN = iota
	BOMB
	NUMBER
	FLAGGED
)

type Block struct {
	*Node
}