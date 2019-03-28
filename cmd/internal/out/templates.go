package out

import (
	"io"
	"text/template"
)

func tmpl(w io.Writer, tmplt string, a interface{}) error {
	t := template.New("apizza")
	template.Must(t.Parse(tmplt))
	return t.Execute(w, a)
}

var defaultOrderTmpl = `{{ .OrderName }}
  products:{{ range .Products }}
    {{.Name}}
      code:     {{.Code}}
      options:  {{.Options}}
      quantity: {{.Qty}}{{end}}
  storeID: {{.StoreID}}
  method:  {{.ServiceMethod}}
`

var cartOrderTmpl = `  {{ .OrderName }} - {{ range .Products }} {{.Code}}, {{end}}
`