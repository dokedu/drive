package mail

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*
var templateFiles embed.FS

type Data struct {
	Name string
	Link string
}

func LoginLinkMailTemplate(name string, link string) (string, error) {
	t, err := template.ParseFS(templateFiles, "templates/*.gohtml")
	if err != nil {
		return "", err
	}

	var data Data

	data.Name = name
	data.Link = link

	out := new(bytes.Buffer)
	err = t.ExecuteTemplate(out, "login_link.gohtml", data)
	if err != nil {
		return "", err
	}

	return out.String(), err
}
