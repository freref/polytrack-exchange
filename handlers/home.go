package handlers

import (
	"bytes"
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	_, ok := headers[http.CanonicalHeaderKey("HX-Request")]

	if ok {
		tmpl := template.Must(template.ParseFiles("./templates/home.html"))
		tmpl.Execute(w, nil)
	} else {
		indexTmpl := template.Must(template.ParseFiles("./templates/index.html"))
		homeTmpl := template.Must(template.ParseFiles("./templates/home.html"))

		var homeContent bytes.Buffer
		homeTmpl.Execute(&homeContent, nil)

		data := map[string]interface{}{
			"MainContent": template.HTML(homeContent.String()),
		}

		indexTmpl.Execute(w, data)
	}
}
