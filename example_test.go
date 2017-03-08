package menu

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

func ExampleListItem_template() {
	const snippet = `<ul>
{{- range .Items}}
	<li>
	{{- if .IsList}}
		{{.Label}}<br>
		<ol>
		{{- range .Items}}
			<li><a href="{{.Action}}">{{.Label}}</a></li>
	{{- end}}
		</ol>
	{{- else}}
		<a href="{{.Action}}">{{.Label}}</a>
	{{- end}}
	</li>
{{- end}}
</ul>`

	m := NewList("Haha", []ListItem{
		NewItem("About", "/about"),
		NewList("Locations", []ListItem{
			NewItem("Berlin", "/locations/Berlin"),
			NewItem("Seoul", "/locations/Seoul"),
		}),
		NewItemMask("Admin", "/admin", 0x80),
		NewItemMask("FAQ", "/faq", 0x01),
	})

	t := template.Must(template.New("snippet").Parse(snippet))
	err := t.Execute(os.Stdout, m.Filtered(0x01))
	if err != nil {
		panic(err)
	}
	// Output:
	// <ul>
	// 	<li>
	// 		<a href="/about">About</a>
	// 	</li>
	// 	<li>
	// 		Locations<br>
	// 		<ol>
	// 			<li><a href="/locations/Berlin">Berlin</a></li>
	// 			<li><a href="/locations/Seoul">Seoul</a></li>
	// 		</ol>
	// 	</li>
	// 	<li>
	// 		<a href="/faq">FAQ</a>
	// 	</li>
	// </ul>
}

func ExampleListItem_jSON() {
	menu := NewList("", []ListItem{
		NewItem("About", "/about"),
		NewList("Locations", []ListItem{
			NewItem("Berlin", "/locations/Berlin"),
			NewItem("Seoul", "/locations/Seoul"),
		}),
		NewItemMask("Admin", "/admin", 0x80),
		NewItemMask("FAQ", "/faq", 0x01),
	})

	menuJSON, err := json.MarshalIndent(menu.Filtered(0x80), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(menuJSON))

	// Output:
	// {
	//   "label": "",
	//   "items": [
	//     {
	//       "label": "About",
	//       "action": "/about"
	//     },
	//     {
	//       "label": "Locations",
	//       "items": [
	//         {
	//           "label": "Berlin",
	//           "action": "/locations/Berlin"
	//         },
	//         {
	//           "label": "Seoul",
	//           "action": "/locations/Seoul"
	//         }
	//       ]
	//     },
	//     {
	//       "label": "Admin",
	//       "action": "/admin"
	//     }
	//   ]
	// }
}
