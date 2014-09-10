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

// Type Item represents a single Item instance. Item.Type references the
// instance's static data and Item.Data references the instance's static data.
//
// If two Item instances cannot be transformed into one another at runtime,
// the difference should not be maintained through the Item.Data, but through
// the Item.Type field.
type Item struct {
	Type Type
	Data [8]byte
}

// Clear removes all data from item and marks it as being uninitialized.
func (item *Item) Clear() {
	item.Type = Uninitialized
	for i := 0; i < len(item.Data); i++ {
		item.Data[i] = 0
	}
}

// IsValid returns true if all fields of item are consistent.
//
// PROGRAMMER NOTE: An unintialized item is considered valid.
func (item *Item) IsValid() bool {
	return item.Type < typeLimit && item.Type >= 0
}