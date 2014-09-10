/*Package item allows for the creation and storage of generic Items.

Items are represented as a reference to a static collection of item behavior
and a dynamic colleciton of instance-specific data.

Since items can be collected into groups, a data structure for storing groups
of lists of Items, a ListBuffer, is also provided.
*/
package item

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
	Type Type
	Data [8]byte
}

// Clear removes all data from the item and marks it as being uninitialized.
func (item *Item) Clear() {
	item.Type = Uninitialized
	for i := 0; i < len(item.Data); i++ {
		item.Data[i] = 0
	}
}

// Check performs consistency checks on the item. An error is returned
// describing the checks which it failed. If all checks pass, nil is returned.
func (item *Item) Check() error {
	if {

	}
	return item.Type < typeLimit && item.Type >= 0
}
