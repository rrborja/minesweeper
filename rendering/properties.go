package rendering

type Position struct{ X, Y int }
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
