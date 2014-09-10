package item

import (
	"fmt"
	"math"

	"github.com/phil-mansfield/rogue/error"
)

// Type BufferIndex is an integer type used to index into ListBuffer.
//
// PROGRAMMER NOTE: Since BufferIndex may be changed to different size or
// to a type of unknown signage, all index literals must be constructed from 
// 0, 1, MaxBufferCount, and NilIndex alone.
type BufferIndex int32

const (
	MaxBufferCount = math.MaxInt32 // Largest possible value of ListBuffer.Count.
	NilIndex = -1 // BufferIndex value which is not a valid index into ListBuffer.
	defaultBufferLength = 1 << 8 // The length of an empty ListBuffer.
)

// Type Node is a wrapper around the type Item which allows it to be an element
// in a linked list.
type Node struct {
	Item Item
	Next, Prev BufferIndex
}

// Type ListBuffer is a data structure which represents numerous lists of Items.
type ListBuffer struct {
	FreeHead BufferIndex
	Buffer []Node
	Count BufferIndex
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
	buf.Buffer[len(buf.Buffer) - 1].Next = NilIndex
}

// Singleton creates a singleton list containing only the given Item.
//
// Singleton returns an error if it is passed a nil pointer, if it is passed
// an uninitialized Item, or if the buf is full.
//
// PROGRAMMER NOTE: It is correct to call buf.IsFull() prior to all calls to
// buf.Singleton(), since it is not possible to switch upon the type of error
// to identify whether the error has a recoverable cause.
func (buf *ListBuffer) Singleton(item *Item) (BufferIndex, *error.Error) {
	if buf.IsFull() {
		desc := fmt.Sprintf(
			"buf has reached maximum capacity of %d Items.",
			MaxBufferCount,
		)
		return error.New(error.Value, desc)
	} else if item == nil {
		return error.New(error.Value, "item is nil.")
	} else if item.Type == Uninitialized {
		return error.New(error.Value, "item is uninitialized.")
	}

	return NilIndex, nil
}

// Link connects the Items at indices prev and next so that prev comes before
// next.
//
// Link returns an error if prev or next are not valid indices into buf. 
func (buf *ListBuffer) Link(prev, next BufferIndex) *error.Error {
	inRange, initialized := buf.legalIndex(prev)
	if !inRange {
		desc := fmt.Sprintf(
			"prev, %d, is out of range for IndexBuffer of length %d.",
			prev, len(buf.Buffer),
		)
		return error.New(error.Value, desc)
	} else if !initialized {
		desc := fmt.Sprintf(
			"Item at prev, %d, has the Type value Uninitialized.",
			prev,
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
			"Item at next, %d, has the Type value Uninitialized.",
			next,
		)
		return error.New(error.Value, desc)
	}

	return nil
}

// Unlink removes the Item at idx from its current list.
//
// An error is returned if idx is not a valid index into buf or if it represents
// an uninitialized Item.
//
// PROGRAMMER NOTE: Unlink does not remove the memory imprint of the Item.
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
			"Item at idx, %d, has the Type value Uninitialized.",
			idx,
		)
		return error.New(error.Value, desc)
	}

	return nil
}

// Delete frees the buffer resources associated with the Item at idx.
//
// An error is returned if idx is not a valid index into buf or if it represents
// an uninitialized Item.
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
			"Item at idx, %d, has the Type value Uninitialized.",
			idx,
		)
		return error.New(error.Value, desc)
	}

	return nil
}

// IsFull returns true if no more items can be added to buf.
func (buf *ListBuffer) IsFull() bool {
	return buf.Count >= MaxBufferCount
}

// Get returns the Item stored at idx within buf.
//
// An error is returned if idx is not a valid index into buf or if it represents
// an uninitialized Item.
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
			"Item at idx, %d, has the Type value Uninitialized.",
			idx,
		)
		return nil, error.New(error.Value, desc)
	}

	return nil, nil
}

// legalIndex determines the legality of accessing buf at idx. inRange is true
// if the index is valid and initialized is true if there is an invalid item at
// idx.
func (buf *ListBuffer) legalIndex(idx BufferIndex) (inRange, initialized bool) {
	inRange = idx >= 0 && idx < BufferIndex(len(buf.Buffer))
	if inRange {
		initialized = buf.Buffer[idx].Item.Type != Uninitialized
	} else {
		initialized = true
	}
	return inRange, initialized
}