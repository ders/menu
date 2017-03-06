// Package menu provides a data structure to represent nested menus such as
// those that might appear on a web page.
//
// The basic unit is a ListItem, which may be either a List or an Item.
// An Item consists of a label (e.g. the text that appears, or perhaps a key
// to a translation table when localized text is needed) and an action (e.g.
// a url path).  A List consists of a label and a list of ListItems.
//
// An Item may optionally have a visibility mask, which may be used when
// certain parts of the menu need to be dynamically hidden.  A example of
// when this might be useful is an admin submenu which is only visible to
// authenticated admin users.
//
// The Filtered() method is used to hide parts of a menu based on a bitmask.
// The provided bitmask is and-ed with the visibility mask of each Item, and
// the item is shown only if the result is nonzero.
//
// A visibility mask may also be applied to a List.  If a List is hidden,
// all the Lists and Items under it are also hidden.
//
// Menus hay have arbitrary depth, although recursive menus are not supported
// and will crash.
//
// Menus play well with json.Marshal (see the examples).
package menu

import (
	"strings"
)

// A ListItem represents either an List or an Item.
type ListItem interface {

	// Display text for this item.  If i18n support is required, put the
	// translation key here.
	Label() string

	// Action to take if this item is selected, e.g. a destination path.
	// If the case of a List, the action should be blank.
	Action() string

	// True if this ListItem is a List.
	IsList() bool

	// For ListItems with a visibility mask, tell us if this ListItem is
	// visible for the given bitmask.  For Items without a visibility mask,
	// this method always returns true.  For Lists, if all the descendants
	// are not visible, then the List is also not visible.
	//
	// Note that empty Lists are always considered not visible.  This has
	// the consequence that a filtered List will never contain ListItems
	// that are empty Lists.
	IsVisible(int64) bool

	// If this is a List, a slice containing all its items.  Returns nil for
	// Items.
	Items() []ListItem

	// If this is a List, return a copy of itself with all the invisible
	// items filtered out.  If this is an Item, return itself.
	Filtered(int64) ListItem

	// Return a string representation with the specified number of tabs prepended
	// to the beginning of every line.  Note that the string representation of a
	// List is multiline.
	IndentedString(int) string

	// Return a string representation.  This is equivalent to IndentedString(0).
	String() string
}

// Type item implements ListItem for single Items.
type item struct {
	Lab        string `json:"label"`
	Act        string `json:"action"`
	visibility int64
}

func (i *item) Label() string             { return i.Lab }
func (i *item) Action() string            { return i.Act }
func (i *item) IsList() bool              { return false }
func (i *item) Items() []ListItem         { return nil }
func (i *item) IsVisible(m int64) bool    { return m&i.visibility != 0 }
func (i *item) Filtered(m int64) ListItem { return i }

func (i *item) IndentedString(n int) string {
	return strings.Repeat("\t", n) + i.String()
}

func (i *item) String() string { return i.Lab + ": " + i.Act }

// Type list implements ListItem for Lists.
type list struct {
	Lab        string     `json:"label"`
	Ims        []ListItem `json:"items"`
	visibility int64
}

func (l *list) Label() string     { return l.Lab }
func (l *list) Action() string    { return "" }
func (l *list) IsList() bool      { return true }
func (l *list) Items() []ListItem { return l.Ims }

func (l *list) IsVisible(m int64) bool {
	if m&l.visibility == 0 {
		return false
	}
	for _, item := range l.Ims {
		if item.IsVisible(m) {
			return true
		}
	}
	return false
}

func (l *list) Filtered(m int64) ListItem {
	filtered := []ListItem{}
	for _, item := range l.Ims {
		if item.IsVisible(m) {
			filtered = append(filtered, item.Filtered(m))
		}
	}
	return &list{Lab: l.Lab, Ims: filtered}
}

func (l *list) IndentedString(n int) string {
	x := []string{strings.Repeat("\t", n) + l.Lab}
	for _, item := range l.Ims {
		x = append(x, item.IndentedString(n+1))
	}
	return strings.Join(x, "\n")
}

func (l *list) String() string { return l.IndentedString(0) }

// NewMenu returns a new List containing the given items.
func NewList(label string, items []ListItem) ListItem {
	return &list{Lab: label, Ims: items, visibility: -1}
}

// NewListMask returns a new List containing the given items and with the
// given visiblity mask.
func NewListMask(label string, items []ListItem, visibility int64) ListItem {
	return &list{Lab: label, Ims: items, visibility: visibility}
}

// NewItem returns a new Item.
func NewItem(label, action string) ListItem {
	return &item{Lab: label, Act: action, visibility: -1}
}

// NewItemMask returns a new Item with the given visibility mask.
func NewItemMask(label, action string, visibility int64) ListItem {
	return &item{Lab: label, Act: action, visibility: visibility}
}
