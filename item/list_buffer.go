package item

import (
	"fmt"
	"math"

	"github.com/phil-mansfield/rogue/error"
)

// BufferIndex is an integer type used to index into ListBuffer.
//
// Since BufferIndex may be changed to different size or to a type of unknown
// signage, all BufferIndex literals must be constructed from 0, 1,
// MaxBufferCount, and NilIndex alone.
type BufferIndex int16

// Note that the usage of "Count" and "Length" in constant names *is* actually
// consistent.

const (
	// MaxBufferCount is the largest possible value of ListBuffer.Count.
	MaxBufferCount      = math.MaxInt16
	// NilIndex is a sentinel ListBuffer index value. It is analogous to a a
	// nil pointer.
	NilIndex            = -1    
	// defaultBufferLength is the length of an empty ListBuffer.        
	defaultBufferLength = 1 << 8
)

// Node is a wrapper around the type Item which allows it to be an element
// in a linked list.
//
// Next and Prev reference the indices of other Items within the same instance
// of ListBuffer.
type Node struct {
	Item       Item
	Next, Prev BufferIndex
}

// ListBuffer is a data structure which represents numerous lists of items.
type ListBuffer struct {
	FreeHead BufferIndex
	Buffer   []Node
	Count    BufferIndex
}

// New creates a new ListBuffer instance.
func New() *ListBuffer {
	buf := new(ListBuffer)
	buf.Init()
	return buf
}

// Init initializes a blank ListBuffer instance.
func (buf *ListBuffer) Init() {
	buf.Buffer = make([]Node, defaultBufferLength)
	buf.FreeHead = 0
	buf.Count = 0

	for i := 0; i < len(buf.Buffer); i++ {
		buf.Buffer[i].Item.Clear()
		buf.Buffer[i].Prev = BufferIndex(i - 1)
		buf.Buffer[i].Next = BufferIndex(i + 1)
	}

	buf.Buffer[0].Prev = NilIndex
	buf.Buffer[len(buf.Buffer)-1].Next = NilIndex
}

// Singleton creates a singleton list containing only the given item.
//
// Singleton returns an error if it is passed
// an uninitialized item, or if the buf is full.
//
// It is correct to call buf.IsFull() prior to all calls to
// buf.Singleton(), since it is not possible to switch upon the type of error
// to identify whether the error has a recoverable cause.
func (buf *ListBuffer) Singleton(item Item) (BufferIndex, *error.Error) {
	if buf.IsFull() {
		desc := fmt.Sprintf(
			"buf has reached maximum capacity of %d Items.",
			MaxBufferCount,
		)
		return NilIndex, error.New(error.Value, desc)
	} else if item.Type == Uninitialized {
		return NilIndex, error.New(error.Value, "item is uninitialized.")
	}

	return buf.internalSingleton(item), nil
}

func (buf *ListBuffer) internalSingleton(item Item) BufferIndex {
	if buf.Count == BufferIndex(len(buf.Buffer)) {
		buf.Buffer = append(buf.Buffer, Node{item, NilIndex, NilIndex})
		buf.Count++
		return BufferIndex(len(buf.Buffer) - 1)
	}

	idx := buf.FreeHead
	buf.Buffer[idx].Item = item
	buf.FreeHead = buf.Buffer[idx].Next

	buf.internalUnlink(idx)
	buf.Count++

	return idx
}

// Link connects the items at indices prev and next so that the item at prev
// comes before the item at next.
//
// Link returns an error if prev or next are not valid indices into buf or if
// the linking would break a pre-existing list or if one of the indices accesses
// a .
func (buf *ListBuffer) Link(prev, next BufferIndex) *error.Error {
	// If your functions don't have 50 lines of error handling for two lines
	// of state altering-code, you aren't cautious enough.

	inRange, initialized := buf.legalIndex(prev)
	if !inRange {
		desc := fmt.Sprintf(
			"prev, %d, is out of range for IndexBuffer of length %d.",
			prev, len(buf.Buffer),
		)
		return error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at prev, %d, has the Type value Uninitialized.", prev,
		)
		return error.New(error.Value, desc)
	}

	inRange, initialized = buf.legalIndex(next)
	if !inRange {
		desc := fmt.Sprintf(
			"next, %d, is out of range for IndexBuffer of length %d.",
			next, len(buf.Buffer),
		)
		return error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at next, %d, has the Type value Uninitialized.", next,
		)
		return error.New(error.Value, desc)
	}

	if buf.Buffer[prev].Next != NilIndex {
		desc := fmt.Sprintf(
			"ItemNode at prev, %d, is already linked to Next ItemNode at %d.",
			prev, buf.Buffer[prev].Next,
		)
		return error.New(error.Value, desc)
	}

	if buf.Buffer[next].Prev != NilIndex {
		desc := fmt.Sprintf(
			"ItemNode at next, %d, is already linked to Prev ItemNode at %d.",
			next, buf.Buffer[next].Prev,
		)
		return error.New(error.Value, desc)
	}

	buf.internalLink(prev, next)
	return nil
}

func (buf *ListBuffer) internalLink(prev, next BufferIndex) {
	buf.Buffer[next].Prev = prev
	buf.Buffer[prev].Next = next
}

// Unlink removes the item at the given index from its current list.
//
// An error is returned if idx is not a valid index into the buffer or if
// it represents an uninitialized item.
func (buf *ListBuffer) Unlink(idx BufferIndex) *error.Error {
	inRange, initialized := buf.legalIndex(idx)
	if !inRange {
		desc := fmt.Sprintf(
			"idx, %d, is out of range for IndexBuffer of length %d.",
			idx, len(buf.Buffer),
		)
		return error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at idx, %d, has the Type value Uninitialized.", idx,
		)
		return error.New(error.Value, desc)
	}

	buf.internalUnlink(idx)
	return nil
}

func (buf *ListBuffer) internalUnlink(idx BufferIndex) {
	next := buf.Buffer[idx].Next
	prev := buf.Buffer[idx].Prev

	if prev != NilIndex {
		buf.Buffer[prev].Next = next
	}
	if next != NilIndex {
		buf.Buffer[next].Prev = prev
	}

	buf.Buffer[idx].Next = NilIndex
	buf.Buffer[idx].Prev = NilIndex
}

// Delete frees the buffer resources associated with the item at the given
// index.
//
// An error is returned if idx is not a valid index into buffer or if it
// represents an uninitialized item.
func (buf *ListBuffer) Delete(idx BufferIndex) *error.Error {
	inRange, initialized := buf.legalIndex(idx)
	if !inRange {
		desc := fmt.Sprintf(
			"idx, %d, is out of range for IndexBuffer of length %d.",
			idx, len(buf.Buffer),
		)
		return error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at idx, %d, has the Type value Uninitialized.", idx,
		)
		return error.New(error.Value, desc)
	}

	buf.internalDelete(idx)
	return nil
}

func (buf *ListBuffer) internalDelete(idx BufferIndex) {
	buf.internalUnlink(idx)
	if buf.FreeHead != NilIndex {
		buf.internalLink(idx, buf.FreeHead)
	}

	node := &buf.Buffer[idx]
	node.Item.Clear()
	node.Next = buf.FreeHead
	node.Prev = NilIndex

	buf.FreeHead = idx

	buf.Count--
}

// IsFull returns true if no more items can be added to the buffer.
func (buf *ListBuffer) IsFull() bool {
	return buf.Count >= MaxBufferCount
}

// Get returns the item stored at the given index within the buffer.
//
// An error is returned if idx is not a valid index into the buffer or if it
// represents an uninitialized item.
func (buf *ListBuffer) Get(idx BufferIndex) (Item, *error.Error) {
	inRange, initialized := buf.legalIndex(idx)
	if !inRange {
		desc := fmt.Sprintf(
			"idx, %d, is out of range for IndexBuffer of length %d.",
			idx, len(buf.Buffer),
		)
		return Item{}, error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at idx, %d, has the Type value Uninitialized.", idx,
		)
		return Item{}, error.New(error.Value, desc)
	}

	return buf.Buffer[idx].Item, nil
}

// Set updates the item stored at the given index within the buffer.
//
// An error is returned if idx is not a valid index into the buffer or if it
// represents an uninitialized item.
func (buf *ListBuffer) Set(idx BufferIndex, item Item) (*error.Error) {
	inRange, initialized := buf.legalIndex(idx)
	if !inRange {
		desc := fmt.Sprintf(
			"idx, %d, is out of range for IndexBuffer of length %d.",
			idx, len(buf.Buffer),
		)
		return error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at idx, %d, has the Type value Uninitialized.", idx,
		)
		return error.New(error.Value, desc)
	}

	buf.Buffer[idx].Item = item
	return nil
}

// legalIndex determines the legality of accessing the buffer at idx. inRange
// is true if the index is valid and initialized is true if there is an valid
// item at idx.
func (buf *ListBuffer) legalIndex(idx BufferIndex) (inRange, initialized bool) {
	inRange = idx >= 0 && idx < BufferIndex(len(buf.Buffer))
	if inRange {
		initialized = buf.Buffer[idx].Item.Type != Uninitialized
	} else {
		initialized = true
	}
	return inRange, initialized
}

// Check performs various consistency checks on the buffer and returns an error
// indicating which the first failed check. If all checks pass, nil is returned.
func (buf *ListBuffer) Check() *error.Error {

	// Check that the buffer isn't too long.
	if len(buf.Buffer) > MaxBufferCount {
		desc := fmt.Sprintf(
			"Buffer length of %d is larger than max of %d.",
			len(buf.Buffer), MaxBufferCount,
		)
		return error.New(error.Sanity, desc)
	}

	// Check that all items are valid.
	for i := 0; i < len(buf.Buffer); i++ {
		if err := buf.Buffer[i].Item.Check(); err != nil { return err }
	}

	// Check that the item count is correct.
	count := 0
	for i := 0; i < len(buf.Buffer); i++ {
		if buf.Buffer[i].Item.Type != Uninitialized {
			count++
		}
	}

	if BufferIndex(count) != buf.Count {
		desc := fmt.Sprintf(
			"buf.Count = %d, but there are %d items in buffer.",
			buf.Count, count,
		)
		return error.New(error.Sanity, desc)
	}

	// Check all Prev indices.
	for i := 0; i < len(buf.Buffer); i++ {
		prev := buf.Buffer[i].Prev
		if prev != NilIndex && buf.Buffer[prev].Next != BufferIndex(i) {
			desc := fmt.Sprintf(
				"Prev index of item %d is %d, but Next index of item %d is %d.",
				i, prev, prev, buf.Buffer[prev].Next,
			)
			return error.New(error.Sanity, desc)
		}
	}

	// Check all Next indices.
	for i := 0; i < len(buf.Buffer); i++ {
		next := buf.Buffer[i].Next
		if next != NilIndex && buf.Buffer[next].Prev != BufferIndex(i) {
			desc := fmt.Sprintf(
				"Next index of item %d is %d, but Prev index of item %d is %d.",
				i, next, next, buf.Buffer[next].Prev,
			)
			return error.New(error.Sanity, desc)
		}
	}

	// Check for cycles.
	checkBuffer := make([]bool, len(buf.Buffer))
	for i := 0; i < len(buf.Buffer); i++ {
		if !checkBuffer[i] && buf.hasCycle(BufferIndex(i), checkBuffer) {
			desc := fmt.Sprintf("List with head at index %d contains cycle.", i)
			return error.New(error.Sanity, desc)
		}
	}

	return nil
}

// hasCycle returns true if there is a cycle after in the list following the
// given index and that cycle has not already been detected.
func (buf *ListBuffer) hasCycle(idx BufferIndex, checkBuffer []bool) bool {
	checkBuffer[idx] = true

	tortise := idx
	hare := idx

	for hare != NilIndex {
		hare = buf.Incr(buf.Incr(hare))
		tortise = buf.Incr(tortise)

		if hare != NilIndex { checkBuffer[hare] = true }
		if tortise != NilIndex { checkBuffer[tortise] = true }

		if hare == tortise && hare != NilIndex {
			return true
		}
	}

	return false
}

// Incr returns the index of item which follows idx in the current list.
func (buf *ListBuffer) Incr(idx BufferIndex) BufferIndex {
	if idx == NilIndex {
		return NilIndex
	} 
	return buf.Buffer[idx].Next
}

func (buf *ListBuffer) Decr(idx BufferIndex) BufferIndex {
	if idx == NilIndex {
		return NilIndex
	}
	return buf.Buffer[idx].Prev
}
