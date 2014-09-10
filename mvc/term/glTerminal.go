package term

import (
	"fmt"
	"runtime"

	gl "github.com/chsc/gogl/gl21"
	"github.com/go-gl/glfw"

	"bitbucket.org/phil-mansfield/roguelike/util/container/queue"
)

// WARNING: THIS CURRENTLY DOES NOT WORK
// Also, it's completely terrifying to read.

type Layer uint8
const (
	Foreground Layer = iota
	Background
)

// Simple struct for holding color info.
type rgb struct { r, g, b uint8 }
// Gives the lower-case and upper-case versions of a given letter.
type shiftInfo struct { lc, uc string }

// Quick snapshot of keyboard state to be passed through callback functions.
type idInfo struct {
	id int
	shiftPressed, altPressed, ctrlPressed bool
}

type KeyInfo struct {
	Key string
	AltPressed, CtrlPressed bool
}

type GlTerminal struct {
	isInitialized, isClosed, capsOn bool
	// isClosed refers to whether or not t.Close() was called, not
	// whether or not the window itself was closed.
	height, width, runeHeight, runeWidth int
	title string

	runes []rune
	fg, bg []Color

	wasChanged []bool //to speed up rendering

	p []uint8 // pixel array
	ids queue.Queue // all pending key presses represented as KeyInfos.

	f *font
}

//var _ Terminal = &GlTerminal{} // typechecking

var (
	colorMap = map[Color]rgb{ //Remember, these are term.Color keys.
		Black: rgb{0, 0, 0},
		White: rgb{255, 255, 255},
		Gray: rgb{128, 128, 128},
		Red: rgb{255, 0, 0},
		Green: rgb{0, 255, 0},
		Blue: rgb{0, 0, 255},
		Cyan: rgb{0, 255, 255},
		Pink: rgb{255, 0, 255},
		Yellow: rgb{255, 255, 0},
		Purple: rgb{128, 0, 128},
		Brown: rgb{139, 69, 19},
		Orange: rgb{255, 128, 0},
	}

	shiftMap = map[int] shiftInfo {
		'0': shiftInfo{"0", ")"},
		'1': shiftInfo{"1", "!"},
		'2': shiftInfo{"2", "@"},
		'3': shiftInfo{"3", "#"},
		'4': shiftInfo{"4", "$"},
		'5': shiftInfo{"5", "%"},
		'6': shiftInfo{"6", "^"},
		'7': shiftInfo{"7", "&"},
		'8': shiftInfo{"8", "*"},
		'9': shiftInfo{"9", "("},
		glfw.KeyKP0: shiftInfo{"0", "0"},
		glfw.KeyKP1: shiftInfo{"1", "1"},
		glfw.KeyKP2: shiftInfo{"2", "2"},
		glfw.KeyKP3: shiftInfo{"3", "3"},
		glfw.KeyKP4: shiftInfo{"4", "4"},
		glfw.KeyKP5: shiftInfo{"5", "5"},
		glfw.KeyKP6: shiftInfo{"6", "6"},
		glfw.KeyKP7: shiftInfo{"7", "7"},
		glfw.KeyKP8: shiftInfo{"8", "8"},
		glfw.KeyKP9: shiftInfo{"9", "9"},
		'`': shiftInfo{"`", "~"},
		'-': shiftInfo{"-", "_"},
		'=': shiftInfo{"=", "+"},
		',': shiftInfo{",", "<"},
		'.': shiftInfo{".", ">"},
		'/': shiftInfo{"/", "?"},
		'[': shiftInfo{"[", "{"},
		']': shiftInfo{"]", "}"},
		'\\': shiftInfo{"\\", "|"},
		';': shiftInfo{";", ":"},
		'\'': shiftInfo{"'", "\""},
		' ': shiftInfo{" ", " "},
		glfw.KeyEsc: shiftInfo{"Esc", "Esc"},
		glfw.KeyInsert: shiftInfo{"Insert", "Insert"},
		glfw.KeyDel: shiftInfo{"Delete", "Delete"},
		glfw.KeyHome: shiftInfo{"Home", "Home"},
		glfw.KeyEnd: shiftInfo{"End", "End"},
		glfw.KeyBackspace: shiftInfo{"Backspace", "Backspace"},
		glfw.KeyTab: shiftInfo{"Tab", "Tab"},
		glfw.KeyEnter: shiftInfo{"Enter", "Enter"},
		glfw.KeyUp: shiftInfo{"Up", "Up"},
		glfw.KeyDown: shiftInfo{"Down", "Down"},
		glfw.KeyLeft: shiftInfo{"Left", "Left"},
		glfw.KeyRight: shiftInfo{"Right", "Right"},
		glfw.KeyPageup: shiftInfo{"PgUp", "PgUp"},
		glfw.KeyPagedown: shiftInfo{"PgDn", "PgDn"},
	}
)


const ( // This needs to be fixed at some point.
	fontPath = "/Users/phil/code/go/src/bitbucket.org/phil-mansfield/roguelike/util/term/font3.png"
)

func (t *GlTerminal) checkValid() {
	if !t.isInitialized {
		panic("Oops! You forgot to initialize the terminal.")
	} else if t.isClosed {
		panic("Oops! You're trying to work with a closed terminal.")
	}
}

// This function gets a little bit crazy. Don't look at it too closely.
func (t *GlTerminal) Init(title string, height, width int) {
	t.isInitialized = true
	t.isClosed = false
	t.capsOn = false
	
	t.height, t.width = height, width
	
	t.ids = queue.New()
	
	t.f = buildFont(fontPath)
	t.runeHeight, t.runeWidth = t.f.height, t.f.width
	
	t.runes = make([]rune, width * height)
	for i := 0; i < len(t.runes); i++ { t.runes[i] = ' ' }
	t.fg = make([]Color, width * height)
	for i := 0; i < len(t.fg); i++ { t.fg[i] = White }
	t.bg = make([]Color, width * height)
	t.wasChanged = make([]bool, width * height)
	t.p = make([]uint8, t.runeHeight * t.runeWidth * height * width * 3)
	
	runtime.LockOSThread()
	if err := gl.Init(); err != nil { panic(err) }
	if err := glfw.Init(); err != nil { panic(err) }

	glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)

	err := glfw.OpenWindow(width * t.runeWidth, height * t.runeHeight, 
		8, 8, 8, 8, 0, 0, glfw.Windowed)
	if err != nil { panic(err) }

	glfw.SetWindowTitle(title)

	glfw.Enable(glfw.KeyRepeat)

	glfw.SetKeyCallback(func(id, state int) {
		if state == glfw.KeyRelease { return }
		
		sp := glfw.Key(glfw.KeyLshift) == glfw.KeyPress || 
			glfw.Key(glfw.KeyRshift) == glfw.KeyPress
		ap := glfw.Key(glfw.KeyLalt) == glfw.KeyPress || 
			glfw.Key(glfw.KeyRalt) == glfw.KeyPress
		cp := glfw.Key(glfw.KeyLctrl) == glfw.KeyPress || 
			glfw.Key(glfw.KeyRctrl) == glfw.KeyPress
		
		v := idInfo{id: id, shiftPressed: sp, altPressed: ap, ctrlPressed: cp}

		t.ids.Enq(v)
	})

	gl.ClearColor(0, 0, 0, 0)
	gl.Enable(gl.DEPTH_TEST)

	t.Refresh()
}

func (t *GlTerminal) Close() { 
	t.checkValid()

	t.isClosed = true
	glfw.Terminate()
}

// Adds a rune to the terminal's internal pixel array.
func (t *GlTerminal) drawRune(brightness []float32, 
	runeX, runeY int, fg, bg rgb) {
	
	// (Note from the future: I only kinda trsut him on this one.)
	h, w := t.f.height, t.f.width // efficiency hack
	pWidth := t.f.width * t.width

	xStart := runeX * t.f.width
	yStart := runeY * t.f.height

	backR := float32(bg.r)
	backG := float32(bg.g)
	backB := float32(bg.b)

	rDiff := float32(fg.r) - backR
	gDiff := float32(fg.g) - backG
	bDiff := float32(fg.b) - backB

	i := 0
	for y := 0; y < h; y++ {
		for x := xStart; x < w + xStart; x++ {
			pi := (x  + (yStart + h - y - 1) * pWidth) * 3

			t.p[pi]     = uint8(rDiff * brightness[i] + backR)
			t.p[pi + 1] = uint8(gDiff * brightness[i] + backG)
			t.p[pi + 2] = uint8(bDiff * brightness[i] + backB)

			i++
		}
	}
}

func (t *GlTerminal) Refresh() { 
	t.checkValid()

	pHeight, pWidth := t.runeHeight * t.height, t.runeWidth * t.width
	
	runeI := 0
	for runeY := 0; runeY < t.height; runeY++ {
		for runeX := 0; runeX < t.width; runeX++ {

			if t.wasChanged[runeI] {
				t.wasChanged[runeI] = false
				r := t.runes[runeI]
				t.drawRune(t.f.brightness[r], runeX, runeY, 
					colorMap[t.fg[runeI]], colorMap[t.bg[runeI]])
			}

			runeI++
		}
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawPixels(gl.Sizei(pWidth), gl.Sizei(pHeight), 
		gl.RGB, gl.UNSIGNED_BYTE, gl.Pointer(&t.p[0]))
	glfw.SwapBuffers()
}

func (t *GlTerminal) Height() int { 
	t.checkValid()
	return t.height
}
func (t *GlTerminal) Width() int { 
	t.checkValid()
	return t.width
}

// This uses the convention that (0, 0) is the top-left corner.
func (t *GlTerminal) PutRune(x, y int, r rune) {
	t.checkValid()

	y = t.height - 1 - y
	i := y * t.width + x
	if i < 0 || i >= t.height * t.width {
		panic(fmt.Sprintf("(%d, %d) out of terminal bounds (%d, %d)",
			x, y, t.width, t.height))
	}

	if r == t.runes[i] { return }
	t.wasChanged[i] = true
	t.runes[i] = r
}

func (t *GlTerminal) PutColor(x, y int, c Color, layer Layer) {
	t.checkValid()

	y = t.height - 1 - y

	i := y * t.width + x
	if i < 0 || i >= t.height * t.width {
		panic(fmt.Sprintf("(%d, %d) out of terminal bounds (%d, %d)",
			x, y, t.width, t.height))
	}

	switch layer {
	default:
		panic(fmt.Sprintf("%d is not a Layer type. What are you doing?",
			layer))
	case Foreground:
		if c == t.fg[i] { return }
		t.wasChanged[i] = true
		t.fg[i] = c 
	case Background:
		if c == t.bg[i] { return }
		t.wasChanged[i] = true
		t.bg[i] = c
	}
}

func (t *GlTerminal) PutString(x, y int, s string) {
	t.checkValid()

	start := x + t.width * y

	if start < 0 || start + len(s) > t.width * t.height {
		panic(fmt.Sprintf("|%s| = %d at (%d, %d) out of terminal"  +
			"bounds (%d, %d)", s, len(s), x, y, t.width, t.height))
	}

	rs := []rune(s)
	for i := 0 ; i < len(s); i++ { 
		t.PutRune((start + i) % t.width, (start + i) / t.width, rs[i])
	}
}

func (t *GlTerminal) PutColors(x, y int, cs []Color, layer Layer) {
	t.checkValid()

	start := y * t.width + x
	if start < 0 || start + len(cs) > t.height * t.width {
		panic(fmt.Sprintf("slice of length %d at (%d,%d) out of terminal" +
			" bounds (%d,%d)", len(cs), x, y, t.width, t.height))
	}

	for i := 0 ; i < len(cs); i++ { 
		t.PutColor((start + i) % t.width, (start + i) / t.width, cs[i], layer)
	}
}

func simpleShiftRange(info *idInfo) (inRange bool) {
	return info.id >= 'A' && info.id <= 'Z'
}

func mapShiftRange(info *idInfo) (inRange bool) {
	_, ok := shiftMap[info.id]
	return ok
}

func simpleShift(info *idInfo, capsOn bool) string {
	var id int

	if (info.shiftPressed && !capsOn) || (!info.shiftPressed && capsOn) {
		id = info.id
	} else {
		id = info.id + 32
	}

	return fmt.Sprintf("%c", id)
}

func mapShift(info *idInfo) string {
	si, ok := shiftMap[info.id]
	
	if !ok {
		panic(fmt.Sprintf("id %d (%c) is not recognized", info.id, info.id))
	}

	if info.shiftPressed {
		return si.uc
	}
	return si.lc
}

func (t *GlTerminal) Keys() ([]*KeyInfo) { 
	var s string

	t.checkValid()

	is := make([]*KeyInfo, 0)

	for elem, empty := t.ids.Deq(); !empty ; elem, empty = t.ids.Deq() {
		info, ok := elem.(idInfo)
		if !ok { panic("What are you doing?") }
		
		if simpleShiftRange(&info) {		
			s = simpleShift(&info, t.capsOn)
			is = append(is, &KeyInfo{s, info.altPressed, info.ctrlPressed})
		} else if mapShiftRange(&info) {
			s = mapShift(&info)
			is = append(is, &KeyInfo{s, info.altPressed, info.ctrlPressed})
		} else if info.id == glfw.KeyCapslock {
			t.capsOn = !t.capsOn
		}
	}

	return is
}

func (t *GlTerminal) IsOpen() bool { 
	t.checkValid()
	return glfw.WindowParam(glfw.Opened) == gl.TRUE
}