# Menu

Package menu provides a data structure to represent nested menus such as
those that might appear on a web page.

The basic unit is a ListItem, which may be either a List or an Item.
An Item consists of a label (e.g. the text that appears, or perhaps a key
to a translation table when localized text is needed) and an action (e.g.
a url path).  A List consists of a label and a list of ListItems.

An Item may optionally have a visibility mask, which may be used when
certain parts of the menu need to be dynamically hidden.  A example of
when this might be useful is an admin submenu which is only visible to
authenticated admin users.

The Filtered() method is used to hide parts of a menu based on a bitmask.
The provided bitmask is and-ed with the visibility mask of each Item, and
the item is shown only if the result is nonzero.

A visibility mask may also be applied to a List.  If a List is hidden,
all the Lists and Items under it are also hidden.

Menus hay have arbitrary depth, although recursive menus are not supported
and will crash.

Menus play well with json.Marshal (see the examples).
