package menu

import (
	"testing"
)

type testCase struct {
	Case     ListItem
	Raw      string
	Filtered []filterCase
}

type filterCase struct {
	Mask int64
	Str  string
}

var testCases = []testCase{
	{
		Case: NewList("Test Case 0", []ListItem{
			NewItem("Item A", "/a"),
			NewList("Item B", []ListItem{
				NewItem("Item B-0", "/b0"),
				NewItem("Item B-1", "/b1"),
			}),
			NewItem("Item C", "/c"),
		}),
		Raw: `Test Case 0
	Item A: /a
	Item B
		Item B-0: /b0
		Item B-1: /b1
	Item C: /c`,
		Filtered: []filterCase{
			{
				Mask: 0xff,
				Str: `Test Case 0
	Item A: /a
	Item B
		Item B-0: /b0
		Item B-1: /b1
	Item C: /c`,
			},
		},
	},
	{
		Case: NewList("Test Case 1", []ListItem{
			NewList("Filter Me Out", []ListItem{
				NewItemMask("Thing 1", "/1", 0x10),
				NewItemMask("Thing 2", "/2", 0x10),
				NewItemMask("Thing 3", "/3", 0x08),
			}),
			NewItem("Keep Me", "/keep"),
		}),
		Raw: `Test Case 1
	Filter Me Out
		Thing 1: /1
		Thing 2: /2
		Thing 3: /3
	Keep Me: /keep`,
		Filtered: []filterCase{
			{
				Mask: 0x22,
				Str: `Test Case 1
	Keep Me: /keep`,
			},
			{
				Mask: 0x08,
				Str: `Test Case 1
	Filter Me Out
		Thing 3: /3
	Keep Me: /keep`,
			},
		},
	},
	{
		Case: NewList("Test Case 2", []ListItem{
			NewItem("Hello", "/hello"),
			NewListMask("你好", []ListItem{
				NewItem("很棒", "/great"),
			}, 0x0f),
		}),
		Raw: `Test Case 2
	Hello: /hello
	你好
		很棒: /great`,
		Filtered: []filterCase{
			{
				Mask: 0xf0,
				Str: `Test Case 2
	Hello: /hello`,
			},
		},
	},
}

func TestMenu(t *testing.T) {
	for i, c := range testCases {

		raw := c.Case.String()
		if raw != c.Raw {
			t.Fatalf("Test case %d, raw:\n===Expected===\n%s\n===Got===\n%s\n", i, c.Raw, raw)
		}

		for _, f := range c.Filtered {
			str := c.Case.Filtered(f.Mask).String()
			if str != f.Str {
				t.Fatalf("Test case %d, filtered by %#x:\n===Expected===\n%s\n===Got===\n%s\n", i, f.Mask, f.Str, str)
			}
		}
	}
}
