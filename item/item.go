/*Package item allows for the creation and storage of generic Items.

Items are represented as a reference to a static collection of item behavior
and a dynamic colleciton of instance-specific data.

Since items can be collected into groups, a data structure for storing groups
of lists of Items, a ListBuffer, is also provided.
*/
package item

import (
	"fmt"

	"github.com/phil-mansfield/rogue/error"
)

// Type Type represents all the data for an Item instance which cannot be
// changed at runtime.
type Type uint32

// Type Item represents a single instance of an item. Item.Type references the
// instance's static data and Item.Data references the instance's static data.
//
// If two Item instances cannot be transformed into one another at runtime,
// the difference should not be maintained through the Item.Data, but through
// the Item.Type field.
type Item struct {
	Count uint32
	Type Type
	Data [6]int8
}

// Clear removes all data from the item and marks it as being uninitialized.
func (item *Item) Clear() {
	item.Count = 0
	item.Type = Uninitialized
	for i := 0; i < len(item.Data); i++ {
		item.Data[i] = 0
	}
}

// Check performs consistency checks on the item. An error is returned
// describing the first failed check. If all checks pass, nil is returned.
func (item *Item) Check() *error.Error {
	if item.Type >= typeLimit || item.Type < 0 {
		desc := fmt.Sprintf("Item.Type value %d is invalid.", item.Type)
		return error.New(error.Sanity, desc)
	}

	return nil
}
