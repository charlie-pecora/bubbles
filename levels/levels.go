package levels

type Coordinates struct {
	X int
	Y int
}

type MoveValue string

const GridSizeX int = 100
const GridSizeY int = 30
const UserChar byte = 'A'
const EmptyChar byte = '.'
const TargetChar byte = '@'
const WallChar byte = '#'

const (
	NilKey   MoveValue = ""
	UpKey    MoveValue = "k"
	DownKey  MoveValue = "j"
	LeftKey  MoveValue = "h"
	RightKey MoveValue = "l"
)

type ReadyMessageType string

const TickMessage ReadyMessageType = "tick"
const ResetMessage ReadyMessageType = "reset"

type Level struct {
	TargetLocation Coordinates
	Walls          []Coordinates
}

var LevelsArray = []Level{
	{
		TargetLocation: Coordinates{20, 5},
	},
	{
		TargetLocation: Coordinates{5, 20},
	},
}
