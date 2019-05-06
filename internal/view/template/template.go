package template

import (
	"bytes"
	"text/template"
)

// Data is the object stores the data to be templated in queries, labels, widgets,
// titles...
type Data map[string]interface{}

// WithData returns the old data + new data in a new Data instance
func (d Data) WithData(data map[string]interface{}) Data {
	if d == nil {
		d = map[string]interface{}{}
	}

	dc := d.deepCopy()
	for k, v := range data {
		dc[k] = v
	}
	return dc
}

func (d Data) deepCopy() Data {
	// Copy vars.
	dc := map[string]interface{}{}
	for k, v := range d {
		dc[k] = v
	}

	return dc
}

// Render will render the template using the object data.
func (d Data) Render(tpl string) string {
	if d == nil {
		d = map[string]interface{}{}
	}

	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return ""
	}

	var b bytes.Buffer
	tmpl.Execute(&b, d)

	return b.String()
}
