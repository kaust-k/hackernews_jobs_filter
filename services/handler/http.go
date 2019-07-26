package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"hnjobs/services/db"
	"hnjobs/services/hnjobs"
)

const tmpl = `
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	<!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

	<title>Happy job searching :)</title>
</head>

<body>
<div class="container-fluid">
{{range .}}
	<div class="card">
		{{htmlSafe .Text .ID .Time}}
	</div>
	<br>
{{end}}
</div>
<!-- Optional JavaScript -->
<!-- jQuery first, then Popper.js, then Bootstrap JS -->
<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
</body>

</html>
`

func HandleFunc() func(w http.ResponseWriter, r *http.Request) {
	return fn
}

var cardTemplate = `
<div class="card-header"> %s </div>
<div class="card-body"> 
	<h6 class="card-subtitle mb-2 text-muted">%s</h6>
	%s 
	<br>
	<a href="https://news.ycombinator.com/item?id=%d" class="btn btn-primary" target="_blank">Discussion Link</a>
</div>
`

func fn(w http.ResponseWriter, r *http.Request) {
	engine := db.GetEngine()

	// Find 'remote' jobs with few extra conditions
	q := engine.Table(&hnjobs.Response{}).
		Where("_text ~* 'backend'").
		Or("_text ~* 'back[-\\ ]end'").
		Or("_text ~* 'software'").
		Or("_text ~* 'developer'").
		Or("_text ~* 'engineer'").
		And("_text ~* 'remote'").
		Omit("kids").
		Desc("_time")
	rows, err := q.Rows(&hnjobs.Response{})

	var arr []*hnjobs.Response
	defer rows.Close()
	for rows.Next() {
		bean := new(hnjobs.Response)
		err = rows.Scan(bean)
		if err == nil {
			arr = append(arr, bean)
		}
	}

	var re = regexp.MustCompile(`(?m)((.*?)(\n|<p>))`)

	t := template.New("fieldname example").Funcs(template.FuncMap{
		"htmlSafe": func(html string, id int, tm hnjobs.JSONDateTime) template.HTML {
			a := re.FindAllStringSubmatch(html, 1)
			title := a[0][2]
			body := html[len(title):]
			tm1 := time.Time(tm)
			return template.HTML(fmt.Sprintf(cardTemplate, title, tm1.Format("2006-01-02 15:04:05"), body, id))
		},
	})
	t, _ = t.Parse(tmpl)

	t.Execute(w, arr)
}
