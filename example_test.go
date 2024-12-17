// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqltemplate

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ExampleTemplate() {
	// Define a template.
	const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.
{{- else}}
It is a shame you couldn't make it to the wedding.
{{- end}}
{{with .Gift -}}
Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`

	// Prepare some data to insert into the template.
	type Recipient struct {
		Name, Gift string
		Attended   bool
	}
	var recipients = []Recipient{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}

	// Create a new template and parse the letter into it.
	t := Must(New("letter").Parse(letter))

	// Execute the template for each recipient.
	for _, r := range recipients {
		args, err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
		println(args)
	}

	// Output:
	// Dear Aunt Mildred,
	//
	// It was a pleasure to see you at the wedding.
	// Thank you for the lovely bone china tea set.
	//
	// Best wishes,
	// Josie
	//
	// Dear Uncle John,
	//
	// It is a shame you couldn't make it to the wedding.
	// Thank you for the lovely moleskin pants.
	//
	// Best wishes,
	// Josie
	//
	// Dear Cousin Rodney,
	//
	// It is a shame you couldn't make it to the wedding.
	//
	// Best wishes,
	// Josie
}

// The following example is duplicated in html/template; keep them in sync.

func ExampleTemplate_block() {
	const (
		master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		overlay = `{{define "list"}} {{join . ", "}}{{end}} `
	)
	var (
		funcs     = FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	masterTmpl, err := New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}
	overlayTmpl, err := Must(masterTmpl.Clone()).Parse(overlay)
	if err != nil {
		log.Fatal(err)
	}
	if err, _ := masterTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
	if err, _ := overlayTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
	// Output:
	// Names:
	// - Gamora
	// - Groot
	// - Nebula
	// - Rocket
	// - Star-Lord
	// Names: Gamora, Groot, Nebula, Rocket, Star-Lord
}

func TestSql(t *testing.T) {
	sqltpl := `insert into posts (created_at, title, content) values
{{range $index, $item := .list -}}
('2024-12-12 16:09:56', {{.Title}}, {{.Content}}) {{if last $index $.list}}  {{else}} , {{end}}
{{end}}`
	type Post struct {
		Title, Content string
	}
	var Posts = []Post{
		{"Aunt Mildred", "bone china tea set"},
		{"Uncle John", "moleskin pants"},
		{"Cousin Rodney", ""},
	}
	tt := New("master")
	tt.Funcs(FuncMap{
		"last": func(index int, list interface{}) bool {
			return index == reflect.ValueOf(list).Len()-1
		},
	})
	tt, err := tt.Parse(sqltpl)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	args, err := tt.Execute(&buf, map[string]any{"list": Posts})
	if err != nil {
		t.Fatal(err)
	}
	println(buf.String())
	println(len(args))
	//insert into posts (created_at, title, content) values
	//('2024-12-12 16:09:56', ?, ?)  ,
	//('2024-12-12 16:09:56', ?, ?)  ,
	//('2024-12-12 16:09:56', ?, ?)
	//
	//6
}
