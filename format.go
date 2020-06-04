package main

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/martinusso/inflect"
)

func format(edges []edge, count int) (string, error) {
	buf := new(bytes.Buffer)

	switch {
	case count <= 0:
		fmt.Fprintf(buf, "No issues needs review.\n")
		return buf.String(), nil

	case count == 1:
		fmt.Fprintf(buf, "@**all**, %s issue needs review:\n", inflect.IntoWords(float64(count)))

	default:
		fmt.Fprintf(buf, "@**all**, %s issues needs review:\n", inflect.IntoWords(float64(count)))
	}

	tmpl, err := template.New("pr").Parse("\n * {{.Repository.NameWithOwner}}#{{.Number}}: **[{{.Title}}]({{.URL}})**")

	if err != nil {
		return buf.String(), err
	}

	for _, i := range edges {
		err = tmpl.Execute(buf, i.Node.PullRequest)
		if err != nil {
			return buf.String(), err
		}
	}

	return buf.String(), nil
}
