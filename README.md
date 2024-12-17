### sqltemplate
[![Go Reference](https://pkg.go.dev/badge/github.com/ystyle/sqltemplate.svg)](https://pkg.go.dev/github.com/ystyle/sqltemplate)

`sqltemplate` is a Go library designed to harness the power of Go Templates for generating SQL statements. It allows developers to write SQL templates using Go Template syntax and dynamically replace template variables with actual parameter values to produce secure parameterized SQL statements and corresponding parameter lists. This parameterization helps prevent SQL injection attacks while maintaining code flexibility and maintainability.

## Features

- **Security**: Generates SQL statements with `?` placeholders, which are automatically handled by database drivers to prevent SQL injection.
- **Flexibility**: Supports complex SQL templates, including loops (`range`) and conditionals (`if`) statements.
- **Usability**: A simple API that is easy to integrate into existing Go projects.
- **Performance**: An optimized template parsing and execution engine for fast template rendering performance.

### usage

To install `sqltemplate`, run the following command:

```shell
go get github.com/ystyle/sqltemplate
```

```go
package main

import (
	"bytes"
	"github.com/ystyle/sqltemplate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Post struct {
	ID             uint
	Title, Content string
}

func main() {
	// Define the SQL template
	sqltpl := `
insert into posts (id, created_at, title, content) values
{{range $index, $item := .list -}}
({{.ID}}, '2024-12-12 16:09:56', {{ .Title | upper}}, {{.Content}}) {{if last $index $.list}}{{else}}, {{end}}
{{end}}
`
	// Prepare the data
	var Posts = []Post{
		{1, "Aunt Mildred", "bone china tea set"},
		{2, "Uncle John", "moleskin pants"},
		{3, "Cousin Rodney", ""},
	}

	// init sqlte
	tt := sqltemplate.New("master")
	tt.Funcs(sqltemplate.FuncMap{
		"last": func(index int, list interface{}) bool {
			return index == reflect.ValueOf(list).Len()-1
		},
		"upper": func(v string) string {
			return strings.ToUpper(v)
		},
	})
	tt, err = tt.Parse(sqltpl)
	if err != nil {
		panic(err)
	}
	// Use sqltemplate to render the SQL template
	var buf bytes.Buffer
	args, err := tt.Execute(&buf, map[string]any{"list": Posts, "data": "test"})

	if err != nil {
		panic(err)
	}
	// Print rendered sql
	println(buf.String())
	println("args len: ", len(args))
	// Create a database connection
	db, err := gorm.Open(mysql.Open("root:12345678@tcp(127.0.0.1:3306)/test"))
	if err != nil {
		panic(err)
	}
	// Execute the SQL statement
	if err := db.Exec(buf.String(), args...).Error; err != nil {
		panic(err)
	}
}
```
log:
```shell
insert into posts (id, created_at, title, content) values
(?, '2024-12-12 16:09:56', ?, ?), 
(?, '2024-12-12 16:09:56', ?, ?), 
(?, '2024-12-12 16:09:56', ?, ?)
args len: 9

2024/12/16 18:11:10 main.go:56 
[5.831ms] [rows:3] 
insert into posts (id, created_at, title, content) values
(1, '2024-12-12 16:09:56', 'AUNT MILDRED', 'bone china tea set'), 
(2, '2024-12-12 16:09:56', 'UNCLE JOHN', 'moleskin pants'), 
(3, '2024-12-12 16:09:56', 'COUSIN RODNEY', '')
```

### Note on Code Origin

A significant portion of the code in this project is derived from the Go standard library's `text/template` package. This library extends and specializes that functionality for the specific use case of generating SQL statements.
