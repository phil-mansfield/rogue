package item

import (
	"fmt"
	"math"

	"github.com/phil-mansfield/rogue/error"
)

// Type BufferIndex is an integer type used to index into ListBuffer.
//
// Since BufferIndex may be changed to different size or to a type of unknown
// signage, all BufferIndex literals must be constructed from 0, 1,
// MaxBufferCount, and NilIndex alone.
type BufferIndex int32

const (
	MaxBufferCount      = math.MaxInt32 // Largest possible value of ListBuffer.Count.
	NilIndex            = -1            // Sentinel ListBuffer index value.
	defaultBufferLength = 1 << 8        // The length of an empty ListBuffer.
)

// Type Node is a wrapper around the type Item which allows it to be an element
// in a linked list.
//
// Next and Prev reference the indices of other Items within the same instance
// of ListBuffer.
type Node struct {
	Item       Item
	Next, Prev BufferIndex
}

// Type ListBuffer is a data structure which represents numerous lists of items.
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
// Singleton returns an error if it is passed a nil pointer, if it is passed
// an uninitialized item, or if the buf is full.
//
// It is correct to call buf.IsFull() prior to all calls to
// buf.Singleton(), since it is not possible to switch upon the type of error
// to identify whether the error has a recoverable cause.
func (buf *ListBuffer) Singleton(item *Item) (BufferIndex, *error.Error) {
	if buf.IsFull() {
		desc := fmt.Sprintf(
			"buf has reached maximum capacity of %d Items.",
			MaxBufferCount,
		)
		return NilIndex, error.New(error.Value, desc)
	} else if item == nil {
		return NilIndex, error.New(error.Value, "item is nil.")
	} else if item.Type == Uninitialized {
		return NilIndex, error.New(error.Value, "item is uninitialized.")
	}

	return buf.internalSingleton(item), nil
}

func (buf *ListBuffer) internalSingleton(item *Item) BufferIndex {
	if buf.Count == len(buf.Buffer) {
		buf.Buffer = append(buf.Buffer, *item)
		return len(buf.Buffer) - 1
	}

	idx := buf.FreeHead
	buf.Buffer[idx].Item = *item
	buf.FreeHead = buf.Buffer[idx].Next

	buf.internalUnlink(idx)
	buf.Count++

	return idx
}

// Link connects the items at indices prev and next so that the item at prev
// comes before the item at next.
//
// Link returns an error if prev or next are not valid indices into buf or if
// the linking would break a pre-existing list.
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

	inRange, initialized = buf.legalIndex(prev)
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
// An error is returned if idx is not a valid index into buferf or if it
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
	buf.internalUnlink()
	if buf.FreeHead != NilIndex {
		err := buf.internalLink(idx, buf.FreeHead)
	}	
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
func (buf *ListBuffer) Get(idx BufferIndex) (*Item, *error.Error) {
	inRange, initialized := buf.legalIndex(idx)
	if !inRange {
		desc := fmt.Sprintf(
			"idx, %d, is out of range for IndexBuffer of length %d.",
			idx, len(buf.Buffer),
		)
		return nil, error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at idx, %d, has the Type value Uninitialized.", idx,
		)
		return nil, error.New(error.Value, desc)
	}

	return &buf.Buffer[idx].Item, nil
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
