package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/markdown"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

var (
	includeTag = regexp.MustCompile(`\{\%\s*include\s*"[^"]+"\s*\%\}`)
	fileQuotes = regexp.MustCompile(`"([^"]+)"`)
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
			if name == `assert/` || (len(name) > 0 && name[0] == '.') {
				return false
			}
			return true
		},
	}
	markdownOptions.Preprocessor = func(c echo.Context, b []byte) []byte {
		if strings.HasPrefix(c.Request().URL().Path(), `/gopl-zh/`) {
			s := string(b)
			s = includeTag.ReplaceAllStringFunc(s, func(v string) string {
				vs := fileQuotes.FindAllString(v, 1)
				if len(vs) > 0 {
					vs[0] = strings.TrimPrefix(vs[0], `"`)
					vs[0] = strings.TrimSuffix(vs[0], `"`)
					fpath := filepath.Join(markdownOptions.Root, c.Request().URL().Path())
					fpath = filepath.Dir(fpath)
					for strings.Contains(vs[0], `../`) {
						fpath = filepath.Dir(fpath)
						vs[0] = strings.Replace(vs[0], `../`, ``, 1)
					}
					fpath = filepath.Join(fpath, vs[0])
					bt, err := ioutil.ReadFile(fpath)
					if err == nil {
						return string(bt)
					}
					println(err.Error())
				}
				return v
			})
			return []byte(s)
		}
		return b
	}
	server := X.Serv()
	server.Use(markdown.Markdown(markdownOptions))

	X.Run(`:` + *port)
}
