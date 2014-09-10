package term

import (
	"os"
	"fmt"
	"strings"

	"github.com/jabb/gocurse/curses"
)
type CursesTerminal struct {
	screen *curses.Window
	height, width int

	colors []Color
	runes []rune
	changed []bool
}

var (
	cursesColorMap = make([]int, int(colorNum))
	cursesColorMap[int(Black)] = curses.COLOR_BLACK
	cursesColorMap[int(White)] = curses.COLOR_WHITE
	cursesColorMap[int(Gray)] = curses.COLOR_WHITE
	cursesColorMap[int(Red)] = curses.COLOR_RED
	cursesColorMap[int(Green)] = curses.COLOR_GREEN
	cursesColorMap[int(Blue)] = curses.COLOR_BLUE
	cursesColorMap[int(Cyan)] = curses.COLOR_CYAN
	cursesColorMap[int(Pink)] = curses.COLOR_MAGENTA
	cursesColorMap[int(Yellow)] = curses.COLOR_YELLOW
	cursesColorMap[int(Purple)] = curses.COLOR_MAGNETA
	cursesColorMap[int(Brown)] = curses.COLOR_YELLOW
	cursesColorMap[int(Orange)] = curses.COLOR_RED
)

// init creates a CursesTerminal instance with the given height and width.
func (term *CursesTerminal) init(height, width int) error  {
	if height <= 0 || width <= 0 {
		
	} else if term.screen, err = curses.Initscr(); err != nil {
		return err
	}

	term.width = width
	term.height = height

	term.colors = make([]Color, term.Width() * term.Height())
	term.runes = make([]rune, term.Width() * term.Height())
	term.changed = make([]bool, term.Width() * term.Height())
	for i := 0; i < len(term.colors); i++ {
		term.colors[i] = Black
		term.runes[i] = '*'
		term.changed[i] = true
	}

	curses.Noecho()
	curses.Start_color()

	curses.Init_pair(curses.COLOR_WHITE, 
		curses.COLOR_WHITE, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_RED, 
		curses.COLOR_RED, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_GREEN, 
		curses.COLOR_GREEN, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_YELLOW, 
		curses.COLOR_YELLOW, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_BLUE, 
		curses.COLOR_BLUE, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_MAGENTA, 
		curses.COLOR_MAGENTA, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_CYAN, 
		curses.COLOR_CYAN, curses.COLOR_BLACK)
	curses.Init_pair(curses.COLOR_WHITE, 
		curses.COLOR_WHITE, curses.COLOR_BLACK)
}

func (term *CursesTerminal) Close() { 
	curses.Endwin()
}

func (term *CursesTerminal) IsOpen() bool  {
	return true
}

func (term *CursesTerminal) Refresh() { 
	for y := 0; y < term.Height(); y++ {
		for x := 0; x < term.Width(); x++ {
			i := x + y * term.Width()			

			if !term.changed[i] { continue }

			flag := curses.Color_pair(colorFlagMap[term.colors[i]])
			term.screen.Addch(x, y, term.runes[i], flag)
			term.changed[i] = false
		}
	}
	term.screen.Refresh()
}

func (term *CursesTerminal) Height() int { return term.height }
func (term *CursesTerminal) Width() int { return term.width }

func (term *CursesTerminal) PutRune(x, y int, r rune, color Color) { 
	i := term.Width() * y + x

	if x < 0 || y < 0 || x >= term.Width() || y >= term.Height() {
		term.Close()
		panic(fmt.Sprintf("Cell (%d %d) is out of window bounds (%d %d).",
			x, y, term.Width(), term.Height()))
	}

	if color != term.colors[i] || r != term.runes[i] {
		term.changed[i] = true
		term.colors[i] = color
		term.runes[i] = r
	}
}

func (term *CursesTerminal) PutString(x, y int, str string, color Color) { 
	lines := strings.Split(str, "\n")
	
	if y + len(lines) >= term.Height() {
		term.Close()
		panic(fmt.Sprintf("str %d + '%s' too tall for terminal height %d.",
			y, str, term.Height()))
	}

	for row, line := range lines {
		if x + len(line) >= term.Width() {
			panic(fmt.Sprintf("Line %d + '%s' too long for terminal width %d.",
				x, line, term.Width()))
		}

		for col, r := range line {
			term.PutRune(x + col, y + row, rune(r), color)
		}
	}
}

func (term *CursesTerminal) NextKey() string {
	for {
		ch := term.screen.Getch()

		if ' ' <= ch && ch <= '~' {
			return string(ch)
		}

		switch ch {
		case curses.KEY_ENTER:
			return "Enter"
		case curses.KEY_LEFT:
			return "Left"
		case curses.KEY_RIGHT:
			return "Right"
		case curses.KEY_UP:
			return "Up"
		case curses.KEY_DOWN:
			return "Down"
		case curses.KEY_BACKSPACE:
			return "Backspace"
		case 127:
			return "Delete"
		case 27:
			return "Esc"
		case 9:
			return "Tab"
		case curses.KEY_NPAGE:
			return "PgDn"
		case curses.KEY_PPAGE:
			return "PgUp"
		case curses.KEY_IC:
			return "Insert"
		}
		fmt.Fprintf(os.Stderr, "Unknown key: %d, '%c'\n", ch, ch)
	}
}