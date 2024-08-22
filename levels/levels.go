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
const EnemyChar byte = '#'

const (
	NilKey   MoveValue = ""
	UpKey    MoveValue = "k"
	DownKey  MoveValue = "j"
	LeftKey  MoveValue = "h"
	RightKey MoveValue = "l"
)

type ReadyMessageType string

const TickMessage ReadyMessageType = "tick"
const UpdateLevelMessage ReadyMessageType = "updateLevel"
const ResetMessage ReadyMessageType = "reset"

type Level struct {
	TargetLocation Coordinates
	Enemies        []Coordinates
	Number         int
	moveEnemies    func([]Coordinates) []Coordinates
}

func (l *Level) Update() {
	if l.moveEnemies != nil {
		l.Enemies = l.moveEnemies(l.Enemies)
	}

}

var levelsArray = []func() Level{
	func() Level {
		l := Level{
			TargetLocation: Coordinates{50, 25},
			Number:         0,
		}
		return l
	},
	func() Level {
		l := Level{
			TargetLocation: Coordinates{20, 20},
			Number:         1,
			Enemies:        initGrid(0, 0, 3, 3),
		}
		return l
	},
	func() Level {
		l := Level{
			TargetLocation: Coordinates{90, 1},
			Number:         2,
			moveEnemies:    moveEnemies(1, 1, 30),
			Enemies:        initGrid(0, 0, 5, 6),
		}
		return l
	},
}

func GetLevel(i int) (Level, bool) {
	if i >= len(levelsArray) {
		return Level{}, false
	}
	levelFunc := levelsArray[i]
	return levelFunc(), true
}

func moveEnemies(translationX, translationY, ticksPerMove int) func(coord []Coordinates) []Coordinates {
	ticker := 0
	return func(coord []Coordinates) []Coordinates {
		if ticker == ticksPerMove {
			newCoord := make([]Coordinates, 0, len(coord))
			for _, coord_i := range coord {
				coord_i.Y = (coord_i.Y + translationY) % GridSizeY
				coord_i.X = (coord_i.X + translationX) % GridSizeX
				newCoord = append(newCoord, coord_i)
			}
			ticker = 0
			return newCoord
		} else {
			ticker += 1
			return coord
		}
	}
}

func initGrid(xStart, yStart, xSpacing, ySpacing int) []Coordinates {
	grid := make([]Coordinates, 0)
	for i := xStart; i < GridSizeX; i += xSpacing {
		for j := yStart; j < GridSizeY; j += ySpacing {
			grid = append(grid, Coordinates{X: i, Y: j})
		}
	}
	return grid
}
