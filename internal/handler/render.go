package handler

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// NewTemplates РёРЅРёС†РёР°Р»РёР·РёСЂСѓРµС‚ Рё РІРѕР·РІСЂР°С‰Р°РµС‚ С€Р°Р±Р»РѕРЅС‹ РґР»СЏ Gin СЃ РєР°СЃС‚РѕРјРЅС‹РјРё С„СѓРЅРєС†РёСЏРјРё
func NewTemplates() *template.Template {
	funcMap := template.FuncMap{
		"in": func(slice []string, item string) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		},
		"join": strings.Join,
		"toggleTag": func(slice []string, item string) []string {
			for i, v := range slice {
				if v == item {
					return append(slice[:i], slice[i+1:]...)
				}
			}
			return append(slice, item)
		},
		"escape": func(s string) string {
			return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`)
		},
	}

	tmpl := template.New("").Funcs(funcMap)

	err := filepath.WalkDir("views", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			name := filepath.ToSlash(path)
			template.Must(tmpl.New(name).Parse(string(content)))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return tmpl
}
