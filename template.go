package main

const tpl = `#### Comodo
{{- with .Results }}
| Infected      | Result      | Engine      | Updated      |
|:-------------:|:-----------:|:-----------:|:------------:|
| {{.Infected}} | {{.Result}} | {{.Engine}} | {{.Updated}} |
{{ end -}}
`

// func printMarkDownTable(comodo Comodo) {
//
// 	fmt.Println("#### Comodo")
// 	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
// 	table.AddRow(map[string]interface{}{
// 		"Infected": comodo.Results.Infected,
// 		"Result":   comodo.Results.Result,
// 		"Engine":   comodo.Results.Engine,
// 		"Updated":  comodo.Results.Updated,
// 	})
// 	table.Markdown = true
// 	table.Print()
// }
