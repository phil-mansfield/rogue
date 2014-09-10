package editor

/* Implements a simple single-line editor, complete with a movable cursor, 
shift-sensitivity, the ability to delete runes, and support for an insert
mode. */

import (
	"fmt"
	"strings"
	"container/list"

	"bitbucket.org/phil-mansfield/roguelike/util/term"
	"bitbucket.org/phil-mansfield/roguelike/util/text"
)

type editor interface {
	SetString(string)
	String() string

	Display(t term.Terminal, x, y int, c term.Color) (string, bool)
}

// Editor contains the internal state required for the line-editor to work.
type Editor struct {
	terminators map[string]bool
	capacity int

	loc int

	ptr *list.Element
	letters *list.List

	terminatorSeen, insert bool

	x, y int
	fg, tg term.Color
}

var _ editor = new(Editor) // typechecking

// New creates a new text editor that can hold a string of up to capacity 
// runes.
func New(capacity int) *Editor { 
	edit := new(Editor)

	if capacity < 0 { panic("LineEditors need positive capacities.") }
	edit.capacity = capacity
		
	edit.terminators = map[string]bool{
		"Enter": true,
		"Esc": true,
	}
	edit.terminatorSeen = false
	edit.insert = false

	edit.SetString("")

	return edit
}

func transportBefore(xs *list.List, elem, target *list.Element) *list.Element{
	v := elem.Value
	xs.Remove(elem)
	return xs.InsertBefore(v, target)
}

func transportAfter(xs *list.List, elem, target *list.Element) *list.Element {
	v := elem.Value
	xs.Remove(elem)
	return xs.InsertAfter(v, target)
}

func (edit *Editor) moveStart() {
	if edit.ptr.Prev() != nil { 
		edit.ptr = transportBefore(edit.letters, edit.ptr, 
			edit.letters.Front())
		edit.loc = 0
	}

}

func (edit *Editor) moveEnd() {
	if edit.ptr.Next() != nil {
		edit.ptr = transportAfter(edit.letters, edit.ptr, edit.letters.Back())
		edit.loc = edit.letters.Len() - 1
	}

}

func (edit *Editor) moveRight() {
	if edit.ptr.Next() != nil {
		edit.ptr = transportAfter(edit.letters, edit.ptr, edit.ptr.Next())
		edit.loc++
	}
}

func (edit *Editor) moveLeft() {
	if edit.ptr.Prev() != nil {
		edit.ptr = transportBefore(edit.letters, edit.ptr, edit.ptr.Prev())
		edit.loc--
	}
}

func (edit *Editor) delRight() {	
	if edit.ptr.Next() != nil { edit.letters.Remove(edit.ptr.Next()) }
}

func (edit *Editor) delLeft() {
	if edit.ptr.Prev() != nil {
		edit.letters.Remove(edit.ptr.Prev())
		edit.loc--
	}
}

func (edit *Editor) addLetter(letter string) {
	if edit.insert { edit.delRight() }
	if edit.capacity + 1 == edit.letters.Len() { return }
	edit.letters.InsertBefore(letter, edit.ptr)
	edit.loc++
}
	
func (edit *Editor) keyPress(key string) {
	if _, ok := edit.terminators[key]; ok {
		edit.terminatorSeen = true
		return
	}
	
	edit.terminatorSeen = false

	if len(key) == 1 {
		edit.addLetter(key) 
		return
	}

	switch key {
	case "Delete":
		edit.delRight()
	case "Backspace":
		edit.delLeft()
	case "Left":
		edit.moveLeft()
	case "Right":
		edit.moveRight()
	case "Up", "Home", "PgUp":
		edit.moveStart()
	case "Down", "End", "PgDn":
		edit.moveEnd()
	case "Insert":
		edit.insert = !edit.insert
	}
}

// SetString sets the editor's default string.
func (edit *Editor) SetString(s string) {
	if edit.capacity < len(s) { 
		panic(fmt.Sprintf("String |%s| = %d too long for capacity %d",
			s, len(s), edit.capacity))
	}

	letterSlice := strings.Split(s, "")
	letterSlice = append(letterSlice, "")
	
	edit.letters = list.New()
	for i := 0; i < len(letterSlice); i ++ {
		edit.letters.PushBack(letterSlice[i])
	}

	edit.ptr = edit.letters.Back()
	if edit.letters.Len() - 1 == edit.capacity {
		edit.loc = edit.capacity - 1
	} else {
		edit.loc = edit.letters.Len() - 1
	}
}

// String returns the string currently inside the editor.
func (edit *Editor) String() string {
	letters := make([]string, 0, edit.letters.Len())

	for elem := edit.letters.Front(); elem != nil; elem = elem.Next() {
		v := elem.Value
		s, ok := v.(string)
		if !ok { panic("What are you doing?") }

		letters = append(letters, s)
	}

	return strings.Join(letters, "")
}

func (edit *Editor) draw(t term.Terminal, x, y int, fg, bg term.Color) {
	s := text.Justify(text.JustLeft, edit.capacity, edit.String())

	fgs := make([]term.Color, len(s))
	bgs := make([]term.Color, len(s))

	loc := edit.loc
	if edit.loc == edit.capacity { loc-- }

	for i := 0; i < edit.capacity; i++ {
		if i == loc {
			fgs[i] = bg
			bgs[i] = fg
		} else {
			fgs[i] = fg
			bgs[i] = bg
		}
	}

	t.PutString(x, y, s)
	t.PutColors(x, y, fgs, term.Foreground)
	t.PutColors(x, y, bgs, term.Background)
	t.Refresh()
}

// Edit places an editor at the specified location in t with the specified
// colors. A possibly nil Output struct and a validity flag are returned. A nil
// Output is returned only if the terminal is closed. The Output struct
// indicates the string written to the editor and whether the user exited it 
// by pressing the Enter key or the Esc key.
func (edit *Editor) Display(t term.Terminal, x, y int, fg, bg term.Color) (*Output, bool) {
	block := t.SaveBlock(x, y, )

	var key string

	for {
		if !t.IsOpen() { return nil, false }

		key = t.NextKey()
		edit.keyPress(key)
		edit.draw(t, x, y, fg, bg)

		if edit.terminatorSeen { break }
	}

	t.RevertBlock(block)
	
	return key, true
}