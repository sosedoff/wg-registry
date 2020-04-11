package service

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/jessevdk/go-assets"
)

func loadTemplate(fs *assets.FileSystem) (*template.Template, error) {
	funcmap := template.FuncMap{
		"time": func(input interface{}) string {
			switch input.(type) {
			case time.Time:
				return input.(time.Time).Format(time.RFC822)
			default:
				return "invalid-value"
			}
		},
	}

	t := template.New("").Funcs(funcmap)

	for name, file := range fs.Files {
		shortname := filepath.Base(name)
		if file.IsDir() || !strings.HasSuffix(shortname, ".html") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(shortname).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
