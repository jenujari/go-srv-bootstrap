package tpl

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed views/*.html
//go:embed views/common/*.html
var views embed.FS

//go:embed assets/*.js
//go:embed assets/*.css
var assets embed.FS

var tpl *template.Template

func init() {
	sub, err := fs.Sub(views, "views")
	if err != nil {
		panic(err)
	}

	tpl = template.Must(template.ParseFS(sub, getTemplateList()...))
}

func GetViewsFs() embed.FS {
	return views
}

func GetAssetsFs() embed.FS {
	return assets
}

func GetTemplateExecutor() *template.Template {
	return tpl
}

func getTemplateList() []string {
	return []string{
		"index.html",
		"common/header.html",
	}
}
