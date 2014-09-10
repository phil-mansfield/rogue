/* package term provides a simplified interface for creating and manipulating
a terminal. package term also supplies several implementations of this
interface. */
package term

// TODO: come up with set that is consistent across keyboard setups.

// Special keys are: "Enter", "Left", "Up", "Down", "Right", "Esc", 
// "Backspace", "Delete", "Tab", "PgDn", "PgUp", and "Insert".

// Color represents a color that can be drawn to the screen.
type Color uint8
const (
	Black Color = iota
	White
	Gray
	Red
	Green
	Blue
	Cyan
	Pink
	Yellow
	Purple
	Brown
	Orange
	colorNum
)

// TerminalType is a flag representing the underlying implementation of the
// terminal.
//
// Implementation note: Currently the only cross-platform terminal 
// implementation is Curses. Although much time has gone into developing Gl,
// there are currently no plans to ever deploy it.
type Type uint8
const (
	Curses Type = iota
	Gl
)

// Block stores the information required to revert a rectangular portion of
// the screen
type Block struct {
	x, y, height, width int
	rs []rune
	cs []Color
}

// Terminal is the basic terminal interface central to this package.
type Terminal interface {
	Close()
	IsOpen() bool

	NextKey() string, error

	PutForeground(x, y int, rs []rune, cs []Color) error
	PutBackground(x, y int, cs []Color) error

	Refresh()

	SaveBlock(x, y, length, height) *Block, error
	RevertBlock(block *Block)

	Height() int
	Width() int
}

// New creates a new terminal of the specified type.
func New(t Type, width, height int) Terminal {
	switch tt {
	case Curses:
		panic("term.Curses is currently non-functional.")
	case Gl:
		panic("term.Gl is currently non-functional.")
	}
	panic("Unrecognized TerminalType.")
}