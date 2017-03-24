package main

import (
	"flag"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/markdown"
	X "github.com/webx-top/webx"
)

var (
	includeTag   = regexp.MustCompile(`\{\%\s*include\s*"[^"]+"\s*\%\}`)
	fileQuotes   = regexp.MustCompile(`"([^"]+)"`)
	fileRel      = regexp.MustCompile(`\([^\)]+\)`)
	imgTag       = regexp.MustCompile(`\[[^\]]*\]\(\.\./[^\)]+\)`)
	imgTagCurDir = regexp.MustCompile(`\[[^\]]*\]\([^/][^\)]+\)`)
	linkNumber   = regexp.MustCompile(`\[[\d]+\]:[\s]+[^/][^\s]+[\s]`)
	linkRel      = regexp.MustCompile(`[\s]([^/][^\s]+)[\s]`)
	Parses       = map[string]*ParseSetting{}
)

type ReplaceSetting struct {
	Old    string `json:"old"`
	New    string `json:"new"`
	Regexp string `json:"regexp"`
	regexp *regexp.Regexp
}

type ParseSetting struct {
	Include bool              `json:"include"`
	Link    bool              `json:"link"`
	Replace []*ReplaceSetting `json:"replace"`
}

func main() {
	port := flag.String(`p`, `8080`, `-p 8080`)
	parse := flag.String(`parse`, `/gopl-zh/:include,link;`, `-parse "路径前缀:解析项目清单"`)
	configFile := flag.String(`c`, `data/config/config.yml`, `-c data/config/config.yml`)
	flag.Parse()

	if len(*parse) > 0 {
		*parse = strings.TrimRight(*parse, `;`)
		for _, item := range strings.Split(*parse, `;`) {
			vs := strings.SplitN(item, `:`, 2)
			if len(vs) > 1 {
				c := &ParseSetting{}
				for _, v := range strings.Split(vs[1], `,`) {
					if v == `include` {
						c.Include = true
						continue
					}
					if v == `link` {
						c.Link = true
						continue
					}
				}
				Parses[vs[0]] = c
			}
		}
	}

	if len(*configFile) > 0 {
		if com.FileExists(*configFile) {
			Parses = map[string]*ParseSetting{}
			_, err := confl.DecodeFile(*configFile, &Parses)
			if err != nil {
				log.Error(err)
			}
		}
		actions := &com.MonitorEvent{
			Modify: func(name string) {
				log.Info(`Reload ` + *configFile)
				Parses = map[string]*ParseSetting{}
				_, err := confl.DecodeFile(*configFile, &Parses)
				if err != nil {
					log.Error(err)
				}
				log.Warnf(com.Dump(Parses, true))
			},
		}
		actions.Create = actions.Modify
		configFileName := filepath.Base(*configFile)
		go func() {
			err := com.Monitor(filepath.Dir(*configFile), actions, func(name string) bool {
				return strings.HasSuffix(name, configFileName)
			})
			if err != nil {
				log.Error(err)
			}
		}()
	}
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
		urlPath := c.Request().URL().Path()
		for prefixPath, config := range Parses {
			if !strings.HasPrefix(urlPath, prefixPath) {
				continue
			}
			s := string(b)
			absPath := filepath.Join(markdownOptions.Root, urlPath)
			absPath = filepath.Dir(absPath)
			ppath := path.Dir(urlPath)
			if config.Include {
				//解析 {% include "file.md" %}
				s = includeTag.ReplaceAllStringFunc(s, func(v string) string {
					vs := fileQuotes.FindAllString(v, 1)
					if len(vs) > 0 {
						vs[0] = strings.TrimPrefix(vs[0], `"`)
						vs[0] = strings.TrimSuffix(vs[0], `"`)
						fpath := absPath
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
			}
			if config.Link {
				//修正markdown图片网址
				s = imgTag.ReplaceAllStringFunc(s, func(v string) string {
					vs := fileRel.FindAllString(v, 1)
					if len(vs) > 0 {
						orig := vs[0]
						vs[0] = strings.TrimPrefix(vs[0], `(`)
						vs[0] = strings.TrimSuffix(vs[0], `)`)
						fpath := ppath
						for strings.Contains(vs[0], `../`) {
							fpath = path.Dir(fpath)
							vs[0] = strings.Replace(vs[0], `../`, ``, 1)
						}
						fpath = path.Join(fpath, vs[0])
						v = strings.Replace(v, orig, `(`+fpath+`)`, 1)
					}
					return v
				})

				s = imgTagCurDir.ReplaceAllStringFunc(s, func(v string) string {
					vs := fileRel.FindAllString(v, 1)
					if len(vs) > 0 && !strings.Contains(vs[0], `://`) {
						orig := vs[0]
						vs[0] = strings.TrimPrefix(vs[0], `(`)
						vs[0] = strings.TrimSuffix(vs[0], `)`)
						fpath := ppath
						vs[0] = strings.TrimPrefix(vs[0], `./`)
						fpath = path.Join(fpath, vs[0])
						v = strings.Replace(v, orig, `(`+fpath+`)`, 1)
					}
					return v
				})
				s = linkNumber.ReplaceAllStringFunc(s, func(v string) string {
					vs := linkRel.FindAllString(v, 1)
					if len(vs) > 0 && !strings.Contains(vs[0], `://`) {
						vs[0] = strings.TrimSpace(vs[0])
						orig := vs[0]
						fpath := ppath
						for strings.Contains(vs[0], `../`) {
							fpath = path.Dir(fpath)
							vs[0] = strings.Replace(vs[0], `../`, ``, 1)
						}
						vs[0] = strings.TrimPrefix(vs[0], `./`)
						fpath = path.Join(fpath, vs[0])
						v = strings.Replace(v, orig, fpath, 1)
					}
					return v
				})
			}
			if config.Replace != nil {
				//自定义替换规则
				for _, re := range config.Replace {
					if len(re.Regexp) == 0 {
						s = strings.Replace(s, re.Old, re.New, -1)
					} else {
						if re.regexp == nil {
							var err error
							re.regexp, err = regexp.Compile(re.Regexp)
							if err != nil {
								log.Error(err)
								continue
							}
						}
						s = re.regexp.ReplaceAllString(s, re.New)
					}
				}
			}
			return []byte(s)
		}
		return b
	}
	server := X.Serv()
	server.Core.Use(markdown.Markdown(markdownOptions))

	X.Run(`:` + *port)
}
