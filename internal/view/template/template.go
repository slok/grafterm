package template

import (
	"bytes"
	"text/template"
)

// Data is the object that all widgets, rows, columns... will have available for their string
// literals so they can populate using the go templating format.
type Data struct {
	Dashboard Dashboard
	Query     Query
}

// Query has the data of the query.
type Query struct {
	DatasourceID string
	Labels       map[string]string
}

// Dashboard has the data of the dashboard.
type Dashboard struct {
	// Range is the range of the dashboard in duration string.
	Range string
}

// WithDashboard Creates a new data struct with the dashboard data applied.
func (d Data) WithDashboard(dash Dashboard) Data {
	dc := d.deepCopy()
	dc.Dashboard = dash
	return dc
}

// WithQuery Creates a new data struct with the query data applied.
func (d Data) WithQuery(q Query) Data {
	dc := d.deepCopy()
	dc.Query = q
	return dc
}

// WithDashboard Creates a new data struct with the dashboard data applied.
func (d Data) deepCopy() Data {
	// Copy dashboard.
	dash := Dashboard{
		Range: d.Dashboard.Range,
	}

	// Copy query.
	labels := map[string]string{}
	for k, v := range d.Query.Labels {
		labels[k] = v
	}
	q := Query{
		DatasourceID: d.Query.DatasourceID,
		Labels:       labels,
	}

	return Data{
		Dashboard: dash,
		Query:     q,
	}
}

// Render will render the template using the object data.
func (d Data) Render(tpl string) string {
	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return ""
	}

	var b bytes.Buffer
	tmpl.Execute(&b, d)

	return b.String()
}
