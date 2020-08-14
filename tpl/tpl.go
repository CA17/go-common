package tpl

import (
	"fmt"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"

	"github.com/ca17/go-common/common"
	"github.com/ca17/go-common/log"
)

type CommonTemplate struct {
	Templates *template.Template
}

func NewCommonTemplate() *CommonTemplate {
	funcMap := template.FuncMap{
		"FenToYuan": common.Fen2Yuan,
		"YuanToFen": common.YuanToFen,
	}
	var templates = template.New("GlobalTemplate").Funcs(funcMap)
	return &CommonTemplate{Templates: templates}
}

func (ct *CommonTemplate) Parse(name, tplcontent string) {
	tplstr := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, tplcontent)
	ct.Templates = template.Must(ct.Templates.Parse(tplstr))
	if log.IsDebug() {
		log.Debugf("parse template %s", name)
	}
}

func (t *CommonTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
