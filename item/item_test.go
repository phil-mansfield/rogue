package item

import (
	"testing"
)

// Note that occasionally when a garbage index needs to be chosen, certain
// conservative assumptions are made about the value of defaultBufferSize.

func TestClear(t *testing.T) {
	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}
	refItem := Item{0, Uninitialized, [6]int8{0, 0, 0, 0, 0, 0}}
	item.Clear()

	if item != refItem {
		t.Errorf("Item.Clear() results in %v.", item)
	}
}

func TestNew(t *testing.T) {
	buf := New()
	err := buf.Check()

	if err != nil {
		t.Error(err)
	} else if buf.Count != 0 {
		t.Errorf("buffer count is %d, but should be 0.", buf.Count)
	}
}

func TestSingleton(t *testing.T) {
	// Valid usages

	counts := []int{
		1, 2, defaultBufferLength, 
		defaultBufferLength + 1,
		MaxBufferCount,
	}

	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}

	for i, count := range counts {
		buf := New()

		for j := 0; j < count; j++ {
			buf.Singleton(item)
		}

		err := buf.Check()
		if err != nil {
			t.Errorf("Test %d: %s", i, err.Error())
			continue
		} else if int(buf.Count) != count {
			t.Errorf(
				"Test %d: buf.Count was %d, but expected %d.",
				i, buf.Count, count,
			)
			continue
		}

		for j := 0; j < len(buf.Buffer); j++ {
			bufItem := buf.Buffer[j].Item
			if bufItem.Type != Uninitialized && item != bufItem {
				t.Errorf(
					"Test %d: Expected Item %v at index %d, but got %v.",
					i, item, j, bufItem,
				)
				break
			}
		}
	}

	// Invalid usages

	buf := New()

	_, err := buf.Singleton(Item{1, Uninitialized, [6]int8{1, 2, 3, 4, 5, 6}})
	if err == nil {
		t.Errorf("No error on Uninitialized input to ItemBuffer.Singleton().")
	}

	for i := 0; i < MaxBufferCount; i++ {
		_, err := buf.Singleton(item)
		if err != nil {
			t.Errorf("After erroneous input: %s", err.Error())
		}
	}

	_, err = buf.Singleton(item)
	if err == nil {
		t.Errorf("No error after trying to add to full ItemBuffer.")
	}
}

func TestIsFull(t *testing.T) {
	buf := New()

	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}

	for i := BufferIndex(0); i < MaxBufferCount; i++ {
		if buf.IsFull() {
			t.Errorf("ItemBuffer.IsFull() returns true after %d items", i)
			return
		}
		buf.Singleton(item)
	}

	if !buf.IsFull() {
		t.Errorf("ItemBuffer.IsFull() returns false when full.")
	}
}

func TestLink(t *testing.T) {
	item := Item{1, TestItem, [6]int8{1,2,3,4,5,6}}

	// Valid usage

	buf := New()

	first, _ := buf.Singleton(item)
	second, _ := buf.Singleton(item)
	third, _ := buf.Singleton(item)

	if err := buf.Link(first, second); err != nil {
		t.Fatalf("Unable to link first and second node: %s", err.Error())
	}

	if err := buf.Link(second, third); err != nil {
		t.Fatalf("Unable to link second and third node: %s", err.Error())
	}

	if err := buf.Check(); err != nil {
		t.Fatalf("ItemBuffer.Check() failed: %s", err.Error())
	}

	fwdRefIndices := []BufferIndex{first, second, third}
	bkdRefIndices := []BufferIndex{third, second, first}
	fwdIndices := listIndices(buf, first, true)
	bkdIndices := listIndices(buf, third, false)

	if !sliceEq(fwdIndices, fwdRefIndices) {
		t.Fatalf(
			"Forward indices are %v instead of %v.",
			fwdIndices, fwdRefIndices,
		)
	}

	if !sliceEq(bkdIndices, bkdRefIndices) {
		t.Fatalf(
			"Backward indices are %v instead of %v.",
			bkdIndices, bkdRefIndices,
		)
	}

	// Invalid usage

	fourth, _ := buf.Singleton(item)

	invalid := BufferIndex(len(buf.Buffer))
	uninitialized := invalid - 1

	if err := buf.Link(first, fourth); err == nil {
		t.Errorf("Reusing forward pointer during Link() succeeded.")
	} else if err := buf.Link(fourth, third); err == nil {
		t.Errorf("Reusing backward pointer during Link() succeeded.")
	} else if err := buf.Link(fourth, invalid); err == nil {
		t.Errorf("Invalid right index to Link() succeeded.")
	} else if err := buf.Link(invalid, fourth); err == nil {
		t.Errorf("Invalid left index to Link() succeeded.")
	} else if err := buf.Link(fourth, uninitialized); err == nil {
		t.Errorf("Uninitialized right index to Link() succeeded.")
	} else if err := buf.Link(uninitialized, fourth); err == nil {
		t.Errorf("Uninitialized left index to Link() succeeded.")
	}
}

func TestUnlink(t *testing.T) {
	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}

	// Valid usage

	tests := []struct {
		nodes int
		unlink BufferIndex
		fwdIndices, bkdIndices []BufferIndex
		fwdStart, bkdStart BufferIndex
	} {
		{3, 2, []BufferIndex{0, 1}, []BufferIndex{1, 0}, 0, 1},
		{3, 1, []BufferIndex{0, 2}, []BufferIndex{2, 0},  0, 2},
		{3, 0, []BufferIndex{1, 2}, []BufferIndex{2, 1}, 1, 2},
	}

	for i, test := range tests {
		buf := New()

		nodes := make([]BufferIndex, test.nodes)
		for j := 0; j < test.nodes; j++ {
			nodes[j], _ = buf.Singleton(item)
			if j > 0 { buf.Link(BufferIndex(j - 1), BufferIndex(j)) }
		}

		err := buf.Unlink(nodes[test.unlink])
		if err != nil {
			t.Errorf("Test %d: %s", i, err.Error())
		}

		fwdIndices := listIndices(buf, test.fwdStart, true)
		if !sliceEq(fwdIndices, test.fwdIndices) {
			t.Errorf(
				"Test %d: Expected %v forward, but got %v",
				i, test.fwdIndices, fwdIndices,
			)
		}

		bkdIndices := listIndices(buf, test.bkdStart, false)
		if !sliceEq(bkdIndices, test.bkdIndices) {
			t.Errorf(
				"Test %d: Expected %v backward, but got %v",
				i, test.bkdIndices, bkdIndices,
			)
		}
	}

	// Invalid usage

	buf := New()

	if err := buf.Unlink(defaultBufferLength); err == nil {
		t.Errorf(
			"Invalid index, %d, did not return error.", 
			defaultBufferLength,
		)
	}

	if err := buf.Unlink(0); err == nil {
		t.Errorf("Uninitialized Item did not result in error.")
	}
}

func listIndices(buf *ListBuffer, idx BufferIndex, fwd bool) []BufferIndex {
	indices := []BufferIndex{}

	for idx != NilIndex {
		indices = append(indices, idx)
		if fwd {
			idx = buf.Incr(idx)
		} else {
			idx = buf.Decr(idx)
		}
	}

	return indices
}

func sliceEq(s1, s2 []BufferIndex) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func TestDelete(t *testing.T) {

	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}

	// Valid usages

	tests := []struct {
		adds, dels []int
	} {
		{[]int{1}, []int{0}},
		{[]int{1}, []int{1}},
		{[]int{2}, []int{1}},
		{[]int{1, 1}, []int{1, 1}},
		{[]int{defaultBufferLength + 1, 50}, []int{10, 10}},
		{[]int{MaxBufferCount}, []int{MaxBufferCount}},
		{[]int{MaxBufferCount, 1}, []int{MaxBufferCount, 0}},
	}

	for i, test := range tests {
		buf := New()

		indices := make([]BufferIndex, MaxBufferCount)
		topIdx := -1
		count := 0

		if len(test.adds) != len(test.dels) {
			panic("Rewrite your tests.")
		}

		for j := 0; j < len(test.adds); j++ {
			count += test.adds[j] - test.dels[j]

			for k := 0; k < test.adds[j]; k++ {
				idx, err := buf.Singleton(item)
				if err != nil {
					t.Fatalf("Test %d, round %d: %s", i, j, err.Error())
				}

				topIdx++
				indices[topIdx] = idx
			}

			for k := 0; k < test.dels[j]; k++ {
				err := buf.Delete(indices[topIdx])
				if err != nil {
					t.Fatalf("Test %d, round %d: %s", i, j, err.Error())
				}

				topIdx--
			}
		}

		if count != int(buf.Count) {
			t.Fatalf(
				"Test %d: Expected count = %d, found count = %d", 
				i, count, buf.Count,
			)
		}

		if err := buf.Check(); err != nil {
			t.Fatalf("Test %d: %s", i, err.Error())
		}
	}

	// Invalid usages

	buf := New()

	if err := buf.Delete(0); err == nil {
		t.Errorf("Deleting uninitialized Item succeeded.")
	} 

	if err := buf.Delete(defaultBufferLength); err == nil {
		t.Errorf("Deleting invalid index succeeded.")
	}
}

func TestGet(t *testing.T) {

	// Valid usages

	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}

	buf := New()

	idx, _ := buf.Singleton(item)
	getItem, err := buf.Get(idx)

	if err != nil {
		t.Errorf("On valid Get() call: %s", err.Error())
	} else if item != getItem {
		t.Errorf("Get Item = %v, but Singleton Item = %v", getItem, item)
	}

	// Invalid usages

	if _, err := buf.Get(defaultBufferLength - 1); err == nil {
		t.Errorf("Invalid index for Get() does not give error.")
	}


	if _, err := buf.Get(MaxBufferCount); err == nil {
		t.Errorf("Uninitialized Item for Get() does not give error.")
	}
}

func TestItemCheck(t *testing.T) {
	item := Item{1, typeLimit, [6]int8{0, 0, 0, 0, 0, 0}}

	if err := item.Check(); err == nil {
		t.Errorf("Invalid item type marked as valid.")
	}
}

// This function is *highly* dependent on the internal representation of
// ItemBuffer. Don't feel bad if you break it, but if it does break
// you need to fix it.
func TestCheck(t *testing.T) {
	buf := New()
	if int(MaxBufferCount) < 1 << 25 { // Or some similarly large number.
		buf.Buffer = make([]Node, MaxBufferCount + 1,)
		if err := buf.Check(); err == nil {
			t.Errorf("Overly long buffer marked as valid.")
		}
	}

	buf.Init()
	buf.Buffer[0].Item.Type = TestItem
	if err := buf.Check(); err == nil {
		t.Errorf("Buffer with invalid Count marked as valid.")
	}

	buf.Init()
	buf.Buffer[0].Item.Type = typeLimit
	if err := buf.Check(); err == nil {
		t.Errorf("Buffer with invalid Item marked as valid.")
	}

	buf.Init()
	buf.Buffer[0].Next = NilIndex
	if err := buf.Check(); err == nil {
		t.Errorf("Buffer with mismatched Next pointer marked as valid.")
	}

	buf.Init()
	buf.Buffer[1].Prev = NilIndex
	if err := buf.Check(); err == nil {
		t.Errorf("Buffer with mismatched Prev pointer marked as valid.")
	}

	buf.Init()
	buf.Buffer[0].Next = 0
	buf.Buffer[0].Prev = 0

	buf.Buffer[1].Prev = NilIndex
	if err := buf.Check(); err == nil {
		t.Errorf("Buffer with circular list marked as valid.")
	}
}

// Bullshit to get 100% coverage.
func TestDecr(t *testing.T) {
	buf := New()
	if NilIndex != buf.Decr(NilIndex) {
		t.Errorf(
			"ItemBuffer.Decr(%d) = %d, not %d",
			NilIndex, buf.Decr(NilIndex), NilIndex,
		)
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkSingletonDefault(b *testing.B) {
	benchmarkSingleton(b, defaultBufferLength)
}

func BenchmarkSingletonDefaultPlusOne(b *testing.B) {
	benchmarkSingleton(b, defaultBufferLength + 1)
}

func BenchmarkSingletonMax(b *testing.B) {
	benchmarkSingleton(b, MaxBufferCount)
}

func benchmarkSingleton(b *testing.B, count int) {
	item := Item{1, TestItem, [6]int8{1, 2, 3, 4, 5, 6}}

	buf := New()
	for i := 0; i < b.N; i++ {
		if i % count == 0 {
			b.StopTimer()
			buf.Init()
			b.StartTimer()
		}

		buf.Singleton(item)
	}
}