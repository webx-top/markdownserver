package main

import (
	"flag"
	"path/filepath"

	"github.com/webx-top/echo/middleware/markdown"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

func main() {
	port := flag.String(`p`, `8080`, ``)
	flag.Parse()

	markdownOptions := &markdown.Options{
		Path:   `/`,
		Root:   filepath.Join(com.SelfDir(), `data/markdown`),
		Index:  `index.html`,
		Browse: true,
		Filter: func(name string) bool {
			if name == `assert/` {
				return false
			}
			return true
		},
	}
	server := X.Serv()
	server.Use(markdown.Markdown(markdownOptions))

	X.Run(`:` + *port)
}
