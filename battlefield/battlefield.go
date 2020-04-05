package battlefield

// maxFieldSize is selected to include all english letters.
const maxFieldSize uint = 'z' - 'a' + 1

// Field contains all battlefield data.
type Field struct {
	field      [][]cell
	size       uint
	isSet      bool
	shipsAdded bool

	gameIsOver bool
	shipsAlive int

	state state
}

type cell struct {
	occupied bool
	ship     *ship
	shot     bool
}

type shotResult struct {
	Destroy bool
	Knock   bool
	End     bool
}

type state struct {
	shipCount int
	destroyed int
	knocked   int
	shotCount int
}

// NewField creates new battlefield with provided size.
func NewField(size uint) Field {
	f := make([][]cell, size)
	for i := range f {
		f[i] = make([]cell, size)
	}

	return Field{
		field: f,
		size:  size,
		isSet: true,
	}
}
