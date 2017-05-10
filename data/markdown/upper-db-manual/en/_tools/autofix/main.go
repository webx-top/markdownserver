package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Replace struct {
	Old    string
	New    string
	Regexp *regexp.Regexp
}

var replaces = []*Replace{
	/*
		&Replace{` // import "upper.io/db.v2"`, ``, nil},
		&Replace{`"upper.io/db.v2"`, `"github.com/webx-top/db"`, nil},
		&Replace{`"upper.io/db.v2/`, `"github.com/webx-top/db/`, nil},
	*/
	&Replace{``, `${1}../$2/index.md$3`, regexp.MustCompile(`([\("\s])(?:https\://upper\.io)?/db\.v[\d]/([^\)".\s]+)([\)"\s])`)},
	&Replace{``, `${1}../../webroot/res/$2$3$4`, regexp.MustCompile(`([\("])/db\.v[\d]/res/([^\)"]+)(\.png)([\)"])`)},
}

func main() {
	root := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/admpub/upper-db-manual/upper-docs/upper.io/v3/content`)
	//root := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/webx-top/markdownserver/data/markdown/upper-db-manual/en/upper-docs/upper.io/v3/content`)
	save := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/webx-top/markdownserver/data/markdown/upper-db-manual/en/upper-docs/upper.io/v3/content`)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == `_tools` || strings.HasPrefix(info.Name(), `.`) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(info.Name(), `.`) {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(b)
		for _, re := range replaces {
			if re.Regexp == nil {
				content = strings.Replace(content, re.Old, re.New, -1)
			} else {
				content = re.Regexp.ReplaceAllString(content, re.New)
			}
		}
		saveAs := strings.TrimPrefix(path, root)
		saveAs = filepath.Join(save, saveAs)
		err = os.MkdirAll(filepath.Dir(saveAs), os.ModePerm)
		if err == nil {
			file, err := os.Create(saveAs)
			if err == nil {
				_, err = file.WriteString(content)
			}
		}
		if err != nil {
			return err
		}
		fmt.Println(`Autofix ` + path + `.`)
		return nil
	})
	defer time.Sleep(5 * time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(`Autofix complete.`)
}
